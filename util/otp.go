package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net/url"
)

// GenerateSeed creates a random 16 byte Base32 string that serves as OTP Seed
func GenerateSeed() string {
	data := make([]byte, 10)
	_, err := rand.Read(data)
	if err != nil {
		return "error"
	}
	seed := base32.StdEncoding.EncodeToString(data)

	return seed
}

// EncryptSeed encrypts the seed with the user's password as key
// We use AES-128 bit encryption so we reuse the user's password
// as much as needed. The end result is Base32 encoded and will be
// appended to the hashed password in the configuration file
func EncryptSeed(seed string, password []byte) string {
	// Make a key that is 16 byte long, using as many repetitions
	// of the password as needed
	rep := 1 + 16/len(password)
	key := bytes.Repeat(password, rep)

	plaintext := []byte(seed)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key[:16])
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base32.StdEncoding.EncodeToString(ciphertext)
}

// DecryptSeed decrypts the base32 encoded encrypted seed
// with the user's password as key. Reversing the EncryptSeed()
// we use AES-128 bit encryption and the key is generated by using
// the user's password as much as needed.
// The end result is the TOTP seed
func DecryptSeed(encryptedSeed string, password []byte) string {
	// Make a key that is 16 byte long, using as many repetitions
	// of the password as needed
	rep := 1 + 16/len(password)
	key := bytes.Repeat(password, rep)

	ciphertext, _ := base32.StdEncoding.DecodeString(encryptedSeed)

	block, err := aes.NewCipher(key[:16])
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext[:])
}

// ProvisionURI creates an otpauth URI commonly used by mobile phone
// authenticator applications
func ProvisionURI(secret string) string {
	auth := "totp/"
	label := "SignMyKey"
	q := make(url.Values)
	q.Add("secret", secret)
	q.Add("issuer", "SignMyKey")

	return "otpauth://" + auth + label + "?" + q.Encode()
}

// GenerateOTPCode calculates a code as defined in rfc6238
// using the seed and a specific time
// the settings are hardcoded to what is most commonly used:
// 30 seconds window (the code changes every 30 seconds)
// generated code is 6 digits long
// the HMAC is using SHA1
func GenerateOTPCode(seed string, timeval int64) string {

	secretBytes, err := base32.StdEncoding.DecodeString(seed)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 8)
	hash := hmac.New(sha1.New, secretBytes)
	binary.BigEndian.PutUint64(buf, uint64(timeval))
	hash.Write(buf)
	sum := hash.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", 6), mod)
}
