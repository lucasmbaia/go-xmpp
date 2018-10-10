package xmpp

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_NewClient(t *testing.T) {
	var ctx = context.Background()

	var options = Options{
		Host:      "172.16.95.179",
		Port:      "5222",
		Mechanism: PLAIN,
		User:      "zeus@localhost",
		Password:  "totvs@123",
	}

	conn, _ := NewClient(options)
	go conn.Receiver(ctx, RecvOptions{})
	time.Sleep(2 * time.Second)
	fmt.Println(conn.Roster())
	conn.DiscoItems("conference.localhost")
	conn.DiscoItems("minions@conference.localhost")
	conn.MucPresence()
	//conn.CreateRoom("chups@conference.localhost")

	<-ctx.Done()
}
