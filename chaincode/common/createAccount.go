package common

import (
	"fmt"
	// "reflect"
	"crypto/aes"
	"crypto/ecdsa"
    "crypto/cipher"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"io"
	"time"
	crand "crypto/rand"
	"path/filepath"
	"math/big"
	"io/ioutil"
	"os"
	// "../secp256k1"

	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
	// "github.com/pborman/uuid"
	// "github.com/ethereum/go-ethereum/common/math"
	// "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/scrypt"
)

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

type keyStorePassphrase struct {
	keysDirPath string
	scryptN     int
	scryptP     int
}

type Key struct {
	// Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey
}
type Account struct {
	Address Address `json:"address"` // Ethereum account address derived from the key
	URL     URL            `json:"url"`     // Optional resource locator within a backend
}
type URL struct {
	Scheme string // Protocol scheme to identify a capable account backend
	Path   string // Path for the backend to identify a unique entity
}
type cipherparamsJSON struct {
	IV string `json:"iv"`
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


// type BitCurve struct {
// 	P       *big.Int // the order of the underlying field
// 	N       *big.Int // the order of the base point
// 	B       *big.Int // the constant of the BitCurve equation
// 	Gx, Gy  *big.Int // (x,y) of the base point
// 	BitSize int      // the size of the underlying field
// }
// var theCurve = new(BitCurve)

// func init() {
// 	fmt.Println("init1-------")
// 	// See SEC 2 section 2.7.1
// 	// curve parameters taken from:
// 	// http://www.secg.org/collateral/sec2_final.pdf
// 	theCurve.P = math.MustParseBig256("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
// 	theCurve.N = math.MustParseBig256("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
// 	theCurve.B = math.MustParseBig256("0x0000000000000000000000000000000000000000000000000000000000000007")
// 	theCurve.Gx = math.MustParseBig256("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
// 	theCurve.Gy = math.MustParseBig256("0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")
// 	theCurve.BitSize = 256
// 	fmt.Println("init2-------")
// }

// // S256 returns a BitCurve which implements secp256k1.
// func S256BySecp256k1() *BitCurve {
// 	return theCurve
// }

func S256() elliptic.Curve {
	return S256Nocgo()
    // return crypto.S256()
    // return secp256k1.S256()
	// return S256BySecp256k1()
}


func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}
func newKey(rand io.Reader) (*Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
        var a Address
        a.SetBytes(b)
        return a
}
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	pubBytes := FromECDSAPub(&p)
	return BytesToAddress(Keccak256(pubBytes[1:])[12:])
}
func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *Key {
	// id := uuid.NewRandom()
	key := &Key{
		// Id:         id,
		Address:    PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}
func keyFileName(keyAddr Address) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), hex.EncodeToString(keyAddr[:]))
}
func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}
func (ks keyStorePassphrase) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}
// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
    i := len(buf)
    for _, d := range bigint.Bits() {
            for j := 0; j < wordBytes && i > 0; j++ {
                    i--
                    buf[i] = byte(d)
                    d >>= 8
            }
    }
}
// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
    if bigint.BitLen()/8 >= n {
            return bigint.Bytes()
    }
    ret := make([]byte, n)
    ReadBits(bigint, ret)
    return ret
}
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

func GetEntropyCSPRNG(n int) []byte {
	mainBuff := make([]byte, n)
	_, err := io.ReadFull(crand.Reader, mainBuff)
	if err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return mainBuff
}
// EncryptKey encrypts a key using the specified scrypt parameters into a json
// blob that can be decrypted later on.
func EncryptKey(key *Key, auth string, scryptN, scryptP int) ([]byte, error) {
	authArray := []byte(auth)
	salt := GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:16]
	keyBytes := PaddedBigBytes(key.PrivateKey.D, 32)

	iv := GetEntropyCSPRNG(aes.BlockSize) // 16
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	mac := Keccak256(derivedKey[16:32], cipherText)

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

func writeKeyFile(file string, content []byte) error {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}

func (ks keyStorePassphrase) StoreKey(filename string, key *Key, auth string) ([]byte, error) {
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
	a := Account{Address: key.Address, URL: URL{Scheme: "keystore", Path: ks.JoinPath(keyFileName(key.Address))}}
	keyStoreByte, err := ks.StoreKey(a.URL.Path, key, auth)
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
	// return key, a, err
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
