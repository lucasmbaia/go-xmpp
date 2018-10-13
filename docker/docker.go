package docker

import (
	"encoding/xml"
	"github.com/satori/go.uuid"
	"sync"
)

const (
	GENERATE_IMAGE       = "generate-image"
	LOAD_IMAGE           = "load-image"
	EXISTS_IMAGE         = "exists-image"
	MASTER_DEPLOY        = "master-deploy"
	APPEND_DEPLOY        = "append-deploy"
	NAME_CONTAINERS      = "name-containers"
	TOTAL_CONTAINERS     = "total-containers"
	OPERATION_CONTAINERS = "operation-containers"
	EMPTY_STR            = ""
)

var (
	odr   map[string]chan Response
	mutex = &sync.RWMutex{}
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

type Response struct {
	Error    error
	Elements []Element
}

type Action struct {
	To        string
	From      string
	Container string
	Action    string
}

type Image struct {
	To   string
	From string
	Path string
	Name string
	Key  string
}

type Deploy struct {
	To              string
	From            string
	Customer        string
	ApplicationName string
	TotalContainers string
	Cpus            string
	Memory          string
	Ports           string
	Path            string
}

type Container struct {
	To   string
	From string
}

func init() {
	odr = make(map[string]chan Response)
}

func RegisterOperationDocker(operation string) {
	mutex.Lock()
	if _, ok := odr[operation]; !ok {
		odr[operation] = make(chan Response)
	}
	mutex.Unlock()
}

func request(from, to string, elements []Element) (IQ, error) {
	var (
		err error
		id  uuid.UUID
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

func getAction(elements []Element) string {
	for _, element := range elements {
		if element.XMLName.Local == "action" {
			return element.Value
		}
	}

	return EMPTY_STR
}

func ResponseDocker(body <-chan []byte) {
	var (
		result  QueryDocker
		respose Response
		err     error
		action  string
	)

	go func() {
		for {
			select {
			case b := <-body:
				if err = xml.Unmarshal(b, &result); err != nil {
					return
				}

				action = getAction(result.Elements)
				respose = Response{Elements: result.Elements}

				mutex.Lock()
				switch action {
				case EMPTY_STR:
					break
				case GENERATE_IMAGE:
					odr[GENERATE_IMAGE] <- respose
				case LOAD_IMAGE:
					odr[LOAD_IMAGE] <- respose
				case EXISTS_IMAGE:
					odr[EXISTS_IMAGE] <- respose
				case MASTER_DEPLOY:
					odr[MASTER_DEPLOY] <- respose
				case APPEND_DEPLOY:
					odr[APPEND_DEPLOY] <- respose
				case NAME_CONTAINERS:
					odr[NAME_CONTAINERS] <- respose
				case TOTAL_CONTAINERS:
					odr[TOTAL_CONTAINERS] <- respose
				case OPERATION_CONTAINERS:
					odr[OPERATION_CONTAINERS] <- respose
				default:
					break
				}
				mutex.Unlock()
			}
		}
	}()
}

func ActionContainer(action Action) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: OPERATION_CONTAINERS},
		{XMLName: xml.Name{Local: "container"}, Value: action.Container},
		{XMLName: xml.Name{Local: "operation"}, Value: action.Action},
	}

	return request(action.From, action.To, elements)
}

func GenerateImage(image Image) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: GENERATE_IMAGE},
		{XMLName: xml.Name{Local: "path"}, Value: image.Path},
		{XMLName: xml.Name{Local: "name"}, Value: image.Name},
		{XMLName: xml.Name{Local: "key"}, Value: image.Key},
	}

	return request(image.From, image.To, elements)
}

func LoadImage(image Image) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: LOAD_IMAGE},
		{XMLName: xml.Name{Local: "path"}, Value: image.Path},
	}

	return request(image.From, image.To, elements)
}

func ExistsImage(image Image) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: EXISTS_IMAGE},
		{XMLName: xml.Name{Local: "name"}, Value: image.Name},
	}

	return request(image.From, image.To, elements)
}

func MasterDeploy(deploy Deploy) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: MASTER_DEPLOY},
		{XMLName: xml.Name{Local: "customer"}, Value: deploy.Customer},
		{XMLName: xml.Name{Local: "application-name"}, Value: deploy.ApplicationName},
		{XMLName: xml.Name{Local: "total-containers"}, Value: deploy.TotalContainers},
		{XMLName: xml.Name{Local: "cpus"}, Value: deploy.Cpus},
		{XMLName: xml.Name{Local: "memory"}, Value: deploy.Memory},
		{XMLName: xml.Name{Local: "ports"}, Value: deploy.Ports},
		{XMLName: xml.Name{Local: "path"}, Value: deploy.Path},
	}

	return request(deploy.From, deploy.To, elements)
}

func AppendDeploy(deploy Deploy) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: APPEND_DEPLOY},
		{XMLName: xml.Name{Local: "customer"}, Value: deploy.Customer},
		{XMLName: xml.Name{Local: "application-name"}, Value: deploy.ApplicationName},
		{XMLName: xml.Name{Local: "total-containers"}, Value: deploy.TotalContainers},
	}

	return request(deploy.From, deploy.To, elements)
}

func NameContainers(c Container) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: NAME_CONTAINERS},
	}

	return request(c.From, c.To, elements)
}

func TotalContainers(c Container) (IQ, error) {
	var elements = []Element{
		{XMLName: xml.Name{Local: "action"}, Value: TOTAL_CONTAINERS},
	}

	return request(c.From, c.To, elements)
}
