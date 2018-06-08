package xmpp

import (
	"fmt"
	"testing"
)

func Test_NewClient(t *testing.T) {
	var options = Options{
		Host:     "192.168.204.131",
		User:     "zeus@localhost",
		Password: "totvs@123",
	}

	fmt.Println(NewClient(options))
}
