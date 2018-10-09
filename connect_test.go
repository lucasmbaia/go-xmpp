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
		Host:      "192.168.204.131",
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
	conn.CreateRoom("chups@conference.localhost")

	<-ctx.Done()
}
