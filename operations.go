package xmpp

import (
  "context"
  "log"
  "fmt"
)

type RecvOptions struct {
  Handlers  []interface{}
}

func (c *Client) Receiver(ctx context.Context, w RecvOptions) {
  go func() {
    for {
      var (
	val interface{}
	err error
      )

      if _, val, err = next(c.dec); err != nil {
	log.Println(err)
	continue
      }

      switch v := val.(type) {
      case *Message:
	fmt.Println(v)
      case *clientPresence:
	fmt.Println(v)
      case *clientIQ:
	fmt.Println(v)
      }
    }
  }()

  <-ctx.Done()
}
