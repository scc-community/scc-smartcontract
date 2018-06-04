package main

import (
    // "bytes"
	"encoding/json"
	"fmt"
	"strconv"
    // "crypto/x509"
    // "encoding/pem"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract example simple Chaincode implementation
type SmartContract struct {
}

type Account struct {
	Balance  float64
}

var adminUserName = "testadmin"

func (s *SmartContract) Init(APIAPIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (t *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return t.initLedger(APIstub)
	} else if function == "createAccount" {
		return t.createAccount(APIstub, args)
	} else if function == "queryAccount" {
		return t.queryAccount(APIstub, args)
	} else if function == "trading" {
		return t.trading(APIstub, args)
	} else if function == "reward" {
		return t.reward(APIstub, args)
	}

	return shim.Error("Invalid invoke function name")
}

// init ledger
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	var account = Account{Balance: 1000000000.00}
    accountAsBytes, _ := json.Marshal(account)
    APIstub.PutState("xxxxxxxxxxxx", accountAsBytes)
	return shim.Success(nil)
}

// Create account
func (t *SmartContract) createAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
    checkAccount, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state")
	}
    if checkAccount != nil {
        return shim.Error("Account already exists!")
    }

    initBalance, err := strconv.ParseFloat(args[1], 64)
    if err != nil {
            return shim.Error(err.Error())
    }

	var account = Account{Balance: initBalance}

	accountAsBytes, _ := json.Marshal(account)
	APIstub.PutState(args[0], accountAsBytes)

	return shim.Success(nil)
}


// to get account info
func (t *SmartContract) queryAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	account := args[0]

	// Get the state from the ledger
	accountInfo, err := APIstub.GetState(account)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + account + "\"}"
		return shim.Error(jsonResp)
	}

	if accountInfo == nil {
		jsonResp := "{\"Error\":\"Info of " + account + " is nil or maybe this account not exist\"}"
		return shim.Error(jsonResp)
	}

	// jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	// fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(accountInfo)
}

// trading from A account to B account
func (s *SmartContract) trading(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
        
    /////////////////////
    // creatorByte,_:= APIstub.GetCreator()
    // certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
    // if certStart == -1 {
    //         return shim.Error("No certificate found")
    // }
    // certText := creatorByte[certStart:]
    // bl, _ := pem.Decode(certText)
    // if bl == nil {
    //         return shim.Error("cannot decode pem")
    // }

    // cert, err := x509.ParseCertificate(bl.Bytes)
    // if err != nil {
    //         return shim.Error(err.Error())
    // }
    // uname:=cert.Subject.CommonName	
    /////////////////////

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

    if args[0] == args[1] {
        return shim.Error("Cannot transfer to oneself")
    }

    // if args[0] != uname {
    //     return shim.Error("No permission")
    // }

    amount, err := strconv.ParseFloat(args[2], 64)
    if err != nil {
        return shim.Error(err.Error())
    }

    if amount <= 0 {
        return shim.Error("Incorrect amount")
    }

    accountFromAsBytes, _ := APIstub.GetState(args[0])
    accountToAsBytes, _ := APIstub.GetState(args[1])

    if accountFromAsBytes == nil {
        return shim.Error("Account not exist!")
    }

    if accountToAsBytes == nil {
        return shim.Error("Account not exist!")
    }

    accountFrom := Account{}
    accountTo := Account{}

    json.Unmarshal(accountFromAsBytes, &accountFrom)
    json.Unmarshal(accountToAsBytes, &accountTo)

    if accountFrom.Balance < amount {
        return shim.Error("Not enough money")
    }

    accountFrom.Balance = accountFrom.Balance - amount
    accountTo.Balance = accountTo.Balance + amount

    accountFromAsBytes, _ = json.Marshal(accountFrom)
    accountToAsBytes, _ = json.Marshal(accountTo)
    APIstub.PutState(args[0], accountFromAsBytes)
    APIstub.PutState(args[1], accountToAsBytes)

	return shim.Success(nil)
}

// give an account some reward
func (t *SmartContract) reward(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	toAccount := args[0]
	// Get the state from the ledger
	accountToAsBytes, err := APIstub.GetState(toAccount)
	if err != nil {
		return shim.Error("Failed to get state")
	}
    if accountToAsBytes == nil {
        return shim.Error("Account not exist!")
    }
    accountTo := Account{}
    json.Unmarshal(accountToAsBytes, &accountTo)

	if toAccount == adminUserName {
		return shim.Error("user must not be admin")	
	}

	// Get admin from the ledger
	adminValbytes, err := APIstub.GetState(adminUserName)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if adminValbytes == nil {
		return shim.Error("admin Entity not found")
	}
	adminVal, err := strconv.ParseFloat(string(adminValbytes), 64)
	if err != nil {
		return shim.Error("Invalid user amount, expecting a float value")
	}

	// Perform the execution
	amt, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a float value")
	}
	if amt <= 0 {
		return shim.Error("Invalid transaction amount, expecting > 0")
	}
	if adminVal < amt {
		return shim.Error("admin amount no enough")
	}
	accountTo.Balance = accountTo.Balance + amt
	adminVal = adminVal - amt
	adminNewVal := strconv.FormatFloat(adminVal, 'f', -1, 64)
	fmt.Printf("user=%s, userNewVal=%s, adminNewVal\n", toAccount, accountTo.Balance, adminNewVal)

	// Write the state back to the ledger
    accountToAsBytes, err = json.Marshal(accountTo)
	if err != nil {
		return shim.Error(err.Error())
	}
    APIstub.PutState(args[1], accountToAsBytes)

	err = APIstub.PutState(adminUserName, []byte(adminNewVal))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
