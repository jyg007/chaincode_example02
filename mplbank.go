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
	"encoding/json"
	"bytes"
	"regexp"
	"crypto/x509"
    "encoding/pem"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	 pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type account struct {
//	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       		  string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
    CurrentBalance    uint64 `json:"currentbalance"`
	TotalForDay       uint64 `json:"totalforday"`
	CurrentDay        uint64 `json:"currentday"`
	Owner             string `json:"owner"`
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	
  	_, args := stub.GetFunctionAndParameters()
    
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Creation of MPLBANK
	i, _ := strconv.ParseUint(args[0],10,64)
    bank := &account { "MPLBANK", i , 0, 0,  "jyg" }

    bankJSONasBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error(err.Error())
	}

    err = stub.PutState("MPLBANK", bankJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
    
    err = stub.PutState("MPLBANK_DAY", []byte("0"))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}



func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(function)

 

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
	} else if function == "changeday" {
		// the old "Query" is now implemtned in invoke
		return t.changeday(stub)
	} else if function == "getaccounts" {
		// the old "Query" is now implemtned in invoke
		return t.getaccounts(stub)
	}


	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}



// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface,args []string) pb.Response {


	var X uint64          // Transaction value
	var err error


	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

    var DebitAccount, CreditAccount account

   var requester string

    creator, err  := stub.GetCreator()
    n := bytes.Index(creator,[]byte("---"))
    ca := creator[n:]
    //fmt.Println(string(ca))
    if err == nil {
    	block, _ := pem.Decode(ca)
    	if block == nil {
	    	fmt.Println("failed to parse certificate PEM")
	    	requester=""
	    } else {
    		cert, err := x509.ParseCertificate(block.Bytes)
    		if err != nil {
    			fmt.Println("failed to parse certificate: " + err.Error())
    		}
    		//fmt.Println(cert.Subject.CommonName)
    		requester = cert.Subject.CommonName
    	}
    }


  
	
	// Perform the execution
	X, err = strconv.ParseUint(args[2],10,64)
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}


	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	DebitAccountbytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state for debut account")
	}
	if DebitAccountbytes == nil {
		return shim.Error("Entity not found")
	}

    err = json.Unmarshal([]byte(DebitAccountbytes), &DebitAccount)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	

    if (DebitAccount.Name != "MPLBANK")  {
    	if (DebitAccount.Owner != requester) {
    		return shim.Error("Sorry but you are not the owner of this account.  Transaction cancelled")
    	}
    } 



	MPLdaybytes, err := stub.GetState("MPLBANK_DAY")
	if err != nil {
		return shim.Error("Failed to get state")
	}
	MPLday, _ := strconv.ParseUint(string(MPLdaybytes),10,64)
	

	if (DebitAccount.CurrentDay != MPLday) {
		DebitAccount.TotalForDay = 0
	}
	


	

	CreditAccountbytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get state for debut account")
	}

	if CreditAccountbytes == nil {
		fmt.Printf("ouverture de compte %s\n", args[1])
		if ( X > 10000 ) {
		       return shim.Error("Montant demandé trop important")
		};

        // jyg à remplacer par creator
	    CreditAccount = account { args[1], 0, 0, MPLday, requester }

	} else {
    	err = json.Unmarshal([]byte(CreditAccountbytes), &CreditAccount)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to decode JSON of: " + args[1] + "\"}"
			return shim.Error(jsonResp)
		}
	}
		

	if (DebitAccount.TotalForDay + X > 1000) && (DebitAccount.Name != "MPLBANK" ) {
	       return shim.Error("Total amount for fund transfer is superior to 1000")
	};

	if X > DebitAccount.CurrentBalance   {
	      return shim.Error("Insufficient funds in debit account")
	}

	DebitAccount.TotalForDay = DebitAccount.TotalForDay + X
	DebitAccount.CurrentBalance = DebitAccount.CurrentBalance - X
	CreditAccount.CurrentBalance = CreditAccount.CurrentBalance + X

	
	fmt.Printf("DebitNewBalance = %d, CreditNewBalance = %d, TotalTransferForTheDay = %d\n",DebitAccount.CurrentBalance , CreditAccount.CurrentBalance, DebitAccount.TotalForDay)
	

    DebitAccountbytes, err = json.Marshal(DebitAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
    CreditAccountbytes, err = json.Marshal(CreditAccount)
	if err != nil {
		return shim.Error(err.Error())
	}


	// Write the state back to the ledger
	err = stub.PutState(DebitAccount.Name, DebitAccountbytes)
	if err != nil {
		return shim.Error("PutState Debit Account failed")
	}

    err = stub.PutState(CreditAccount.Name, CreditAccountbytes)
	if err != nil {
		return shim.Error("PutState Credit Account failed")
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

// Deletes an entity from state
func (t *SimpleChaincode) changeday(stub shim.ChaincodeStubInterface) pb.Response {
	var  err error
	var MPLday int
	
	//fmt.Println("coucou")
	
	MPLdaybytes, err := stub.GetState("MPLBANK_DAY")
	if err != nil {
		return shim.Error("Failed to get state")
	}
	MPLday, _ = strconv.Atoi(string(MPLdaybytes))
	MPLday++
	
	err = stub.PutState("MPLBANK_DAY", []byte(strconv.Itoa(MPLday)))
	if err != nil {
		return shim.Error(err.Error());
	}	

	return shim.Success([]byte(string(MPLday)))
}



func (t *SimpleChaincode) getaccounts(stub shim.ChaincodeStubInterface) pb.Response {

	resultsIterator, err := stub.GetStateByRange("\"", "}")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	match, _ := regexp.Compile("MPLBANK")

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
	//	queryResultKey, queryResultValue, err := resultsIterator.Next()
		queryResultKey, _ , err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if (!match.MatchString(queryResultKey)) {
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
	//	buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

	//	buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
	//	buffer.WriteString(string(queryResultValue))
	//	buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	fmt.Printf("-  queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	var acc account // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}


	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Accountbytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state for debut account")
	}
	if Accountbytes == nil {
		return shim.Error("Entity not found")
	}

    err = json.Unmarshal([]byte(Accountbytes), &acc)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

    i := strconv.FormatUint(acc.CurrentBalance,10)
	jsonResp := "{\"Name\":\"" + acc.Name + "\",\"Amount\":\"" + i + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) queryplafond(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	var acc account // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}


	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Accountbytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state for debut account")
	}
	if Accountbytes == nil {
		return shim.Error("Entity not found")
	}

    err = json.Unmarshal([]byte(Accountbytes), &acc)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}


	MPLdaybytes, err := stub.GetState("MPLBANK_DAY")
	if err != nil {
		return shim.Error("Failed to get state")
	}
	MPLday, _ := strconv.ParseUint(string(MPLdaybytes),10,64)
	

	if (acc.CurrentDay != MPLday) {
		acc.TotalForDay = 0
	}
	
	i := strconv.FormatUint(acc.TotalForDay,10)
	jsonResp := "{\"Name\":\"" + args[0] + "\",\"Total FT\":\"" + i + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
