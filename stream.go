package xmpp

import (
  "encoding/xml"
)

type Stream struct {
  XMLName     xml.Name	`xml:"stream:stream"`
  XMLNS       string	`xml:"xmlns,attr,omitemtpy"`
  XMLNSStream string	`xml:"xmlns:stream,attr,omitempty"`
  Language    string	`xml:"xml:lang,attr,omitempty"`
  To          string	`xml:"to,attr,omitempty"`
  From	      string	`xml:"from,attr,omitempty"`
  ID	      string	`xml:"id,attr,omitempty"`
  Version     string	`xml:"version,attr,omitempty"`
}

type streamFeatures struct {
  XMLName   xml.Name `xml:"features"`
  TLS	    *startTLS
  Mechanism mechanism
}

type mechanism struct {
  XMLName   xml.Name `xml:"mechanisms,omitempty"`
  Mechanism []string `xml:"mechanism,omitempty"`
}

type startTLS struct {
  XMLName   xml.Name  `xml:"starttls,omitempty"`
  XMLNS	    string    `xml:"xmlns,attr,omitempty"`
  Required  *string   `xml:"required,omitempty"`
}

type tlsProceed struct {
  XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type authPLAIN struct {
  XMLName   xml.Name  `xml:"auth"`
  XMLNS	    string    `xml:"xmlns,attr,omitempty"`
  Mechanism string    `xml:"mechanism,attr,omitempty"`
  XMLNSGA   string    `xml:"xmlns:ga,attr,omitempty"`
  GA	    bool      `xml:"ga:client-uses-full-bind-result,attr,omitempty"`
}

type saslSuccess struct {
  XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type clientIQ struct {
  XMLName xml.Name `xml:"jabber:client iq"`
  From    string   `xml:"from,attr"`
  ID      string   `xml:"id,attr"`
  To      string   `xml:"to,attr"`
  Type    string   `xml:"type,attr"` // error, get, result, set
  Query   []byte   `xml:",innerxml"`
  Error   clientError
  Bind    bindBind
}

type clientError struct {
  XMLName xml.Name `xml:"jabber:client error"`
  Code    string   `xml:",attr"`
  Type    string   `xml:",attr"`
  Any     xml.Name
  Text    string
}

type bindBind struct {
  XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
  Resource string
  Jid      string `xml:"jid"`
}
