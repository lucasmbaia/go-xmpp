package xmpp

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"strings"
)

const ()

type RecvOptions struct {
	Handlers []interface{}
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
				fmt.Println(string(v.Query))
			}
		}
	}()

	<-ctx.Done()
}

func (c *Client) Roster() error {
	var (
		err error
		id  uuid.UUID
	)

	if id, err = uuid.NewV4(); err != nil {
		return err
	}

	if _, err = c.conn.Write([]byte(fmt.Sprintf("<iq type='get' from='%s' id='%s'><query xmlns='jabber:iq:roster'/></iq>", c.jid, id.String()))); err != nil {
		return err
	}

	return nil
}

//XEP-0030: Service Discovery
func (c *Client) DiscoItems(to string) error {
	var (
		err error
		id  uuid.UUID
	)

	if id, err = uuid.NewV4(); err != nil {
		return err
	}

	if _, err = c.conn.Write([]byte(fmt.Sprintf("<iq to='%s' type='get' id='%s'><query xmlns='http://jabber.org/protocol/disco#items'/></iq>", to, id.String()))); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateRoom(name string) error {
	var (
		err  error
		id   uuid.UUID
		form string
	)

	if id, err = uuid.NewV4(); err != nil {
		return err
	}

	form = fmt.Sprintf(FORM_CREATE_ROOM, id.String())

	if _, err = c.conn.Write([]byte(strings.Replace(form, "{name}", name, -1))); err != nil {
		return err
	}

	fmt.Println(strings.Replace(form, "{name}", name, -1))
	return nil
}
