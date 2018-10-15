package xmpp

import (
  "encoding/xml"
)

type Stream struct {
  XMLName     xml.Name `xml:"stream:stream"`
  XMLNS       string   `xml:"xmlns,attr,omitempty"`
  XMLNSStream string   `xml:"xmlns:stream,attr,omitempty"`
  Language    string   `xml:"xml:lang,attr,omitempty"`
  To          string   `xml:"to,attr,omitempty"`
  From        string   `xml:"from,attr,omitempty"`
  ID          string   `xml:"id,attr,omitempty"`
  Version     string   `xml:"version,attr,omitempty"`
}

type streamFeatures struct {
  XMLName   xml.Name `xml:"features"`
  TLS       *startTLS
  Mechanism mechanism
}

type mechanism struct {
  XMLName   xml.Name `xml:"mechanisms,omitempty"`
  Mechanism []string `xml:"mechanism,omitempty"`
}

type startTLS struct {
  XMLName  xml.Name `xml:"starttls,omitempty"`
  XMLNS    string   `xml:"xmlns,attr,omitempty"`
  Required *string  `xml:"required,omitempty"`
}

type tlsProceed struct {
  XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type authPLAIN struct {
  XMLName   xml.Name `xml:"auth"`
  XMLNS     string   `xml:"xmlns,attr,omitempty"`
  Mechanism string   `xml:"mechanism,attr,omitempty"`
  XMLNSGA   string   `xml:"xmlns:ga,attr,omitempty"`
  GA        bool     `xml:"ga:client-uses-full-bind-result,attr,omitempty"`
}

type saslSuccess struct {
  XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type saslFailure struct {
  XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
  Any     xml.Name `xml:",any"`
  Text    string   `xml:"text"`
}

type ClientIQ struct {
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

type Message struct {
  XMLName xml.Name `xml:"message"`
  XMLNSGA string   `xml:"xmlns:ga,attr,omitempty"`
  To      string   `xml:"to,attr,omitempty"`
  From    string   `xml:"from,attr,omitempty"`
  Type    string   `xml:"type,attr,omitempty"`
  ID      string   `xml:"id,attr,omitempty"`
  Subject string   `xml:"subject,omitempty"`
  Body    string   `xml:"body,omitempty"`
}

type clientPresence struct {
  XMLName xml.Name `xml:"jabber:client presence"`
  From    string   `xml:"from,attr"`
  ID      string   `xml:"id,attr"`
  To      string   `xml:"to,attr"`
  Type    string   `xml:"type,attr"` // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
  Lang    string   `xml:"lang,attr"`

  Show     string `xml:"show"`   // away, chat, dnd, xa
  Status   string `xml:"status"` // sb []clientText
  Priority string `xml:"priority,attr"`
  Error    *clientError
  User	  MucUser
}

type clientQuery struct {
  Item []RosterItem
}

/*type rosterItem struct {
  XMLName      xml.Name `xml:"jabber:iq:roster item"`
  Jid          string   `xml:",attr"`
  Name         string   `xml:",attr"`
  Subscription string   `xml:",attr"`
  Group        []string
}*/

/******************************************** NEW ********************************************/

type IQ struct {
  XMLName xml.Name `xml:"iq"`
  From    string   `xml:"from,attr,omitempty"`
  To      string   `xml:"to,attr,omitempty"`
  Type    string   `xml:"type,attr,omitempty"`
  ID      string   `xml:"id,attr,omitempty"`
  Query   []byte   `xml:"query,omitempty"`
}

type QueryRoster struct {
  XMLName xml.Name      `xml:"jabber:iq:roster query"`
  Item    []RosterItem  `xml:"item"`
}

type QueryDiscoItems struct {
  XMLName xml.Name      `xml:"http://jabber.org/protocol/disco#items query"`
  Item    []RosterItem  `xml:"item"`
}

type RosterItem struct {
  Jid		string    `xml:"jid,attr,omitempty"`
  Subscription	string	  `xml:"subscription,attr,omitempty"`
  Name		string	  `xml:"name,attr,omitempty"`
  Group		[]string  `xml:"group,omitempty"`
}

type Query struct {
  XMLName xml.Name
}

type Presence struct {
  XMLName   xml.Name `xml:"jabber:client presence"`
  From	    string   `xml:"from,attr,omitempty"`
  ID	    string   `xml:"id,attr,omitempty"`
  To	    string   `xml:"to,attr,omitempty"`
  Type	    string   `xml:"type,attr,omitempty"` // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
  Lang	    string   `xml:"lang,attr,omitempty"`

  Show	    string    `xml:"show,omitempty"`   // away, chat, dnd, xa
  Status    string    `xml:"status,omitempty"` // sb []clientText
  Priority  string    `xml:"priority,attr,omitempty"`
  Error	    *clientError
  User	    MucUser
}

type MucUser struct {
  XMLName xml.Name	  `xml:"http://jabber.org/protocol/muc#user x"`
  Item	  []itemPresence  `xml:"item"`
}

type itemPresence struct {
  Jid	      string  `xml:"jid,attr,omitempty"`
  Role	      string  `xml:"role,attr,omitempty"`
  Affiliation string  `xml:"affiliation,attr,omitempty"`
}

