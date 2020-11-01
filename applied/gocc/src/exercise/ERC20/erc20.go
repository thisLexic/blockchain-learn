package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	peer "github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
	"fmt"
)

// TokenChaincode Represents our chaincode object
type TokenChaincode struct {
}

func (token *TokenChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetArgs()
	if len(args) != 5 {
		return shim.Error("Invalid number of arguments provided! There must be 4 arguments")
	}
	if _, err := strconv.ParseFloat(string(args[2]), 64); err != nil {
		return shim.Error("Total Supply must be a valid float!")
	}
	stub.PutState("symbol", args[1])
	stub.PutState("totalSupply", args[2])
	stub.PutState("description", args[3])
	stub.PutState(string(args[4]), args[2])
	return shim.Success([]byte("true"))
}

func (token *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters()

	fmt.Println(funcName)
	fmt.Println(args)

	if funcName == "totalSupply" {
		return totalSupply(stub)
	} else if funcName == "balanceOf" {
		return balanceOf(stub, args) 
	} else if funcName == "transfer" {
		return transfer(stub, args) 
	}

	// This is not good
	return shim.Error(("Bad Function Name = "+funcName+"!!!"))
}


func totalSupply(stub shim.ChaincodeStubInterface) peer.Response {
	value, err := stub.GetState("totalSupply")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("{\"response\": " + string(value) +", \"code\": 0 }"))
}


func balanceOf(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("{\"response\": " + "{\"balance\":"+string(value)+ "}" +", \"code\": 0 }"))
}



func transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fromValue, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	fromFloatValue, err := strconv.ParseFloat(string(fromValue), 64)

	toValue, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	var toFloatValue float64
	if toValue == nil {
		toFloatValue = 0.0
	}else {
		toFloatValue, err = strconv.ParseFloat(string(toValue), 64)
	}

	amount, err := strconv.ParseFloat(string(args[3]), 64)
	newFromFloatValue := fromFloatValue - amount
	newToFloatValue := toFloatValue + amount
	stub.PutState(args[1], []byte(fmt.Sprintf("%d", newFromFloatValue)))
	stub.PutState(args[2], []byte(fmt.Sprintf("%d", newToFloatValue)))
	eventPayload := "{\"from\": "+args[0]+", \"to\": " +args[1] + ",  \"amount\": " +  args[2]+ "}"
	stub.SetEvent("transfer",[]byte(eventPayload))
	return shim.Success([]byte("Success!"))
}


func main() {
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}