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
	REMOVE_CONTAINER     = "remove-container"
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
	Error   *IQError
	Query   interface{}
}

type IQError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:"code,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Text    string   `xml:"text,omitempty"`
}

type QueryDocker struct {
	XMLName xml.Name `xml:"jabber:iq:docker query"`
	Action  string   `xml:"action,attr,omitempty"`
	Elements
}

type Elements struct {
	ID              string           `xml:"id,omitempty"`
	Name            string           `xml:"name,omitempty"`
	Operation       string           `xml:"operation,omitempty"`
	Path            string           `xml:"path,omitempty"`
	Key             string           `xml:"key,omitempty"`
	Customer        string           `xml:"customer,omitempty"`
	ApplicationName string           `xml:"application-name,omitempty"`
	TotalContainers int              `xml:"total-containers,omitempty"`
	Cpus            string           `xml:"cpus,omitempty"`
	Memory          string           `xml:"memory,omitempty"`
	BuildName       string           `xml:"build-name,omitempty"`
	Tag             string           `xml:"tag,omitempty"`
	Image           string           `xml:"image,omitempty"`
	Address         string           `xml:"address,omitempty"`
	CreateImage     bool             `xml:"create-image,omitempty"`
	Minion          string           `xml:"minion,omitempty"`
	PortsContainer  []PortsContainer `xml:"portsContainer,omitempty"`
	Ports           []Ports          `xml:"ports,omitempty"`
	Containers      []Container      `xml:"containers,omitempty"`
	Args            []Args           `xml:"arg,omitempty"`
}

type Args struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Ports struct {
	Port     int    `xml:"port,omitempty"`
	Protocol string `xml:"protocol,omitempty"`
}

type PortsContainer struct {
	Source       string   `xml:"source,omitempty"`
	Destinations []string `xml:"destinations,omitempty"`
}

type Container struct {
	Name string `xml:"name,omitempty"`
}

type Element struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Response struct {
	Error    error
	Elements Elements
}

type Action struct {
	To        string
	From      string
	Container string
	Action    string
}

type Image struct {
	To        string
	From      string
	Name      string
	BuildName string
	Tag       string
	Path      string
}

type Deploy struct {
	To              string
	From            string
	Customer        string
	ApplicationName string
	Name            string
	Cpus            string
	Memory          string
	Ports           []Ports
	Args            []Args
	CreateImage     bool
	Image           string
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

func request(from, to, action string, elements Elements) (IQ, error) {
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
			Action:   action,
			Elements: elements,
		},
	}, nil
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

				action = result.Action
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
	var elements = Elements{
		Name:      action.Container,
		Operation: action.Action,
	}

	return request(action.From, action.To, OPERATION_CONTAINERS, elements)
}

func RemoveContainer(action Action) (IQ, error) {
	var elements = Elements{
		Name: action.Container,
	}

	return request(action.From, action.To, REMOVE_CONTAINER, elements)
}

func GenerateImage(image Image) (IQ, error) {
	var elements = Elements{
		Name:      image.Name,
		BuildName: image.BuildName,
		Tag:       image.Tag,
		Path:      image.Path,
	}

	return request(image.From, image.To, GENERATE_IMAGE, elements)
}

func LoadImage(image Image) (IQ, error) {
	var elements = Elements{
		Path: image.Path,
		Name: image.Name,
	}

	return request(image.From, image.To, LOAD_IMAGE, elements)
}

func ExistsImage(image Image) (IQ, error) {
	var elements = Elements{
		Name: image.Name,
	}

	return request(image.From, image.To, EXISTS_IMAGE, elements)
}

func MasterDeploy(deploy Deploy) (IQ, error) {
	var elements = Elements{
		Customer:        deploy.Customer,
		ApplicationName: deploy.ApplicationName,
		Name:            deploy.Name,
		Cpus:            deploy.Cpus,
		Memory:          deploy.Memory,
		Ports:           deploy.Ports,
		Args:            deploy.Args,
		Image:           deploy.Image,
		CreateImage:     deploy.CreateImage,
	}

	return request(deploy.From, deploy.To, MASTER_DEPLOY, elements)
}

func AppendDeploy(deploy Deploy) (IQ, error) {
	var elements = Elements{
		Customer:        deploy.Customer,
		ApplicationName: deploy.ApplicationName,
	}

	return request(deploy.From, deploy.To, APPEND_DEPLOY, elements)
}

func NameContainers(from, to string) (IQ, error) {
	return request(from, to, NAME_CONTAINERS, Elements{})
}

func TotalContainers(from, to string) (IQ, error) {
	return request(from, to, TOTAL_CONTAINERS, Elements{})
}
