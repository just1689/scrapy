package tools

import (
	"encoding/base64"
)

func Base64ToStr(in string) string {
	outB, _ := base64.StdEncoding.DecodeString(in)
	return string(outB)

}

func StrToBase64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}
