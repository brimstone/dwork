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
var jobs map[string]*Job

const MaxInt64 = int(^uint(0) >> 1)

type Job struct {
	Size     int64
	Status   bool
	Shards   [1000]*pb.WorkUnit
	Code     string
	Lock     *sync.Mutex
	Solution interface{}
}

func (s *server) GiveWork(ctx context.Context, in *pb.WorkerID) (*pb.WorkUnit, error) {
	var jobid string

	// TODO validate in

	// Look for a job that needs workers
	for id, job := range jobs {
		if job.Status {
			jobid = id
			break
		}
	}
	if jobid == "" {
		return nil, fmt.Errorf("No work")
	}
	job := jobs[jobid]
	job.Lock.Lock()
	defer job.Lock.Unlock()
	var shard *pb.WorkUnit
	var i int
	// Find next workUnit in workUnits with status = Unworked (0) or Timestamp is too old
	// TODO figure out how this could end up unbounded or something?
	for i = range job.Shards {
		shard = job.Shards[i]
		if shard == nil {
			job.Shards[i] = &pb.WorkUnit{}
			shard = job.Shards[i]
		}
		if shard.Status == 0 {
			break
		}
		// TODO hand out shards that are stale
	}
	log.Printf("Distributing work unit %d\n", i)
	shard.JobID = jobid
	shard.Id = int64(i)
	shard.Offset = int64(i)
	shard.Status = 1
	shard.Size = job.Size
	shard.Code = job.Code
	shard.Timestamp = time.Now().Unix()
	// TODO something about shard log/history
	return shard, nil
}

func (s *server) ReceiveResults(ctx context.Context, r *pb.Results) (*pb.ResultsSuccess, error) {
	if r.Found {
		log.Println("Someone found it!")
		log.Printf("%#v\n", r)
		job := jobs[r.JobID]
		job.Lock.Lock()
		defer job.Lock.Unlock()
		job.Solution = r.Location
		job.Status = false
	}
	return &pb.ResultsSuccess{
		Success: true,
	}, nil
}

func Main(cmd *cobra.Command, args []string) {
	jobs = make(map[string]*Job)

	jobs["sha256"] = &Job{
		Size:   1000000,
		Status: true,
		Lock:   &sync.Mutex{},
		Code: `
func work(x) {
	msg = x + ""
	hash = sprintf("%x\n", sha256(msg))
	if hash[0:6] == "ffffff" {
		return hash
	}
}
`,
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
