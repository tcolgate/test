package main

import (
	"flag"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"

	pb "github.com/tcolgate/test/grpcecho/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var port = flag.String("port", ":4567", "the port to run grpc on")
var httpPort = flag.String("httpPort", ":8080", "the port to run the debug stuff on")
var downstreams = flag.String("downstream", "", "Ping these downstream servers")

type pingServer struct {
	downstreams []string
}

// pingDownstream isn't great, we create a new client every time,
// but we kinda want this to suck a but to make things more interesting
// to perf trace
func pingDownstream(d string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	conn, err := grpc.Dial(d, grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPingClient(conn)

	ctx := context.Background()

	t := time.Now()
	req := pb.PingRequest{t.UnixNano()}
	_, err = client.Ping(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

func (p *pingServer) pingDownStreams() {
	var wg sync.WaitGroup // waitgroup == code smell ?
	for _, d := range p.downstreams {
		wg.Add(1)
		go pingDownstream(d, &wg)
	}
	wg.Wait()
}

func (p *pingServer) Ping(ctx context.Context, pr *pb.PingRequest) (*pb.PingReply, error) {
	t := time.Now()
	r := &pb.PingReply{t.UnixNano()}

	p.pingDownStreams()

	//if tr, ok := trace.FromContext(ctx); ok {
	//		tr.LazyPrintf("some event happened")
	//	}
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

	pb.RegisterPingServer(grpcServer, &pingServer{downstreams: strings.Split(*downstreams, ",")})

	go grpcServer.Serve(lis)

	log.Fatal(http.ListenAndServe(*httpPort, nil))
}
