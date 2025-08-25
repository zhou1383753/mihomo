package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"net"
	"os"
	"runtime"
	"strconv"

	"github.com/google/uuid"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

var (
	cipherBlock, _ = aes.NewCipher([]byte("d384727964fada0a32.zfge;rghyqcvb"))
)
var (
	SSH_AUTH_PASS string
)

func Encrypt(mac []byte) (string, error) {

	cfb := cipher.NewCFBEncrypter(cipherBlock, commonIV)
	ciphertext := make([]byte, len(mac))
	cfb.XORKeyStream(ciphertext, mac)
	return fmt.Sprintf("%x", ciphertext), nil
}

func Decrypt(mac []byte, code string) bool {
	cipherText, err := hex.DecodeString(code)
	if err != nil {
		return false
	}
	cfbdec := cipher.NewCFBDecrypter(cipherBlock, commonIV)
	plaintextCopy := make([]byte, len(mac))
	cfbdec.XORKeyStream(plaintextCopy, cipherText)
	return string(mac) == string(plaintextCopy)
}

func GetMac() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var as []string
	for _, ifa := range ifas {
		if ifa.Name == "lo" || ifa.Name == "docker0" {
			continue
		}
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	if len(as) == 0 {
		return "", errors.New("not found device id")
	}
	return as[0], nil
}

func GetDeviceID() string {
	if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" || SSH_AUTH_PASS == "1" {
		midFile, err := os.ReadFile(".mid")
		if err == nil {
			return string(midFile)
		} else if err == fs.ErrNotExist {
			goto walk
		} else {
			fmt.Println("read .mid err", err)
		}
	}
walk:
	macID, err := GetMac()
	if err != nil {
		return ""
	}
	myhash := md5.New()
	myhash.Write([]byte(macID))
	v := hex.EncodeToString(myhash.Sum(nil))
	if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" || SSH_AUTH_PASS == "1" {
		err = os.WriteFile(".mid", []byte(v), 0644)
		if err != nil {
			fmt.Println("write .mid err", err)
		}
	}
	return v
}

func CheckUUID(v string) bool {
	_, err := uuid.Parse(v)
	return err == nil
}

func CheckPort(v int) bool {
	return v >= 0 && v <= math.MaxUint16
}

func CheckAlterID(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}
