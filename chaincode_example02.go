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
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	 pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	
  	_, args := stub.GetFunctionAndParameters()
        var A,M string    // Entities
	var Aval int // Asset holdings
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	//fmt.Printf("Aval = %d, Bval = %d\n", Aval)

	M = A+"_OUTTOT@"
	
	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error());
	}
	err = stub.PutState(M, []byte("0"))
	if err != nil {
		return shim.Error(err.Error());
	}

	return shim.Success(nil)
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "queryplafond" {
		// the old "Query" is now implemtned in invoke
		return t.queryplafond(stub, args)
	}


	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}


// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	var A, B, M string    // Entities
	var Aval, Bval, Mval int // Asset holdings
	var X int          // Transaction value
	var err error


	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]
	M = A+"_OUTTOT@"
	
	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

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
		return shim.Error("Failed to get state B")
	}
	
	Mvalbytes, err := stub.GetState(M)
	if err != nil {
		return shim.Error("Failed to get state for M")
	}	
	Mval, _ = strconv.Atoi(string(Mvalbytes))
	
	if Bvalbytes == nil {
		fmt.Printf("ouverture de compte %s\n", B)
		if ( X > 10000 ) {
		       return shim.Error("Montant demandé trop important")
		};
		Bvalbytes=[]byte("0")
	        //return shim.Error("Entity not foud")

	}
		
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	
	Mval = Mval + X
	if (Mval > 1000) && (A != "MPLBANK" ) {
	       return shim.Error("Total amount for fund transfer is superior to 1000")
	};

	if (Bval==0) {
	      // dans le cadre d une creation de compte
	      err = stub.PutState(B+"_OUTTOT@", []byte("0"))
	      if err != nil {
		      return shim.Error("PutStat Failed")
	      }
	};
	
	
	Aval = Aval - X
	Bval = Bval + X

	if Aval < 0  {
	      return shim.Error("Insufficient funds in debit account")
	}
	
	
	fmt.Printf("Aval = %d, Bval = %d, Mval= %d\n", Aval, Bval,Mval)

	
	
	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error("PutState failed")
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error("PutStat Failed")
	}

	err = stub.PutState(M, []byte(strconv.Itoa(Mval)))
	if err != nil {
		return shim.Error(err.Error());
	}

	
	
	
	
	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) queryplafond(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]+"_OUTTOT@"

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + args[0] + "\",\"Total FT\":\"" + string(Avalbytes) + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
