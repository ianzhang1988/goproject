package main

import (
	"fmt"
	"os"

	goopenssl "github.com/Luzifer/go-openssl/v4"
	openssl "github.com/spacemonkeygo/openssl"
)

func goOpenssl() {
	plaintext := "Hello World!"
	passphrase := "password"

	o := goopenssl.New()

	enc, err := o.EncryptBytes(passphrase, []byte(plaintext), goopenssl.PBKDF2SHA256)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}

	fmt.Printf("Encrypted text: %s\n", string(enc))

	//data, _ := base64.StdEncoding.DecodeString(string(enc))
	// os.WriteFile("encry", data, os.ModeExclusive)

	// encrytedData, _ := os.ReadFile("../encrypted.txt")
	o2 := goopenssl.New()
	dec, err := o2.DecryptBytes(passphrase, enc, goopenssl.PBKDF2SHA256)

	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}

	fmt.Printf("Decrypted text: %s\n", string(dec))
}

func opensslwarper() {
	// CGO_ENABLED=1 go get github.com/spacemonkeygo/openssl

	key := []byte("never gonna give you up, never g")
	iv := []byte("onna let you dow")
	plaintext1 := "n, never gonna run around"
	plaintext2 := " and desert you"

	cipher, err := openssl.GetCipherByName("aes-256-cbc")
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

	os.ReadFile("../")
	plainbytes, err = dCtx.DecryptUpdate([]byte(ciphertext[15:]))
	if err != nil {
		fmt.Println("DecryptUpdate(ciphertext part 3) failure: ", err)
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

func main() {
	goOpenssl()
	opensslwarper()
}
