package message

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var mac = hmac.New(sha256.New, []byte("Spoon!"))

func VerifySignature(message, signature string) (bool, error) {
	mac.Reset()
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	messageMAC, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	equal := hmac.Equal(messageMAC, expectedMAC)
	return equal, nil
}
