package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/andeya/erpc/v7"
	"github.com/andeya/erpc/v7/examples/bench/msg"
)

//go:generate go build $GOFILE

var (
	port      = flag.Int64("p", 8972, "listened port")
	delay     = flag.Duration("delay", 0, "delay to mock business processing")
	debugAddr = flag.String("d", "127.0.0.1:9981", "server ip and port")
	network   = flag.String("network", "tcp", "network")
)

func main() {
	flag.Parse()

	defer erpc.SetLoggerLevel("ERROR")()
	erpc.SetGopool(1024*1024*100, time.Minute*10)

	go func() {
		log.Println(http.ListenAndServe(*debugAddr, nil))
	}()

	erpc.SetServiceMethodMapper(erpc.RPCServiceMethodMapper)
	server := erpc.NewPeer(erpc.PeerConfig{
		Network:          *network,
		DefaultBodyCodec: "protobuf",
		ListenPort:       uint16(*port),
	})
	server.RouteCall(new(Hello))
	server.ListenAndServe()
}

type Hello struct {
	erpc.CallCtx
}

func (t *Hello) Say(args *msg.BenchmarkMessage) (*msg.BenchmarkMessage, *erpc.Status) {
	s := "OK"
	var i int32 = 100
	args.Field1 = s
	args.Field2 = i
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	return args, nil
}
