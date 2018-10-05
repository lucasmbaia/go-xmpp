package xmpp

import (
  "testing"
  "fmt"
)

func Test_NewClient(t *testing.T) {
  var block = make(chan bool, 1)

  var options = Options{
    Host:	"xmpp",
    Port:	"5222",
    Mechanism:	PLAIN,
    User:	"zeus@localhost",
    Password:	"totvs@123",
  }

  fmt.Println(NewClient(options))
  <-block
}
