// package secp256k1
package commoncgo

import (
	"fmt"
	// "crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"crypto/sha256"
	// "unsafe"
	// "bytes"
	// "crypto/ecdsa"
	// "reflect"
	// "testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/common/math"
)
var (
//	testmsg     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testmsg     = hexutil.MustDecode("0x304c3ee287d9dd1af4a8cb49bfac6608761e6f3d3a7914bb34765ed3028e4076")
//  testsig     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testsig     = hexutil.MustDecode("0x5ed7884c070b2e2d1cc782e2f40075be53b8e243e01a0485a765b63733ea060f34bb61bb23890fab17dc5ba4407db6218bf3bff7a4479e3458d84a291c96bd7201")
	testsig1     = hexutil.MustDecode("0x5ed7884c070b2e2d1cc782e2f40075be53b8e243e01a0485a765b63733ea060f34bb61bb23890fab17dc5ba4407db6218bf3bff7a4479e3458d84a291c96555501")
	// r=5ed7884c070b2e2d1cc782e2f40075be53b8e243e01a0485a765b63733ea060f
	// s=34bb61bb23890fab17dc5ba4407db6218bf3bff7a4479e3458d84a291c96bd72
	// sign=r+s
//	testpubkey  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkey  = hexutil.MustDecode("0x5da6d2d2ac9d38d82e263ffce703485fdbdc45951e24410fb7c8ae9f03427d467e4d4c50fd59423979f2de9421c8be52b144564434a60a8c17aa767477eb048f")
	// testpubkeyc = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
	testseckey  = hexutil.MustDecode("0xc793e7a9a52074f2574ca09e945ae2a73aaf441c741befd8bcd5ad9a57a8f282")
)
type Transaction struct {
	From string `json:"from"`
    To string `json:"to"`
    Amount float64 `json:"amount"`
    Timestamp int `json:"timestamp"`
    Version int `json:"version"`
}

func createAccount() {
	dir  := "/Users/heliuxing/Documents/hyperledger/scc-smartcontract/chaincode/secp256k1/keystore"
	password :="heliuxing123"
	address, err := keystore.StoreKey(dir, password, 4096, 6)

	if err != nil {
		fmt.Printf("Failed to create account: %v", err)
	}
	// fmt.Println("address="+hex.EncodeToString(address[0]))
	fmt.Println(address)
}

func main() {
	fmt.Println("begin")
	// sig := testsig[:len(testsig)-1]
	// // sig = testsig
	// // fmt.Println(hex.EncodeToString(sig))
	// if VerifySignature(testpubkey, testmsg, sig) {
	// 	fmt.Println("verify success")
	// } else {
	// 	fmt.Println("verify fail")
	// }

	// ///////////////////
	// sign, _ := Sign(testmsg, testseckey)
	// fmt.Println("sign="+hex.EncodeToString(sign))	//5ed7884c070b2e2d1cc782e2f40075be53b8e243e01a0485a765b63733ea060f34bb61bb23890fab17dc5ba4407db6218bf3bff7a4479e3458d84a291c96bd7201

	// // 传进来msg、sign
	// pubkey, _ := RecoverPubkey(testmsg, testsig1)
	// fmt.Println("pubkey="+hex.EncodeToString(pubkey))

	// if VerifySignature(pubkey, testmsg, testsig1[:len(testsig1)-1]) {
	// 	fmt.Println("verify success")
	// } else {
	// 	fmt.Println("verify fail")
	// }

	/////////////////
	// 使用sha256对消息进行hash
	msg123 := Transaction{
        From: "0x84d224882ddae50b760d22efe739adeb0f58e9f0",
        To: "b",
        Amount: 3.81,
        Timestamp: 1408782304939,
        Version: 1,
	}
	sign := hexutil.MustDecode("0x11e39a48e651fb7c8f0e8a2dd79c394ae8939da4af2283b958ecfc27f50ae8852e8b1d0510ee32785b4d09c1245267a7a36441605888cef17ff0841d00e72fe601")

	msgjsonstr, err := json.Marshal(msg123)
    if err != nil {  
      fmt.Println(err.Error())  
      return
    }  
	// {"from":"0x37d21ec4f826e1c14ef8a29b3aa5bb971d0ea1db","to":"b","amount":3.8,"timestamp":1528782304939,"version":1}
	fmt.Println(string(msgjsonstr))//614cdafa2f2767711c9c44da8367e24112f898c72a8c1a96cf143946f846cdc0
    h := sha256.New()
    h.Write([]byte(string(msgjsonstr)))
    msgHash := h.Sum(nil)
    fmt.Println("msg hash="+hex.EncodeToString(msgHash))

	pubkey, _ := secp256k1.RecoverPubkey(msgHash, sign)
	fmt.Println("pubkey="+hex.EncodeToString(pubkey))
	
	// 对比sign
	// if secp256k1.VerifySignature(pubkey, msgHash, sign[:len(sign)-1]) {
	// 	fmt.Println("verify success")
	// } else {
	// 	fmt.Println("verify fail")
	// }

	// pubkey to address
	// abc := *(*ecdsa.PublicKey)(unsafe.Pointer(&pubkey))
	// fmt.Printf("%+v\n", abc)
	pub, err := crypto.Ecrecover(msgHash[:], sign)  
     if err != nil || len(pub) == 0 || pub[0] != 4 {  
         return 
     }  
     fmt.Println("pubkey1="+hex.EncodeToString(pub))
 // convert pubKey to Address  
     var addr common.Address  
     copy(addr[:], crypto.Keccak256(pub[1:])[12:])  
	fmt.Println("address="+"0x" + hex.EncodeToString(addr[:]))

	// address := crypto.PubkeyToAddress(*abc)
	// fmt.Println(address)


	fmt.Println("finish")

	
}
// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["createAccount","0xaaaa"]}'
// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["reward","0xaaaa","8.2"]}'
// peer chaincode invoke -o orderer:7050 -C sccchannel -n scc -c '{"Args":["trading","0xaaaa","0xbbbb","3.8"]}'