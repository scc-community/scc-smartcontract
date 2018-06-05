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

var adminUserName = "xxxx"
var initAdminBalance = 100.00

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
	var account = Account{Balance: initAdminBalance}
    accountAsBytes, _ := json.Marshal(account)
    APIstub.PutState(adminUserName, accountAsBytes)
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

func transfer(APIstub shim.ChaincodeStubInterface, from string, to string, strAmt string) sc.Response {
    if from == to {
        return shim.Error("Cannot transfer to oneself")
    }

    // if args[0] != uname {
    //     return shim.Error("No permission")
    // }

    amount, err := strconv.ParseFloat(strAmt, 64)
    if err != nil {
        return shim.Error(err.Error())
    }

    if amount <= 0 {
        return shim.Error("Incorrect amount")
    }

    accountFromAsBytes, _ := APIstub.GetState(from)
    accountToAsBytes, _ := APIstub.GetState(to)

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
    APIstub.PutState(from, accountFromAsBytes)
    APIstub.PutState(to, accountToAsBytes)

	return shim.Success(nil)
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

	from := args[0]
	to := args[1]
	strAmt := args[2]
	fmt.Printf("trading BEGIN: from=%s, to=%s, amt=%s\n", from, to, strAmt)
	return transfer(APIstub, from, to, strAmt)
}

// give an account some reward
func (t *SmartContract) reward(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	from := adminUserName
	to := args[0]
	strAmt := args[1]
	fmt.Printf("reward BEGIN: to=%s, amt=%s\n", to, strAmt)
	return transfer(APIstub, from, to, strAmt)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
