package common

import (
	"fmt"
	"testing"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func TestCreateAccount(t *testing.T) {

	password :="star"
	privateKey, address, keyStoreJson, err := CreateAccount(password)
	// address, err := keystore.StoreKey(dir, password, 4096, 6)

	if err != nil {
		t.Errorf("Failed to create account: %v", err)
		// fmt.Printf("Failed to create account: %v", err)
	}

	fmt.Println("privateKey=" + privateKey)
	fmt.Println("address=" + "0x" + address)
	fmt.Println("keyStoreJson=" + keyStoreJson)
}


func TestCreateAccountByGeth(t *testing.T) {

	dir  := "./keystore"
	password :="heliuxing"
	// address, err := createAccount(dir, password, 4096, 6)
	address, err := keystore.StoreKey(dir, password, 4096, 6)

	if err != nil {
		fmt.Printf("Failed to create account: %v", err)
	}
	// keystore.storeNewKey(nil,nil,nil)
	fmt.Println(address)
	fmt.Println("address(address.Hex())="+address.Hex())
	fmt.Println("address(hex.EncodeToString)="+ "0x" + hex.EncodeToString(address[:]))
}
