// package secp256k1
package main

import (
	"fmt"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcec"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Transaction struct {
	From string `json:"from"`
    To string `json:"to"`
    Amount float64 `json:"amount"`
    Timestamp int64 `json:"timestamp"`
    Version int64 `json:"version"`
}
// Ecrecover returns the uncompressed public key that created the given signature.
func EcrecoverNocgo(hash, sig []byte) ([]byte, error) {
	pub, err := SigToPubNocgo(hash, sig)
	if err != nil {
		return nil, err
	}
	bytes := (*btcec.PublicKey)(pub).SerializeUncompressed()
	return bytes, err
}

// SigToPub returns the public key that created the given signature.
func SigToPubNocgo(hash, sig []byte) (*ecdsa.PublicKey, error) {
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig)

	pub, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, hash)
	return (*ecdsa.PublicKey)(pub), err
}

func SignVerify(from, to string, amount float64, timestamp, version int64, sign string) (bool, error) {
	// hashing for message
	msg := Transaction{
        From: from,
        To: to,
        Amount: amount,
        Timestamp: timestamp,
        Version: version,
	}
	signByte := hexutil.MustDecode(sign)

	msgjsonstr, err := json.Marshal(msg)
    if err != nil {
      fmt.Println(err.Error())  
      return false, err
    }  

    h := sha256.New()
    h.Write([]byte(string(msgjsonstr)))
    msgHash := h.Sum(nil)
    fmt.Println("msg hash="+hex.EncodeToString(msgHash))
	
	// get pubKey by msg and sign
	pubkey, err := EcrecoverNocgo(msgHash[:], signByte)  
	if err != nil || len(pubkey) == 0 || pubkey[0] != 4 {    
      fmt.Println(err.Error())  
      return false, err
	}
	// fmt.Println("pubkey="+hex.EncodeToString(pubkey))

 	// convert pubKey to Address  
    var addr Address  
    copy(addr[:], Keccak256(pubkey[1:])[12:])  
	fmt.Println("address="+"0x" + hex.EncodeToString(addr[:]))

	// address := crypto.PubkeyToAddress(*abc)
	if (("0x"+hex.EncodeToString(addr[:])) == (msg.From)) {
		return true, nil
	} else {
		return false, nil
	}
	
}

// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["createAccount","1222333"]}'
// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["reward","0xaaaa","8.2"]}'
// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["trading","0xaaaa","0xbbbb","3.8"]}'