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

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("query"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
//	fmt.Println(string(res.Payload))  
}

func checkQuery2(t *testing.T, stub *shim.MockStub, fonc string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(fonc), []byte(value)})
	if res.Status != shim.OK {
		fmt.Println("Query", fonc, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", fonc, "failed to get value")
		t.FailNow()
	}
	//fmt.Println(string(res.Payload))
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println(string(res.Message))
		//t.FailNow()
	}
}

func TestExample02_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=123 B=234
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("MPLBANK"), []byte("9000000000")})

	checkState(t, stub, "MPLBANK", "9000000000")
}



func TestExample02_Invoke(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=567 B=678
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("MPLBANK"), []byte("900000000")})

	// Invoke A->B for 123
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("MPLBANK"), []byte("COMPTE_JYG"), []byte("2000")})
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("MPLBANK"), []byte("COMPTE_KARINE"), []byte("1000")})
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("MPLBANK"), []byte("COMPTE_FABIEN"), []byte("100000")})

	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("COMPTE_JYG"), []byte("COMPTE_KARINE"), []byte("10")})
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("COMPTE_JYG"), []byte("COMPTE_KARINE"), []byte("2")})

	
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("COMPTE_JYG"), []byte("COMPTE_KARINE"), []byte("1100")})

	checkQuery(t, stub, "COMPTE_JYG")
	checkQuery2(t, stub, "queryplafond", "COMPTE_JYG")
	
	checkQuery(t, stub, "COMPTE_KARINE")

	// Invoke B->A for 234
	//checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("B"), []byte("A"), []byte("234")})
	//checkQuery(t, stub, "A", "678")
	//checkQuery(t, stub, "B", "567")
	//checkQuery(t, stub, "A", "678")
	//checkQuery(t, stub, "B", "567")
}


func TestExample02_Query(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=345 B=456
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("MPLBANK"), []byte("900000000")})

	// Query A
	checkQuery(t, stub, "A")

}