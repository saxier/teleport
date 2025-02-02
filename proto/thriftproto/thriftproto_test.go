package thriftproto_test

import (
	"testing"
	"time"

	"github.com/andeya/erpc/v7/proto/thriftproto"
	"github.com/andeya/erpc/v7/xfer/gzip"
)

type Home struct {
	erpc.CallCtx
}

func (h *Home) Test(arg *Test) (*Test, *erpc.Status) {
	if string(h.PeekMeta("peer_id")) != "110" {
		panic("except meta: peer_id=110")
	}
	return &Test{
		Author: arg.Author + "->OK",
	}, nil
}

func TestBinaryProto(t *testing.T) {
	gzip.Reg('g', "gizp-5", 5)

	// server
	srv := erpc.NewPeer(erpc.PeerConfig{ListenPort: 9090, DefaultBodyCodec: "thrift"})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(thriftproto.NewBinaryProtoFunc())
	defer srv.Close()
	time.Sleep(1e9)

	// client
	cli := erpc.NewPeer(erpc.PeerConfig{DefaultBodyCodec: "thrift"})
	sess, stat := cli.Dial(":9090", thriftproto.NewBinaryProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result Test
	stat = sess.Call("Home.Test",
		&Test{Author: "andeya"},
		&result,
		erpc.WithAddMeta("peer_id", "110"),
		erpc.WithXferPipe('g'),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	if result.Author != "andeya->OK" {
		t.FailNow()
	}
	t.Logf("result:%v", result)
}

func TestStructProto(t *testing.T) {
	// server
	srv := erpc.NewPeer(erpc.PeerConfig{ListenPort: 9090})
	srv.RouteCall(new(Home))
	go srv.ListenAndServe(thriftproto.NewStructProtoFunc())
	defer srv.Close()
	time.Sleep(1e9)

	// client
	cli := erpc.NewPeer(erpc.PeerConfig{})
	sess, stat := cli.Dial(":9090", thriftproto.NewStructProtoFunc())
	if !stat.OK() {
		t.Fatal(stat)
	}
	var result Test
	stat = sess.Call("Home.Test",
		&Test{Author: "andeya"},
		&result,
		erpc.WithAddMeta("peer_id", "110"),
	).Status()
	if !stat.OK() {
		t.Error(stat)
	}
	if result.Author != "andeya->OK" {
		t.FailNow()
	}
	t.Logf("result:%v", result)
}
