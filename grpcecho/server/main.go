package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"

	pb "github.com/tcolgate/test/grpcecho/ping"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
)

type pingServer struct{}

func (p *pingServer) Ping(ctx context.Context, pr *pb.PingRequest) (*pb.PingReply, error) {
	t := time.Now()
	r := &pb.PingReply{t.UnixNano()}

	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("some event happened")
	}
	//tr := trace.New("blah", "blahping")
	//defer tr.Finish()
	//tr.LazyPrintf("some event happened")

	return r, nil
}

func main() {
	grpc.EnableTracing = true

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4567))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterPingServer(grpcServer, &pingServer{})

	go grpcServer.Serve(lis)

	time.Sleep(time.Second * 1)

	conn, err := grpc.Dial("localhost:4567", grpc.WithBlock())
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewPingClient(conn)

	ctx := context.Background()
	tr := trace.New("mything", "otherthing")

	ctx = trace.NewContext(ctx, tr)

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

	tr.Finish()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
