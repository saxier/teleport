package main

import (
	"time"

	"github.com/andeya/erpc/v7"
)

//go:generate go build $GOFILE

func main() {
	defer erpc.SetLoggerLevel("DEBUG")()

	cli := erpc.NewPeer(erpc.PeerConfig{RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	cli.SetTLSConfig(erpc.GenerateTLSConfigForClient())

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}

	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		erpc.WithAddMeta("author", "andeya"),
	).Status()
	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}
	erpc.Printf("result: %d", result)
	erpc.Printf("Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

// Push push handler
type Push struct {
	erpc.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *erpc.Status {
	erpc.Printf("%s", *arg)
	return nil
}
