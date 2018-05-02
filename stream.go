package xmpp

import (
  "encoding/xml"
)

type Stream struct {
  XMLName     xml.Name	`xml:"stream:stream"`
  XMLNS       string	`xml:"xmlns,attr,omitemtpy"`
  XMLNSStream string	`xml:"xmlns:stream,attr,omitempty"`
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
