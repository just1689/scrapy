package tools

import "testing"

func TestBase64ToStr(t *testing.T) {

	in := "test"
	out := StrToBase64(in)
	backAgain := Base64ToStr(out)
	if in != backAgain {
		t.Fatal(in, " not equal to ", backAgain)
	}

}
