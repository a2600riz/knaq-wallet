package security

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func EncodeKeyToHexAndBase64String(keyData []byte) string {
	encodedKey := hex.EncodeToString([]byte(fmt.Sprintf("%s", keyData)))
	return base64.URLEncoding.EncodeToString([]byte(encodedKey))
}
