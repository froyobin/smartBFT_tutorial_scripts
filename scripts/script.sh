#!/bin/bash
peer chaincode install -n mycc -v 1.2 -p github.com/hyperledger/fabric/examples/chaincode/go
sleep 3	
peer chaincode instantiate -o  bft.frontend.1000:7050 -C channel47 -n mycc -l golang -v 1.2 -c '{"Args":["init","100"]}'
sleep 3
#peer chaincode invoke  -C channel47 -n mycc --waitForEvent -c '{"Args":[ "uploaddomain","helloworld","111111"]}'
#peer chaincode invoke  -C channel47 -n mycc --waitForEvent -c '{"Args":[ "uploadblucktest", "0","100"]}'
# peer chaincode query  -C channel47 -n mycc  -c '{"Args":["query","4"]}'
