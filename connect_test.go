package xmpp

import (
  "testing"
  "fmt"
)

func Test_NewClient(t *testing.T) {
  var options = Options{
    Host:	"xmpp",
    Mechanism:	PLAIN,
  }

  fmt.Println(options.NewClient())
}
