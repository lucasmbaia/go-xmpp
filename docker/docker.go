package docker

import (
	"encoding/xml"
	"github.com/satori/go.uuid"
)

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	From    string   `xml:"from,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	ID      string   `xml:"id,attr,omitempty"`
	Query   interface{}
}

type QueryDocker struct {
	XMLName  xml.Name `xml:"jabber:iq:docker query"`
	Elements []Element
}

type Element struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Action struct {
	To        string
	From      string
	Container string
	Action    string
}

func ActionContainer(action Action) (IQ, error) {
	var (
		err      error
		elements []Element
		id       uuid.UUID
	)

	if id, err = uuid.NewV4(); err != nil {
		return IQ{}, err
	}

	elements = []Element{
		{XMLName: xml.Name{Local: "container"}, Value: action.Container},
		{XMLName: xml.Name{Local: "action"}, Value: action.Action},
	}

	return IQ{
		From: action.From,
		To:   action.To,
		Type: "set",
		ID:   id.String(),
		Query: QueryDocker{
			Elements: elements,
		},
	}
}
