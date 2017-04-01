package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/brimstone/dwork/pb"
	"github.com/spf13/cobra"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

var port = ":9000"
var working = true
var workMutex sync.Mutex
var workUnits []pb.WorkUnit

func (s *server) GiveWork(ctx context.Context, in *pb.WorkerID) (*pb.WorkUnit, error) {
	if !working {
		return nil, fmt.Errorf("No work")
	}
	workMutex.Lock()
	defer workMutex.Unlock()
	// TODO Find next workUnit in workUnits with status = Unworked
	for i := range workUnits {
		if workUnits[i].Status != 1 {
			continue
		}
		log.Printf("Distributing work unit %d\n", i)
		workUnits[i].Id = int64(i)
		workUnits[i].Status = 2
		workUnits[i].Timestamp = time.Now().Unix()
		return &workUnits[i], nil
	}
	return nil, fmt.Errorf("No work")
}

// TODO
func (s *server) ReceiveResults(ctx context.Context, r *pb.Results) (*pb.ResultsSuccess, error) {
	if r.Found {
		log.Println("Someone found it!")
		working = false
	}
	return &pb.ResultsSuccess{
		Success: true,
	}, nil
}

func Main(cmd *cobra.Command, args []string) {
	workSize := 10
	for i := int64(0); i < 100; i++ {
		workUnits = append(workUnits, pb.WorkUnit{
			Offset: i,
			Size:   int64(workSize),
			Search: 42,
			Status: 1,
			Code: `
func work(x) {
	return sprintf("%x\n", sha256(x + ""))
}
func check(x) {
	if x[0:1] == "2" {
		return true
	}
	return false
}
`,
		})
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDworkServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("Ready to go")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
