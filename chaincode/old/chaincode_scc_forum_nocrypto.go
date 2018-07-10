package main

import (
    // "bytes"
	"encoding/json"
	"fmt"
	// "strconv"
    // "crypto/x509"
    // "encoding/pem"

	"github.com/shopspring/decimal"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract example simple Chaincode implementation
type SmartContract struct {
}

type TokenAddress struct {
	TokenName  string
	Address  string
}
type Account struct {
	Balance  string
	OtherTokenAddrList	[]TokenAddress
}

// type PubAndPriKey struct {
// 	PrivateKey  string
// 	PublicKey	string
// }

var adminUserName = "xxxx"
var initAdminBalance = "100"

func (s *SmartContract) Init(APIAPIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (t *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return t.initLedger(APIstub)
	// } else if function == "createPubAndPriKey" {
	// 	return t.createPubAndPriKey(APIstub)
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

// create public key and private key using Secp256k1
// func (s *SmartContract) createPubAndPriKey(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	wallet := NewWallet()
// 	if wallet == nil || wallet.PublicKey == nil || wallet.PrivateKey == nil {
// 		return shim.Error("Failed to get keys")
// 	}
// 	var keyInfo = PubAndPriKey{PublicKey: byteString(wallet.PublicKey), PrivateKey: byteString(wallet.PrivateKey)}
// 	var result, _ = json.Marshal(keyInfo)
//     return shim.Success(result)
// }

// Create account
func (t *SmartContract) createAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var address = args[0]
    checkAccount, err := APIstub.GetState(address)
	if err != nil {
		return shim.Error("Failed to get state")
	}
    if checkAccount != nil {
        return shim.Error("Address already exists!")
    }

	var account = Account{Balance: "0"}

	accountAsBytes, _ := json.Marshal(account)
	APIstub.PutState(address, accountAsBytes)
	fmt.Printf("createAccount Success! address=%s\n", address)

	return shim.Success(nil)
}


// to get account info
func (t *SmartContract) queryAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	address := args[0]

	// Get the state from the ledger
	accountInfo, err := APIstub.GetState(address)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + address + "\"}"
		return shim.Error(jsonResp)
	}

	if accountInfo == nil {
		jsonResp := "{\"Error\":\"Info of " + address + " is nil or maybe this address not exist\"}"
		return shim.Error(jsonResp)
	}

	// jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	// fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(accountInfo)
}

func transfer(APIstub shim.ChaincodeStubInterface, from string, to string, strAmt string, msg []byte, sign []byte) sc.Response {

    if from == to {
        return shim.Error("Cannot transfer to oneself")
    }

    // verify sign
    pubKey, err := RecoverPubkey(msg, sign)
    if err != nil {
        return shim.Error(err.Error())
    }

	if !VerifySignature(pubKey, msg, sign[:len(sign)-1]) {
        return shim.Error("verify sign fail!")
	}

    // if args[0] != uname {
    //     return shim.Error("No permission")
    // }

    // amount, err := strconv.ParseFloat(strAmt, 64)
    amount, err := decimal.NewFromString(strAmt)
    if err != nil {
        return shim.Error(err.Error())
    }

    if(!amount.IsPositive()) {
        return shim.Error("Incorrect amount < 0")
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
    fromBal, err := decimal.NewFromString(accountFrom.Balance)
    if err != nil {
        return shim.Error(err.Error())
    }
    toBal, err := decimal.NewFromString(accountTo.Balance)
    if err != nil {
        return shim.Error(err.Error())
    }

    if fromBal.Cmp(amount) < 0 {
        return shim.Error("Not enough money")
    }

    fromBal = fromBal.Sub(amount)
    toBal = toBal.Add(amount)
    accountFrom.Balance = fromBal.String()
    accountTo.Balance = toBal.String()

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

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	from := args[0]
	to := args[1]
	strAmt := args[2]
	msg := args[3]
	sign := args[4]
	fmt.Printf("trading BEGIN: from=%s, to=%s, amt=%s\n", from, to, strAmt)
	return transfer(APIstub, from, to, strAmt, msg, sign)
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
	return transfer(APIstub, from, to, strAmt, nil, nil)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
