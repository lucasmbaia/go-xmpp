package xmpp

import (
  "testing"
)

func Test_NewClient(t *testing.T) {
  var options = Options{
    Host: "172.16.95.111",
  }

  options.NewClient()
}
