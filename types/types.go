package types

import (
  "encoding/xml"
)

type IQ struct {
  XMLName xml.Name    `xml:"iq"`
  From    string      `xml:"from,attr,omitempty"`
  To      string      `xml:"to,attr,omitempty"`
  Type    string      `xml:"type,attr,omitempty"`
  ID      string      `xml:"id,attr,omitempty"`
  Query   []byte `xml:",innerxml"`
  Error   *IQError
  Bind	  IQBind
}

type IQError struct {
  XMLName xml.Name  `xml:"jabber:client error"`
  Code    string    `xml:"code,attr,omitempty"`
  Type    string    `xml:"type,attr,omitempty"`
  Text    string    `xml:"text,omitempty"`
}

type IQBind struct {
  XMLName   xml.Name  `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
  Resource  string
  Jid	    string    `xml:"jid"`
}
