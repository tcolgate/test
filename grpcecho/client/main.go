package main

import _ "net/http/pprof"
import _ "runtime/pprof"

import (
	"log"
	"os"
	"time"

	pb "github.com/tcolgate/test/grpcecho/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	grpc.EnableTracing = false

	conn, err := grpc.Dial("localhost:4567", grpc.WithBlock())
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewPingClient(conn)

	ctx := context.Background()
	t := time.Now()
	req := pb.PingRequest{t.UnixNano()}
	resp, err := client.Ping(ctx, &req)
	if err != nil {
		os.Exit(1)
	}

	t2 := time.Now()

	then := time.Unix(0, resp.TimeBack)

	reqdelay := then.Sub(t)
	log.Println("req ", reqdelay)

	repdelay := t2.Sub(then)
	log.Println("rep ", repdelay)
}
