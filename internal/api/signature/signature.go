package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Verifier struct {
	Key string
}

func (v *Verifier) VerifySignature(message, signature string) (bool, error) {
	mac := hmac.New(sha256.New, []byte(v.Key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	messageMAC, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	equal := hmac.Equal(messageMAC, expectedMAC)
	return equal, nil
}
