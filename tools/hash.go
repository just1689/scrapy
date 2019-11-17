package tools

import (
	"crypto/sha1"
	"fmt"
)

func HashString(in string) string {
	h := sha1.New()
	h.Write([]byte(in))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
