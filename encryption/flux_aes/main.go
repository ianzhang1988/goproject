package main

import (
	"encoding/binary"
	"fmt"
	"os"

	openssl "github.com/spacemonkeygo/openssl"
)

func encrypt(data []byte, cipher *openssl.Cipher, key, iv []byte) ([]byte, error) {

	eCtx, err := openssl.NewEncryptionCipherCtx(cipher, nil, key, iv)
	if err != nil {
		fmt.Println("Could not create encryption context: ", err)
		return nil, err
	}

	cipherbytes, err := eCtx.EncryptUpdate(data)
	if err != nil {
		fmt.Println("EncryptUpdate(plaintext1) failure: ", err)
		return nil, err
	}

	finalBytes, err := eCtx.EncryptFinal()
	if err != nil {
		fmt.Println("EncryptFinal() failure: ", err)
		return nil, err
	}

	cipherbytes = append(cipherbytes, finalBytes...)
	return cipherbytes, nil
}

func decrypt(data []byte, cipher *openssl.Cipher, key, iv []byte) ([]byte, error) {

	dCtx, err := openssl.NewDecryptionCipherCtx(cipher, nil, key, iv)
	if err != nil {
		fmt.Println("Could not create decryption context: ", err)
		return nil, err
	}

	plainbytes, err := dCtx.DecryptUpdate([]byte(data))
	if err != nil {
		fmt.Println("DecryptUpdate(ciphertext part 1) failure: ", err)
		return nil, err
	}

	finalBytes, err := dCtx.DecryptFinal()
	if err != nil {
		fmt.Println("DecryptFinal() failure: ", err)
		return nil, err
	}

	plainbytes = append(plainbytes, finalBytes...)
	return plainbytes, nil
}

func opensslwarper() {
	// CGO_ENABLED=1 go get github.com/spacemonkeygo/openssl

	key := []byte("never gonna give")
	iv := []byte("onna let you dow")
	plaintext1 := "n, never gonna run around"
	plaintext2 := " and desert you"

	cipher, err := openssl.GetCipherByName("aes-128-ctr")
	if err != nil {
		fmt.Println("Could not get cipher: ", err)
		return
	}

	eCtx, err := openssl.NewEncryptionCipherCtx(cipher, nil, key, iv)
	if err != nil {
		fmt.Println("Could not create encryption context: ", err)
		return
	}
	cipherbytes, err := eCtx.EncryptUpdate([]byte(plaintext1))
	if err != nil {
		fmt.Println("EncryptUpdate(plaintext1) failure: ", err)
		return
	}
	ciphertext := string(cipherbytes)
	cipherbytes, err = eCtx.EncryptUpdate([]byte(plaintext2))
	if err != nil {
		fmt.Println("EncryptUpdate(plaintext2) failure: ", err)
		return
	}
	ciphertext += string(cipherbytes)
	cipherbytes, err = eCtx.EncryptFinal()
	if err != nil {
		fmt.Println("EncryptFinal() failure: ", err)
		return
	}
	ciphertext += string(cipherbytes)

	fmt.Println("ciphertext: ", ciphertext)

	dCtx, err := openssl.NewDecryptionCipherCtx(cipher, nil, key, iv)
	if err != nil {
		fmt.Println("Could not create decryption context: ", err)
		return
	}
	plainbytes, err := dCtx.DecryptUpdate([]byte(ciphertext[:15]))
	if err != nil {
		fmt.Println("DecryptUpdate(ciphertext part 1) failure: ", err)
		return
	}
	plainOutput := string(plainbytes)
	plainbytes, err = dCtx.DecryptUpdate([]byte(ciphertext[15:]))
	if err != nil {
		fmt.Println("DecryptUpdate(ciphertext part 2) failure: ", err)
		return
	}

	plainOutput += string(plainbytes)
	plainbytes, err = dCtx.DecryptFinal()
	if err != nil {
		fmt.Println("DecryptFinal() failure: ", err)
		return
	}
	plainOutput += string(plainbytes)

	fmt.Println("output ", plainOutput)
}

func makeKeyIV(secret []byte, keyLen int) ([]byte, []byte, error) {
	secretLen := len(secret)
	if !(secretLen >= 4 && keyLen%8 == 0 && keyLen >= secretLen) {
		fmt.Println("invalid secret lenth")
		return nil, nil, fmt.Errorf("invalid secret lenth:%d/%d", secretLen, keyLen)
	}

	tempData := make([]byte, keyLen*2)

	copy(tempData, secret)

	for idx, _ := range tempData[secretLen:] {
		tempData[secretLen+idx] = 0xff
	}

	seed := make([]byte, 8)

	copy(seed, tempData)

	seedInt := binary.LittleEndian.Uint64(seed)

	random := make([]byte, 8)
	orgin := make([]byte, 8)

	for i := 0; i < keyLen*2; i += 8 {
		// LCG with a = 6364136223846793005, c = 12345, m = 2**64
		seedInt = seedInt*6364136223846793005 + 12345
		binary.LittleEndian.PutUint64(random, seedInt)
		copy(orgin, tempData[i:])
		for idx, _ := range random {
			random[idx] ^= orgin[idx]
		}
		copy(tempData[i:], random)
	}

	return tempData[:keyLen], tempData[keyLen:], nil
}

func main() {
	// opensslwarper()

	cipher, err := openssl.GetCipherByName("aes-128-ctr")
	if err != nil {
		fmt.Println("Could not get cipher: ", err)
		return
	}

	plainText := []byte("Hello")
	key := []byte("never gonna give")
	iv := []byte("onna let you dow")

	plainbytes, _ := encrypt(plainText, cipher, key, iv)
	fmt.Println("plain bytes: ", plainbytes)
	text, _ := decrypt(plainbytes, cipher, key, iv)
	fmt.Println("plain text: ", string(text))

	encryptbytes, _ := os.ReadFile("../encrypted.txt")
	text, _ = decrypt(encryptbytes, cipher, []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01})
	fmt.Println("openssl plain text: ", string(text))

	mykey, myiv, _ := makeKeyIV([]byte{0x01, 0x01, 0x01, 0x01}, 16)
	fmt.Println("key: ", mykey, " iv: ", myiv)

	fluxtest()
}
