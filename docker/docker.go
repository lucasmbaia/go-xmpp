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

type Image struct {
  To	string
  From	string
  Path	string
  Name	string
  Key	string
}

type Deploy struct {
  To		  string
  From		  string
  Customer	  string
  ApplicationName string
  TotalContainers string
  Cpus		  string
  Memory	  string
  Ports		  string
  Path		  string
}

type Container struct {
  To	string
  From	string
}

func request(from, to string, elements []Element) (IQ, error) {
  var (
    err      error
    id       uuid.UUID
  )

  if id, err = uuid.NewV4(); err != nil {
    return IQ{}, err
  }

  return IQ{
    From: from,
    To:   to,
    Type: "set",
    ID:   id.String(),
    Query: QueryDocker{
      Elements: elements,
    },
  }, nil
}

func ActionContainer(action Action) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "container"}, Value: action.Container},
    {XMLName: xml.Name{Local: "action"}, Value: action.Action},
  }

  return request(action.From, action.To, elements)
}

func GenerateImage(image Image) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "action"}, Value: "generate-image"},
    {XMLName: xml.Name{Local: "path"}, Value: image.Path},
    {XMLName: xml.Name{Local: "name"}, Value: image.Name},
    {XMLName: xml.Name{Local: "key"}, Value: image.Key},
  }

  return request(image.From, image.To, elements)
}

func LoadImage(image Image) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "action"}, Value: "load-image"},
    {XMLName: xml.Name{Local: "path"}, Value: image.Path},
  }

  return request(image.From, image.To, elements)
}

func MasterDeploy(deploy Deploy) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "action"}, Value: "master-deploy"},
    {XMLName: xml.Name{Local: "customer"}, Value: deploy.Customer},
    {XMLName: xml.Name{Local: "application-name"}, Value: deploy.ApplicationName},
    {XMLName: xml.Name{Local: "total-containers"}, Value: deploy.TotalContainers},
    {XMLName: xml.Name{Local: "cpus"}, Value: deploy.Cpus},
    {XMLName: xml.Name{Local: "memory"}, Value: deploy.Memory},
    {XMLName: xml.Name{Local: "ports"}, Value: deploy.Ports},
    {XMLName: xml.Name{Local: "path"}, Value: deploy.path},
  }

  return request(deploy.From, deploy.To, elements)
}

func NameContainers(c Container) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "action"}, Value: "name-containers"},
  }

  return request(c.From, c.To, elements)
}

func TotalContainers(c Container) (IQ, error) {
  var elements = []Element{
    {XMLName: xml.Name{Local: "action"}, Value: "total-containers"},
  }

  return request(c.From, c.To, elements)
}
