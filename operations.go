package xmpp

import (
  "encoding/xml"
  "context"
  "fmt"
  "github.com/satori/go.uuid"
  "log"
  "strings"
  "reflect"
)

const (
  JABBER_IQ_ROSTER  = "jabber:iq:roster"
  DISCO_ITEMS	    = "http://jabber.org/protocol/disco#items"
  MUC_ONLINE	    = "muc_online"
  MUC_OFFLINE	    = "muc_offline"
)

var (
  minions map[string]Minions
)

type RecvOptions struct {
  Handlers map[string]func(interface{})
}

type Minions struct {
  Containers  []string
}

func NewRecv() RecvOptions {
  return RecvOptions{
    Handlers: make(map[string]func(interface{})),
  }
}

func init() {
  minions = make(map[string]Minions)
}

func (c *Client) Receiver(ctx context.Context, w RecvOptions) {
  go func() {
    for {
      var (
	val	interface{}
	err	error
	q	= Query{}
      )

      if _, val, err = next(c.dec); err != nil {
	log.Println(err)
	continue
      }

      switch v := val.(type) {
      case *Message:
	fmt.Println(v)
      case *Presence:
	fmt.Println(v)

	if !reflect.DeepEqual(v.User, MucUser{}) && !strings.Contains(v.From, c.user) && strings.Contains(v.From, minions@conference.localhost) {
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
	}

	fmt.Println(minions)
      case *clientIQ:
	if v.Type == "result" {
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
	  }
	}

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

func (c *Client) MucPresence() {
  c.conn.Write([]byte("<presence to='minions@conference.localhost/zeus' xml:lang='en'><x xmlns='http://jabber.org/protocol/muc'><history maxchars='0' /></x></presence>"))
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
