package main

import (
	"flag"
	"log"
	"net"
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

var port = flag.String("port", ":4567", "the port to run grpc on")
var httpPort = flag.String("httpPort", ":8080", "the port to run the debug stuff on")
var downstreams = flag.String("downstream", "", "Ping these downstream servers")

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
	flag.Parse()

	grpc.EnableTracing = true

	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterPingServer(grpcServer, &pingServer{})

	go grpcServer.Serve(lis)

	log.Fatal(http.ListenAndServe(*httpPort, nil))
}
