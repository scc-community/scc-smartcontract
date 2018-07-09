package commoncgo

import (
	"fmt"
	// "reflect"
	"crypto/aes"
	"crypto/ecdsa"
    "crypto/cipher"
	// "crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"io"
	// "time"
	crand "crypto/rand"
	// "path/filepath"
	"math/big"
	// "io/ioutil"
	// "os"
	// "../secp256k1"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
	// // "github.com/pborman/uuid"
	"github.com/ethereum/go-ethereum/common/math"
	"golang.org/x/crypto/scrypt"
	"github.com/ethereum/go-ethereum/crypto/randentropy"
	"github.com/ethereum/go-ethereum/crypto"
)


type keyStorePassphrase struct {
	keysDirPath string
	scryptN     int
	scryptP     int
}

type Key struct {
	// Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey
}

type encryptedKeyJSONV3 struct {
	Address string     `json:"address"`
	Crypto  cryptoJSON `json:"crypto"`
	// Id      string     `json:"id"`
	Version int        `json:"version"`
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

const (
    // number of bits in a big.Word
    wordBits = 32 << (uint64(^big.Word(0)) >> 63)
    // number of bytes in a big.Word
    wordBytes = wordBits / 8
	scryptR     = 8
	scryptDKLen = 32
	keyHeaderKDF = "scrypt"
	version = 3
	// maxRate is the maximum size of the internal buffer. SHAKE-256
	// currently needs the largest buffer.
	maxRate = 168
)


func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
    // AES-128 is selected due to size of encryptKey.
    aesBlock, err := aes.NewCipher(key)
    if err != nil {
            return nil, err
    }
    stream := cipher.NewCTR(aesBlock, iv)
    outText := make([]byte, len(inText))
    stream.XORKeyStream(outText, inText)
    return outText, err
}
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}
func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *Key {
	// id := uuid.NewRandom()
	key := &Key{
		// Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}
func newKey(rand io.Reader) (*Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}
// EncryptKey encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on.
func EncryptKey(key *Key, auth string, scryptN, scryptP int) ([]byte, error) {
	authArray := []byte(auth)
	salt := randentropy.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:16]
	keyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)

	iv := randentropy.GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	encryptedKeyJSONV3 := encryptedKeyJSONV3{
		hex.EncodeToString(key.Address[:]),
		cryptoStruct,
		// key.Id.String(),
		version,
	}
	return json.Marshal(encryptedKeyJSONV3)
}

func (ks keyStorePassphrase) StoreKey(key *Key, auth string) ([]byte, error) {
	return EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	// keyjson, err := EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	// if err != nil {
	// 	return err
	// }
	// return writeKeyFile(filename, keyjson)
}
func storeNewKey(ks keyStorePassphrase, rand io.Reader, auth string) (string, string, string, error) {
// func storeNewKey(ks keyStorePassphrase, rand io.Reader, auth string) (*Key, Account, error) {
	key, err := newKey(rand)
	if err != nil {
		fmt.Printf("Failed to create address: %v", err)
		return "", "", "", err
	}
	// a := Account{Address: key.Address, URL: URL{Scheme: "keystore", Path: ks.JoinPath(keyFileName(key.Address))}}
	keyStoreByte, err := ks.StoreKey(key, auth)
	if err != nil {
		fmt.Printf("Failed to create keystore: %v", err)
		return "", "", "", err
	}
	// fmt.Printf("----key:%+v\n", key)
	d := key.PrivateKey.D.Bytes()
	b := make([]byte, 0, 32)
	priKeyByte := paddedAppend(32, b, d)
	privateKey := hex.EncodeToString(priKeyByte[:])
	address := hex.EncodeToString(key.Address[:])
	keyStoreJson := string(keyStoreByte[:])

	fmt.Println("privateKey=" + privateKey)
	fmt.Println("address=" + "0x" + address)
	fmt.Println("keyStoreJson=" + keyStoreJson)
	
	return privateKey, address, keyStoreJson, err
}
func CreateAccount(auth string) (string, string, string, error) {
	return storeNewKey(keyStorePassphrase{"", 4096, 6}, crand.Reader, auth)
}

// func main() {

// 	password :="star"
// 	privateKey, address, keyStoreJson, err := CreateAccount(password)
// 	// address, err := keystore.StoreKey(dir, password, 4096, 6)

// 	if err != nil {
// 		fmt.Printf("Failed to create account: %v", err)
// 	}

// 	fmt.Println("privateKey=" + privateKey)
// 	fmt.Println("address=" + "0x" + address)
// 	fmt.Println("keyStoreJson=" + keyStoreJson)
// }
