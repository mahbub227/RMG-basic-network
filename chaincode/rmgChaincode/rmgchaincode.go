/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the product structure, with 4 properties.  Structure tags are used by encoding/json library
type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Owner       string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "fabproduct" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabproduct"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryProduct" {
		return s.queryProduct(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createProduct" {
		return s.createProduct(APIstub, args)
	} else if function == "queryAllProducts" {
		return s.queryAllProducts(APIstub)
	} else if function == "changeProductOwner" {
		return s.changeProductOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	productAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(productAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	products := []Product{
		Product{Name: "T-Shirt", Description: "Color:Blue, Size:XL", Price: "600", Owner: "Yellow"},
		Product{Name: "Full Shirt", Description: "Color:White, Size:M", Price: "1700", Owner: "Ecstacy"},
		Product{Name: "Polo Shirt", Description: "Color:Red, Size:L", Price: "1300", Owner: "Tanjim"},
		Product{Name: "Silk Saree", Description: "Color:Maroon", Price: "10300", Owner: "Rajshahi Silk"},
		Product{Name: "Pant Jeans", Description: "Color:Dark Blue, Size: 34", Price: "3200", Owner: "Yellow"},
		Product{Name: "Jamdani Saree", Description: "Color:Orange", Price: "7600", Owner: "Tangail Kutir"},
	}

	i := 0
	for i < len(products) {
		fmt.Println("i is ", i)
		productAsBytes, _ := json.Marshal(products[i])
		APIstub.PutState("product"+strconv.Itoa(i), productAsBytes)
		fmt.Println("Added", products[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var product = Product{Name: args[1], Description: args[2], Price: args[3], Owner: args[4]}

	productAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[0], productAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllProducts(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "product0"
	endKey := "product999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllProducts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeProductOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	productAsBytes, _ := APIstub.GetState(args[0])
	product := Product{}

	json.Unmarshal(productAsBytes, &product)
	product.Owner = args[1]

	productAsBytes, _ = json.Marshal(product)
	APIstub.PutState(args[0], productAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
