package xmpp

import (
  "context"
  "testing"
  //"fmt"
)

func Test_NewClient(t *testing.T) {
  var ctx = context.Background()

  var options = Options{
    Host:	"xmpp",
    Port:	"5222",
    Mechanism:	PLAIN,
    User:	"zeus@localhost",
    Password:	"totvs@123",
  }

  conn, _ := NewClient(options)
  conn.Receiver(ctx, RecvOptions{})
}
