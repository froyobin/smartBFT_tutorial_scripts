#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e
SOURCEPATH=$(pwd)

DES=$SOURCEPATH/bftsmart-orderer/config/system.config
cp $SOURCEPATH/bftsmart-orderer/config/systemtemplate.config $DES
rm  $SOURCEPATH/bftsmart-orderer/config/currentView ||true
sed -i s/SERVERNAME/$1/g $DES
nodes=$(($1 - 1))
str="0"
for j in `seq 1 $nodes`;
do
	str="$str,$j"
done
sed -i s/INITVIEW/$str/g $DES
set -x
echo "setup the network done"
docker network create -d bridge bft_network

for i in `seq 0 $nodes`;
do 
    echo $i
    docker run -d --rm --network=bft_network --name=bft.node.$i -v $SOURCEPATH/bftsmart-orderer:/etc/bftsmart-orderer bftsmart/fabric-orderingnode:amd64-1.2.0 $i
    sleep 1
done
echo "create $1 pbf nodes"


docker run -d --rm --network=bft_network --name=bft.frontend.1000 bftsmart/fabric-frontend:amd64-1.2.0 1000
#docker run -d --rm --network=bft_network --name=bft.frontend.2000 bftsmart/fabric-frontend:amd64-1.2.0 2000
echo "start the fronted server done"

docker run -d --rm --network=bft_network -v /var/run/:/var/run/  --name=bft.peer.0 hyperledger/fabric-peer:amd64-1.2.0
#docker run -d --rm --network=bft_network -v /var/run/:/var/run/  --name=bft.peer.1 hyperledger/fabric-peer:amd64-1.2.0
echo "start two peers done"
docker run -dit --rm --network=bft_network --name=bft.cli.0  -v $SOURCEPATH/scripts:/scripts -v $SOURCEPATH/CA:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/ -e CORE_PEER_ADDRESS=bft.peer.0:7051 bftsmart/fabric-tools:amd64-1.2.0
#docker run -dit --rm --network=bft_network --name=bft.cli.1  -v $SOURCEPATH/scripts:/scripts -v $SOURCEPATH/CA:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/ -e CORE_PEER_ADDRESS=bft.peer.0:7051 bftsmart/fabric-tools:amd64-1.2.0


echo "create two cli client"

docker exec bft.cli.0 configtxgen -profile SampleSingleMSPChannel -outputCreateChannelTx channel.tx -channelID channel47
sleep 1
docker exec bft.cli.0 configtxgen -profile SampleSingleMSPChannel -outputAnchorPeersUpdate anchor.tx -channelID channel47 -asOrg SampleOrg
echo "genenete artifacts"
echo "we need to sleep for 10 seconds to wait for everything ready"
sleep 10
docker exec bft.cli.0 peer channel create -o bft.frontend.1000:7050 -c channel47 -f channel.tx 
sleep 3
docker exec bft.cli.0 peer channel update -o bft.frontend.1000:7050 -c channel47 -f anchor.tx
echo "channel created."
sleep 3
docker exec bft.cli.0 peer channel join -b channel47.block
echo "join the channel successfully"
sleep 20
docker exec bft.cli.0 /scripts/script.sh
exit
# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
LANGUAGE=${1:-"golang"}
CC_SRC_PATH=github.com/fabcar/go
if [ "$LANGUAGE" = "node" -o "$LANGUAGE" = "NODE" ]; then
	CC_SRC_PATH=/opt/gopath/src/github.com/fabcar/node
fi

# clean the keystore
rm -rf ./hfc-key-store

# launch network; create channel and join peer to channel
cd ../basic-network
./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
docker-compose -f ./docker-compose.yml up -d cli

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n fabcar -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n fabcar -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 10
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n fabcar -c '{"function":"initLedger","Args":[""]}'

printf "\nTotal setup execution time : $(($(date +%s) - starttime)) secs ...\n\n\n"
printf "Start by installing required packages run 'npm install'\n"
printf "Then run 'node enrollAdmin.js', then 'node registerUser'\n\n"
printf "The 'node invoke.js' will fail until it has been updated with valid arguments\n"
printf "The 'node query.js' may be run at anytime once the user has been registered\n\n"
