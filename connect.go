package xmpp

import (
  "fmt"
  "net"
  "errors"
  "crypto/tls"
  "crypto/x509"
  "crypto/rand"
  "encoding/base64"
  "encoding/binary"
  "strings"
  //"errors"
  "encoding/xml"

  "github.com/lucasmbaia/go-xmpp/utils"
)

const (
  nsSASL	    = "urn:ietf:params:xml:ns:xmpp-sasl"
  PLAIN		    = "PLAIN"
  BINARY_SALS	    = "\x00"
  XML_STREAM	    = "http://etherx.jabber.org/streams"
  XML_CLIENT	    = "jabber:client"
  XML_TLS	    = "urn:ietf:params:xml:ns:xmpp-tls"
  VERSION	    = "1.0"
  XMPP_DEFAULT_PORT = ":5222"
  STREAM	    = "stream"
  rootPEM	    = `
-----BEGIN CERTIFICATE-----
MIIFnDCCA4SgAwIBAgIJAJrQSXXz8sTrMA0GCSqGSIb3DQEBCwUAMGMxCzAJBgNV
BAYTAkZSMQ4wDAYDVQQIDAVQYXJpczEOMAwGA1UEBwwFUGFyaXMxEzARBgNVBAoM
ClByb2Nlc3NPbmUxCzAJBgNVBAsMAklUMRIwEAYDVQQDDAlsb2NhbGhvc3QwHhcN
MTcxMjI4MTQyOTA4WhcNNDUwNTE1MTQyOTA4WjBjMQswCQYDVQQGEwJGUjEOMAwG
A1UECAwFUGFyaXMxDjAMBgNVBAcMBVBhcmlzMRMwEQYDVQQKDApQcm9jZXNzT25l
MQswCQYDVQQLDAJJVDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0B
AQEFAAOCAg8AMIICCgKCAgEAwVgMQ9yhHEq+pa6HCDIkq7mk5HVzFhsmDqO0EUPY
PAgQITFZ9ODxJ88t7Q44IQJGHpxedzgvCLsJBymqZDFIh2N35+bLmuDjRXBa+iO8
MIgoSiAToK0gD3s8CQb6AGR52U8+Qe9F/TQxZwG1vPKPv2RQTVnukSumR3lM0EJX
BWAemB/X55BlZhmphsUhMGxLZmPkrR0Kt1gvGpBdSvdaiWVaMAMUcxExEblD3spb
xhBynrZNFGNVcE1uKzjdwJoqG7KsuN3I6d5FLa6SGW/lz+qfedMiU8Heve9qjoFE
prFAzgTntowzPM5O8FulldF/6VwRwvCprxeuwg07gP444mLujuUYJIUe9tw9VWqq
UsQ03tXzPaoqnXQQvNiIhimmlqRqnvLLIqGmPQ+yVbRka1MK+GWE/voQu1MttjhX
sgha3iRezoVxKZxYsJxdwwrTJYKjmXB+m8Q/oemjx+/rs91ivPAQCwU+i0YboUd8
ntxR54/vrNwN+s0Jq42tHNPRLAk+aBd47fba3Jdz0OKcxae+X7TXkdjgMDFdtMJg
iqzPiBHI9QnrQDaxufPPwC+XRnx7GrfDBjXsNHJingApuxisg4mIfIsCHkQChJvJ
uVJUN69jYigwTeJbEaFmJDM6gljSp6VKJUTrERgvt+NFmP9RH9jIy7s7CyAFvpNA
3n8CAwEAAaNTMFEwHQYDVR0OBBYEFNM1dMqG43cHn9fLDOY7hf57aCbwMB8GA1Ud
IwQYMBaAFNM1dMqG43cHn9fLDOY7hf57aCbwMA8GA1UdEwEB/wQFMAMBAf8wDQYJ
KoZIhvcNAQELBQADggIBAAkRAw4oFokXDWbG0DTY27P42lyeQsWzD3v98nkovCSW
cvIDt5JYG+YFPMTqGC/9OdghMFdiJe+t47R4gbeAUbXw8Ckv7MHHTpDLd3e2oOf5
ByZ57mguEyrcLiZ+08ZcPwtV62WptuU4hi3P+gnXQt9ebIrb9128CvbponmgvUH3
e0Tr1CiJcDUbYOY5bZWv5K1OxAoXVuVVL1nDADyVZXbaR+bIpe+1y3g+PPp4MQ0M
SVAEBI8QY0NR0AFhEsWGalmE9hSF7zK/d0+WIvN1p2l7lKTHeEdPuGKHLhBxl6w0
gXVAP+7isA8sAlTIzbC4fXV9TvUBGR8sXTddOjp5KVc3YY5KUrCsnO3DKHA1zfTQ
j7lV6GxUn4dXP5WkE3ht4W9gEljlcC1s0x+xcMlTOGnWhJY+ol4cuVjKnnOEShHD
CyPyqRxv9gQpXksl8W9pEG/eJk/pmY//xCYMvdpgdbN+RjB9Vtl5gQUnfUqoymoj
URhptXt48xaHxZkt0CM9OjsFjkJnLab3FB1TyzCZjkL7RWZmLqtqkkg3yuxm3d8p
1r+BsD7SLX1g0hH1Ln3ySg/esYPd3WQMuGB2NTtDc3rfxDwzhKw0w+UE0PNDJBNk
ujeC1Vs6ItWJ/hB/2qnzqZBqdddY1FwB+ziEjYoW914svBJYxwLk5HbXNV+CpxEh
-----END CERTIFICATE-----`
)

type Client struct {
  conn	  net.Conn
  domain  string
  enc	  *xml.Encoder
  dec	  *xml.Decoder
}

type Options struct {
  Host	    string
  Port	    string
  User	    string
  Password  string
  Mechanism string
}

type Cookie uint64

func getCookie() Cookie {
  var buf [8]byte
  if _, err := rand.Reader.Read(buf[:]); err != nil {
    panic("Failed to read random bytes: " + err.Error())
  }
  return Cookie(binary.LittleEndian.Uint64(buf[:]))
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
    client	= new(Client)
    err		error
    conn	net.Conn
    sf		= new(streamFeatures)
    authBase64	[]byte
    user	string
    //name	xml.Name
    //val		interface{}
  )

  if conn, err = connect(o); err != nil {
    return client, err
  }

  client.conn = conn
  client.enc = xml.NewEncoder(client.conn)
  client.dec = xml.NewDecoder(client.conn)
  client.domain = strings.Split(strings.Split(o.User, "@")[1], "/")[0]
  user = strings.Split(o.User, "@")[0]

  if sf, err = client.startStream(); err != nil {
    return client, err
  }

  if sf, err = client.startTLSStream(sf); err != nil {
    return client, err
  }

  for _, mechanism := range sf.Mechanism.Mechanism {
    switch mechanism {
    case "PLAIN":
      authBase64 = []byte(fmt.Sprintf("%s%s%s%s", BINARY_SALS, user, BINARY_SALS, o.Password))
      if _, err = client.conn.Write([]byte(fmt.Sprintf("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>",  base64.StdEncoding.EncodeToString(authBase64)))); err != nil {
	return client, err
      }
      //fmt.Fprintf(client.conn, fmt.Sprintf("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>", base64.StdEncoding.EncodeToString(authBase64)))
      break
    }
  }

  if _, _, err = next(client.dec); err != nil {
    return client, err
  }

  if _, err = client.startStream(); err != nil {
    return client, err
  }

  cookie := getCookie()

  fmt.Fprintf(client.conn, fmt.Sprintf("<iq type='set' id='%x'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>", cookie))

  var iq clientIQ

  if err = client.dec.DecodeElement(&iq, nil); err != nil {
    return client, err
  }

  fmt.Println(iq)
  return new(Client), nil
}

func (c *Client) startTLSStream(sf *streamFeatures) (*streamFeatures, error) {
  /*if sf.TLS == nil {
    return sf, nil
  }

  if sf.TLS.Required == nil {
    return sf, nil
  }*/

  var (
    err	      error
    configTLS tls.Config
    tlsconn   *tls.Conn
    proceed   tlsProceed
    certs     = x509.NewCertPool()
    ok	      bool
    //xtls      []byte
  )

  /*if xtls, err = utils.MarshalSelfClosingTag(startTLS{XMLNS: XML_TLS}); err != nil {
    return sf, err
  }*/

  fmt.Fprintf(c.conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>")
  /*if err = c.enc.Encode(&xtls); err != nil {
    return sf, err
  }*/

  if err = c.dec.DecodeElement(&proceed, nil); err != nil {
    return sf, err
  }

  if ok = certs.AppendCertsFromPEM([]byte(rootPEM)); !ok {
    return sf, errors.New("Failed to parse root certificate")
  }

  configTLS.ServerName = "localhost"
  configTLS.InsecureSkipVerify = false
  configTLS.RootCAs = certs
  tlsconn = tls.Client(c.conn, &configTLS)

  if err = tlsconn.Handshake(); err != nil {
    fmt.Println("DEU BOSTA AQUI")
    return sf, err
  }

  if err = tlsconn.VerifyHostname("localhost"); err != nil {
    fmt.Println("DEU MERDA AQUI")
    return sf, err
  }

  c.conn = tlsconn
  fmt.Println("PORRA")

  return c.startStream()
  //fmt.Fprintf(c.conn, "<stream:stream to='localhost' xmlns:stream='http://etherx.jabber.org/streams' xmlns='jabber:client' xml:lang='en' version='1.0'>")

  //return sf, err
}

func (c *Client) startStream() (*streamFeatures, error) {
  var (
    stream  []byte
    err	    error
    sf	    = new(streamFeatures)
    se	    xml.StartElement
  )

  c.dec = xml.NewDecoder(c.conn)
  c.enc = xml.NewEncoder(c.conn)

  if stream, err = utils.MarshalWithOutEndTag(Stream{XMLNSStream: XML_STREAM, XMLNS: XML_CLIENT, To: "localhost", Language: "en", Version: VERSION}, true); err != nil {
    return sf, err
  }

  if _, err = c.conn.Write(stream); err != nil {
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

func next(dec *xml.Decoder) (xml.Name, interface{}, error) {
  var (
    se	xml.StartElement
    err	error
    nv	interface{}
  )

  if se, err = startStream(dec); err != nil {
    return se.Name, nv, err
  }

  switch fmt.Sprintf("%s %s", se.Name.Space, se.Name.Local) {
  case fmt.Sprintf("%s success", nsSASL):
    nv = &saslSuccess{}
  }

  if err = dec.DecodeElement(nv, &se); err != nil {
    return se.Name, nv, err
  }

  return se.Name, nv, nil
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
