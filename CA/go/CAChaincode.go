/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CAChaincode example simple Chaincode implementation
type URLData struct {
	Key          string
	Locationdata string
	Padding      []byte
}

type CAChaincode struct {
}

func (t *CAChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	if err != nil {
		return shim.Error("Incorrect parameter for uint64!!")
	}

	println("Initialization done")

	return shim.Success(nil)
}

func (t *CAChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()

	//if function == "uploaddomain" {
	//	return t.uploaddomain(stub, args)
	//}

	if function == "uploaddomain" {
		// Make payment of X units from A to B
		return t.uploaddomain(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "uploadbulktest" {
		return t.uploadbulktest(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *CAChaincode) uploadbulktest(stub shim.ChaincodeStubInterface,
	args[]string) pb.Response{
		fmt.Println("uploading bulk test")
		var i uint64
		start,err := strconv.ParseUint(args[0], 10, 64)
		if err != nil{
			return shim.Error("error in parse argument 0 ")
		}
		end,err := strconv.ParseUint(args[1], 10, 64)
		if err != nil{
			return shim.Error("error in parse argument 1")
		}

		for i=start;i<end;i++{
			args=[]string{strconv.FormatUint(i,10), "TEST"}
			ret := t.uploaddomain(stub,args)
			if ret.Payload != nil{
				return ret
			}
		}

	return shim.Success([]byte("Done"))
}

func (t *CAChaincode) uploaddomain(stub shim.ChaincodeStubInterface,
	args []string) pb.Response {

	var domainname string
	var mappings string
	domainname = args[0]
	mappings = args[1]

	urldata := URLData{
		domainname,
		mappings,
		make([]byte, 1024, 1024),
	}

	storedata, err := json.Marshal(urldata)
	if err != nil {
		return shim.Error("marshal error!")
	}
	err = stub.PutState(domainname, storedata)
	if err != nil {
		fmt.Println("error in put the key")
		return shim.Error("error in put the key")
	}
	//fmt.Println("Successfully upload the domain")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *CAChaincode) invoke2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *CAChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *CAChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var domainname string
	domainname = args[0]
	//fmt.Println(domainname)
	blockdata, err := stub.GetState(domainname)
	if err != nil {
		fmt.Println("error in find the record")
		shim.Error("error in find the record")
	}
	return shim.Success(blockdata)

}

func main() {
	err := shim.Start(new(CAChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
