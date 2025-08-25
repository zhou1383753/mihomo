package util

import (
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {
	fmt.Println(Decrypt([]byte("90:65:84:87:67:1b"), "21eaf6de5abc2008024031c81a1ee265b5"))
}
