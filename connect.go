package xmpp

import (
  "fmt"
  "net"
  //"errors"
  "encoding/xml"

  //"github.com/lucasmbaia/go-xmpp/utils"
)

const (
  XML_STREAM	    = "http://etherx.jabber.org/streams"
  XML_CLIENT	    = "jabber:client"
  XML_TLS	    = "urn:ietf:params:xml:ns:xmpp-tls"
  VERSION	    = "1.0"
  XMPP_DEFAULT_PORT = ":5222"
  STREAM	    = "stream"
)

type Client struct {
  conn	  net.Conn
  domain  string
  enc	  *xml.Encoder
  dec	  *xml.Decoder
}

type Options struct {
  Host	    string
  User	    string
  Password  string
}

func (o *Options) connect() (net.Conn, error) {
  return net.Dial("tcp", o.Host + XMPP_DEFAULT_PORT)
}

func (o *Options) NewClient() (*Client, error) {
  var (
    client  = new(Client)
    err	    error
    conn    net.Conn
    sf	    = new(streamFeatures)
  )

  if conn, err = o.connect(); err != nil {
    return client, err
  }

  client.conn = conn
  client.enc = xml.NewEncoder(client.conn)
  client.dec = xml.NewDecoder(client.conn)

  if sf, err = client.startStream(); err != nil {
    return client, err
  }

  client.startTLSStream(sf)
  return new(Client), nil
}

func (c *Client) init(o *Options) error {
  return nil
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

func (c *Client) startStream() (*streamFeatures, error) {
  var (
    stream  Stream
    err	    error
    sf	    = new(streamFeatures)
    se	    xml.StartElement
  )

  stream = Stream{XMLNSStream: XML_STREAM, XMLNS: XML_CLIENT, To: "localhost", Version: VERSION}

  if err = c.enc.Encode(&stream); err != nil {
    return sf, err
  }

  if se, err = startStream(c.dec); err != nil {
    return sf, err
  }

  if se.Name.Local != STREAM {
    return sf, fmt.Errorf("expected <stream> but got <%v> in %v", se.Name.Local, se.Name.Space)
  }

  if err = c.dec.DecodeElement(sf, nil); err != nil {
    return sf, err
  }

  return sf, nil
}

func startStream(dec *xml.Decoder) (xml.StartElement, error) {
  var (
    t	xml.Token
    err	error
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
