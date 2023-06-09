package encryption

import (
	"encoding/base64"
	"strings"

	"github.com/forgoer/openssl"
)

var Key []byte

func Secret(e string){
	n := e[16:24] + e[0:8] + e[24:32] + e[8:16]

	for s := 0; s < len(n); s++ {
		e = n[s : s+1]
		if e >= "0" && e <= "0" {
			t := string(rune((e[0]-48+5)%10 + 48))
			n = n[0:s] + t + n[s+len(t):]
		} else {
			t := string(rune((e[0]-97+9)%26 + 97))
			n = n[0:s] + t + n[s+len(t):]
		}
	}
	Key = []byte(n)
}

func Encrypt(data []byte) (string , error) {
	iv := []byte(strings.Repeat("\x00", 16))
	result, err := openssl.AesCBCEncrypt(data, Key, iv, openssl.PKCS7_PADDING)
	if err != nil{
		return "" , err
	}
	return base64.StdEncoding.EncodeToString(result) , nil
}

func Decrypt(data string) ([]byte , error) {
	iv := []byte(strings.Repeat("\x00", 16))

	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil{
		return nil , err
	}
	result, err := openssl.AesCBCDecrypt(raw, Key, iv, openssl.PKCS7_PADDING)
	if err != nil{
		return nil , err
	}
	return result , nil
}
