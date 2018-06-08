package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/lucasmbaia/go-xmpp/utils"
	"net"
	"strings"
)

const (
	XML_STREAM        = "http://etherx.jabber.org/streams"
	XML_CLIENT        = "jabber:client"
	XML_TLS           = "urn:ietf:params:xml:ns:xmpp-tls"
	VERSION           = "1.0"
	XMPP_DEFAULT_PORT = "5222"
	STREAM            = "stream"
)

type Client struct {
	conn   net.Conn
	domain string
	enc    *xml.Encoder
	dec    *xml.Decoder
}

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
}

func connect(o Options) (net.Conn, error) {
	var coon net.Conn

	if o.User == "" || o.Password == "" {
		return coon, errors.New("The Option's user and password is required")
	}

	if !strings.Contains(o.User, "@") {
		return coon, errors.New("The format of user is equal the JID format \"user@domain/Resource\"")
	}

	if o.Host == "" {
		return coon, errors.New("The Option's host is required")
	}

	if o.Port == "" {
		o.Port = XMPP_DEFAULT_PORT
	}

	return net.Dial("tcp", fmt.Sprintf("%s:%s", o.Host, o.Port))
}

func NewClient(o Options) (*Client, error) {
	var (
		client = new(Client)
		err    error
		conn   net.Conn
	)

	if conn, err = connect(o); err != nil {
		return client, err
	}

	client.conn = conn
	client.enc = xml.NewEncoder(client.conn)
	client.dec = xml.NewDecoder(client.conn)
	client.domain = strings.Split(strings.Split(o.User, "@")[1], "/")[0]

	if err = client.newClient(); err != nil {
		return client, err
	}

	return client, nil
}

func (c *Client) newClient() error {
	var (
		sf  = new(streamFeatures)
		err error
	)

	if sf, err = c.startStream(); err != nil {
		return err
	}

	if sf, err = c.startTLSStream(sf); err != nil {
		return err
	}

	/*for _, m := range sf.Mechanism.Mechanism {
	}*/
	return nil
}

func (c *Client) startStream() (*streamFeatures, error) {
	var (
		stream Stream
		err    error
		sf     = new(streamFeatures)
		se     xml.StartElement
		st     []byte
	)

	stream = Stream{XMLNSStream: XML_STREAM, XMLNS: XML_CLIENT, To: c.domain, Version: VERSION}

	//remove the end tag stream before the init request
	if st, err = utils.MarshalWithOutEndTag(stream); err != nil {
		return sf, err
	}

	//send the init stream
	if _, err = c.conn.Write(st); err != nil {
		return sf, err
	}

	// get the response of init stream server
	if se, err = startStream(c.dec); err != nil {
		return sf, err
	}

	//check if the server answer is stream
	if se.Name.Local != STREAM {
		return sf, fmt.Errorf("expected <stream> but got <%v> in %v", se.Name.Local, se.Name.Space)
	}

	//server inform of avaliable authenticate mechanisms
	if err = c.dec.DecodeElement(sf, nil); err != nil {
		return sf, err
	}

	return sf, nil
}

func (c *Client) startTLSStream(sf *streamFeatures) (*streamFeatures, error) {
	if sf.TLS == nil {
		return sf, nil
	}

	if sf.TLS.Required == nil {
		return sf, nil
	}

	var err error
	var tls = startTLS{XMLNS: XML_TLS}

	if err = c.enc.Encode(&tls); err != nil {
		return sf, err
	}

	return sf, nil
}

func startStream(dec *xml.Decoder) (xml.StartElement, error) {
	var (
		t   xml.Token
		err error
	)

	for {
		if t, err = dec.Token(); err != nil {
			return xml.StartElement{}, err
		}

		switch t := t.(type) {
		case xml.StartElement:
			return t, nil
		}
	}
}
