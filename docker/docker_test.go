package docker

import (
  "testing"
  "fmt"
)

func TestGenerateActionContainer(t *testing.T) {
  var (
    iq	IQ
    err	error
  )

  if iq, err = ActionContainer(Action{
    To:		"to@localhost",
    From:	"from@localhost",
    Container:	"teste",
    Action:	"STOP",
  }); err != nil {
    t.Fatal(err)
  }

  fmt.Println(iq)
}
