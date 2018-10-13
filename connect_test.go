package xmpp

import (
	"context"
	"fmt"
	"github.com/lucasmbaia/go-xmpp/docker"
	"testing"
	"time"
)

func Test_NewClient(t *testing.T) {
	var ctx = context.Background()

	var options = Options{
		Host:      "192.168.204.131",
		Port:      "5222",
		Mechanism: PLAIN,
		User:      "minion-1@localhost",
		Password:  "totvs@123",
	}

	conn, _ := NewClient(options)
	go conn.Receiver(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println(conn.Roster())
	conn.DiscoItems("conference.localhost")
	conn.DiscoItems("minions@conference.localhost")
	conn.MucPresence("minions@conference.localhost")
	//conn.CreateRoom("chups@conference.localhost")

	<-ctx.Done()
}

func Test_Docker(t *testing.T) {
	var ctx = context.Background()

	var options = Options{
		Host:      "192.168.204.131",
		Port:      "5222",
		Mechanism: PLAIN,
		User:      "zeus@localhost",
		Password:  "totvs@123",
	}

	message := func(i interface{}) {
		fmt.Println(i)
	}

	conn, _ := NewClient(options)
	go conn.Receiver(ctx)
	conn.RegisterHandler(MESSAGE_HANDLER, message)

	time.Sleep(2 * time.Second)
	fmt.Println(conn.Roster())
	conn.DiscoItems("conference.localhost")
	conn.DiscoItems("minions@conference.localhost")
	conn.MucPresence("minions@conference.localhost")
	time.Sleep(2 * time.Second)

	iq, err := docker.GenerateImage(docker.Image{
		From: conn.Jid,
		To:   "minion-1@localhost/145197768943085423711250",
		Path: "teste",
		Name: "teste",
		Key:  "teste",
	})

	fmt.Println(iq, err)

	if err == nil {
		fmt.Println(conn.Send(iq))
	}
	<-ctx.Done()
}
