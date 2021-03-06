package main

import (
    // "bytes"
	"encoding/json"
	"fmt"
	"time"
	// "./common"
	"strconv"
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
type AccountInCC struct {
	Balance  string
	OtherTokenAddrList	[]TokenAddress
}
type AccountInfoEx struct {
	PrivateKey  string    `json:"privateKey"`
	Address	string    `json:"address"`
	Keystore string    `json:"Keystore"`
}

const (
	EXPIRE_TIME = 600000
	adminUserName = "xxxx"
	initAdminBalance = "100"
)

// type PubAndPriKey struct {
// 	PrivateKey  string
// 	PublicKey	string
// }

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

	var account = AccountInCC{Balance: initAdminBalance}
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
	var password = args[0]
	if(password == "") {
		return shim.Error("password invalid")
	}

	privateKey, address, keyStoreJson, err := CreateAccount(password)
	if err != nil {
		return shim.Error("Failed to create account! err=" + err.Error())
	}
	if(privateKey == "" || address == "" || keyStoreJson == "") {
		return shim.Error("Failed to create account! return data has null! privateKey=" + privateKey + ",address=" + address + ",keyStoreJson=" + keyStoreJson)
	}

    checkAccount, err := APIstub.GetState(address)
	if err != nil {
		return shim.Error("Failed to get state")
	}
    if checkAccount != nil {
        return shim.Error("Address already exists!")
    }

	var account = AccountInCC{Balance: "0"}

	accountAsBytes, _ := json.Marshal(account)
	APIstub.PutState(address, accountAsBytes)
	fmt.Printf("createAccount Success! address=%s\n", address)

	newAccountInfo := AccountInfoEx{privateKey, address, keyStoreJson}
	result, _ := json.Marshal(newAccountInfo)
    return shim.Success(result)
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

func transfer(APIstub shim.ChaincodeStubInterface, from, to, strAmt, strTimestamp, strVersion, sign string, checkSign bool) sc.Response {

	if(from == "" || to == "" || strAmt == "" || strTimestamp == "" || strVersion == "" || (checkSign == true && sign == "")) {
		return shim.Error("transfer param has null value")
	}

    if from == to {
        return shim.Error("from==to, cannot transfer to oneself")
    }

    amount, err := decimal.NewFromString(strAmt)
	if err != nil {
        return shim.Error("amount format invalid, string convert to decimal fail")
    }
    if(!amount.IsPositive()) {
        return shim.Error("Incorrect amount <= 0")
    }
    floatAmount, bolRes := amount.Float64()
    if bolRes == false {
    	fmt.Printf("decimal to float64 not exactly res=%v, amount=%f\n", bolRes, floatAmount)
    }
	timestamp, err := strconv.ParseInt(strTimestamp, 10, 64)
	if err != nil {
		return shim.Error("timestamp ParseInt fail")
	}

	curTimestamp := time.Now().UnixNano()/1e6
    if(curTimestamp < timestamp || (curTimestamp - timestamp) > EXPIRE_TIME) {
		return shim.Error("Timestamp Expire!")
    }
	version, err := strconv.ParseInt(strVersion, 10, 64)
	if err != nil {
		return shim.Error("version ParseInt fail")
	}
	if (checkSign == true) {
		signRes, err := SignVerify(from, to, floatAmount, timestamp, version, sign)
		if err != nil {
			fmt.Println("SignVerify has Error")
			return shim.Error(err.Error())
		}
		if signRes == false {
			return shim.Error("sign verify fail")
		}
		fmt.Println("SignVerify Success")
	}

    accountFromAsBytes, _ := APIstub.GetState(from)
    accountToAsBytes, _ := APIstub.GetState(to)

    if accountFromAsBytes == nil {
        return shim.Error("Account of from not exist!")
    }

    if accountToAsBytes == nil {
        return shim.Error("Account of to not exist!")
    }

    accountFrom := AccountInCC{}
    accountTo := AccountInCC{}

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

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	from := args[0]
	to := args[1]
	strAmt := args[2]
	strTimestamp := args[3]
	strVersion := args[4]
	sign := args[5]
	fmt.Printf("trading BEGIN: from=%s, to=%s, amt=%s, timestamp=%s, version=%s, sign=%s\n", from, to, strAmt, strTimestamp, strVersion, sign)
	return transfer(APIstub, from, to, strAmt, strTimestamp, strVersion, sign, true)
}

// give an account some reward
func (t *SmartContract) reward(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	from := adminUserName
	to := args[0]
	strAmt := args[1]

	curTimestamp := time.Now().UnixNano()/1e6
	strTimestamp := strconv.FormatInt(curTimestamp,10)
	strVersion := "1"
	sign := ""
	fmt.Printf("reward BEGIN: from=%s, to=%s, amt=%s, timestamp=%s, version=%s, sign=%s\n", from, to, strAmt, strTimestamp, strVersion, sign)
	return transfer(APIstub, from, to, strAmt, strTimestamp, strVersion, sign, false)
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
