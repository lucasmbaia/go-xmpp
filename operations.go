package xmpp

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"strings"
)

const (
	JABBER_IQ_ROSTER = "jabber:iq:roster"
	JABBER_IQ_DOCKER = "jabber:iq:docker"
	DISCO_ITEMS      = "http://jabber.org/protocol/disco#items"
	MUC_ONLINE       = "muc_online"
	MUC_OFFLINE      = "muc_offline"

	PRESENCE_HANDLER = "PRESENCE"
	MESSAGE_HANDLER  = "MESSAGE"
	IQ_HANDLER       = "IQ"
)

var (
	handlers map[string]func(interface{})
	minions  map[string]Minions
)

type Minions struct {
	Containers []string
}

func init() {
	minions = make(map[string]Minions)
	handlers = make(map[string]func(interface{}))
}

func (c *Client) RegisterHandler(plugin string, f func(interface{})) {
	c.Lock()
	if _, ok := handlers[plugin]; !ok {
		handlers[plugin] = f
	}
	c.Unlock()
}

func (c *Client) Receiver(ctx context.Context) {
	go func() {
		for {
			var (
				val interface{}
				err error
				//q   = Query{}
				//ok  bool
			)

			if _, val, err = next(c.dec); err != nil {
				log.Println(err)
				continue
			}

			switch v := val.(type) {
			case *Message:
				c.sendHandler(MESSAGE_HANDLER, v)
			case *Presence:
				c.sendHandler(PRESENCE_HANDLER, v)
				/*if !reflect.DeepEqual(v.User, MucUser{}) && !strings.Contains(v.From, c.user) && strings.Contains(v.From, "minions@conference.localhost") {
					for _, item := range v.User.Item {
						if v.Type == "unavailable" {
							c.Lock()
							if _, ok := minions[item.Jid]; ok {
								delete(minions, item.Jid)
							}
							c.Unlock()
						} else {
							c.Lock()
							if _, ok := minions[item.Jid]; !ok {
								minions[item.Jid] = Minions{}
							}
							c.Unlock()
						}
					}
				}*/
			case *clientIQ:
				c.sendHandler(IQ_HANDLER, v)
				/*if v.Type == "result" {
					if err = xml.Unmarshal(v.Query, &q); err != nil {
						continue
					}

					switch q.XMLName.Space {
					case JABBER_IQ_ROSTER:
						var roster QueryRoster
						if err = xml.Unmarshal(v.Query, &roster); err != nil {
							continue
						}

						fmt.Println(roster)
					case DISCO_ITEMS:
						var di QueryDiscoItems
						if err = xml.Unmarshal(v.Query, &di); err != nil {
							continue
						}

						fmt.Println(di)
					case JABBER_IQ_DOCKER:

					}
				}

				fmt.Println(v)
				fmt.Println(string(v.Query))*/
			}
		}
	}()

	<-ctx.Done()
}

func (c *Client) sendHandler(plugin string, v interface{}) {
	c.Lock()
	if _, ok := handlers[plugin]; ok {
		handlers[plugin](v)
	}
	c.Unlock()
}

func (c *Client) Send(i interface{}) error {
	var (
		body []byte
		err  error
	)

	if body, err = xml.Marshal(i); err != nil {
		return err
	}

	if _, err = c.conn.Write(body); err != nil {
		return err
	}

	return nil
}

func (c *Client) Roster() error {
	var (
		err error
		id  uuid.UUID
	)

	if id, err = uuid.NewV4(); err != nil {
		return err
	}

	if _, err = c.conn.Write([]byte(fmt.Sprintf("<iq type='get' from='%s' id='%s'><query xmlns='jabber:iq:roster'/></iq>", c.Jid, id.String()))); err != nil {
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

func (c *Client) MucPresence(room string) error {
	var body = fmt.Sprintf("<presence to='%s/%s' xml:lang='en'><x xmlns='http://jabber.org/protocol/muc'><history maxchars='0' /></x></presence>", room, c.user)

	if _, err := c.conn.Write([]byte(body)); err != nil {
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

	return nil
}
