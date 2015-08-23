package main

import (
	"flag"
	"log"
	"math/rand"
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
var delay = flag.Int64("delay", 0, "Time, in ms, to delay after downstreams have been pinged")
var normDelay = flag.Bool("normDelay", false, "The delay will be given a bit of variation")
var justOne = flag.Bool("justOne", false, "One response from a downstream is enough, cancel the rest")

type pingServer struct {
	downstreams []pb.PingClient
}

// pingDownstream isn't great, we create a new client every time,
// but we kinda want this to suck a but to make things more interesting
// to perf trace
func pingOneDownstream(ctx context.Context, d pb.PingClient, resp chan<- error) {
	done := make(chan error, 1)

	go func() {
		defer close(done)
		t := time.Now()
		req := pb.PingRequest{t.UnixNano()}
		_, err := d.Ping(ctx, &req)

		done <- err
	}()

	select {
	case err := <-done:
		resp <- err
	case <-ctx.Done():
		resp <- ctx.Err()
	}
}

// pingAllDownStreams waits for all the downstream pings to finish, then
// returns
func (p *pingServer) pingAllDownStreams(ctx context.Context) {
	resp := make(chan error)
	wg := sync.WaitGroup{}
	defer close(resp)

	ctx, cancel := context.WithCancel(ctx)

	for _, d := range p.downstreams {
		wg.Add(1)
		go func(d pb.PingClient) {
			pingOneDownstream(ctx, d, resp)
			wg.Done()
		}(d)
	}

	maxResp := len(p.downstreams)

	for i := 0; i < maxResp; i++ {
		_ = <-resp
		if *justOne {
			cancel()
		}
	}

	wg.Wait()
}

func (p *pingServer) Ping(ctx context.Context, pr *pb.PingRequest) (*pb.PingReply, error) {
	d := *delay * 1000
	if *normDelay {
		d = int64(rand.NormFloat64()*1000) + d
	}

	select {
	case <-time.After(time.Duration(d * int64(time.Microsecond))):
		p.pingAllDownStreams(ctx)
	case <-ctx.Done():
	}

	return &pb.PingReply{time.Now().UnixNano()}, nil
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

	var downstreamClients []pb.PingClient
	for _, d := range strings.Split(*downstreams, ",") {
		if d == "" {
			continue
		}
		log.Println(d)
		conn, err := grpc.Dial(d, grpc.WithBlock())
		if err != nil {
			log.Fatalf("In here " + err.Error())
		}
		defer conn.Close()

		downstreamClients = append(downstreamClients, pb.NewPingClient(conn))
	}
	pb.RegisterPingServer(grpcServer, &pingServer{downstreams: downstreamClients})

	go grpcServer.Serve(lis)

	log.Fatal(http.ListenAndServe(*httpPort, nil))
}
