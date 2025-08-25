package main

import "testing"

func Test_sshDial(t *testing.T) {
	sshDial("root", "123456", "172.17.0.5:22")
	sshDial("root", "1234567", "172.17.0.5:22")
}
