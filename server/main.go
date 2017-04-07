package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
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

type Shard struct {
	Size              int64
	Status            int64
	StartTimestamp    int64
	CompleteTimestamp int64
	Worker            string
}

type Job struct {
	Size     int64
	Status   bool
	Shards   [1000]*Shard
	Code     string
	Lock     *sync.Mutex
	Solution string // I'd like this to be more generic, but oh well
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
	shard := &pb.WorkUnit{}
	var i int
	// Find next workUnit in workUnits with status = Unworked (0) or Timestamp is too old
	// TODO figure out how this could end up unbounded or something?
	for i = range job.Shards {
		if job.Shards[i] == nil {
			job.Shards[i] = &Shard{}
		}
		if job.Shards[i].Status == 0 {
			break
		}
		// TODO hand out shards that are stale
		// TODO maybe find the oldest shard not submitted?
	}
	log.Printf("Distributing work unit %d\n", i)
	job.Shards[i].Worker = in.UUID
	job.Shards[i].StartTimestamp = time.Now().Unix()
	job.Shards[i].Status = 1
	shard.JobID = jobid
	shard.ID = int64(i)
	shard.Offset = int64(i)
	shard.Size = job.Size
	shard.Code = job.Code
	// TODO something about shard log/history
	return shard, nil
}

func (s *server) ReceiveResults(ctx context.Context, r *pb.Results) (*pb.Success, error) {
	job := jobs[r.JobID]
	job.Lock.Lock()
	defer job.Lock.Unlock()
	shard := job.Shards[r.WorkID]
	shard.Status = 2
	shard.CompleteTimestamp = time.Now().Unix()
	if r.Found {
		log.Println("Someone found it!")
		log.Printf("%#v\n", r)
		job.Solution = strconv.FormatInt(r.Location, 10)
		job.Status = false
	}
	return &pb.Success{
		Success: true,
	}, nil
}

func (s *server) SubmitJob(ctx context.Context, job *pb.Job) (*pb.Success, error) {
	if job.Name == "" {
		return &pb.Success{
			Success: false,
		}, nil
	}
	if job.Code == "" {
		return &pb.Success{
			Success: false,
		}, nil
	}
	log.Printf("Got a job submission: %s\n", job.Name)
	jobs[job.Name] = &Job{
		Size:   1000000,
		Status: true,
		Lock:   &sync.Mutex{},
		Code:   job.Code,
	}
	return &pb.Success{
		Success: true,
	}, nil
}

func (s *server) GetAllJobs(ctx context.Context, _ *pb.JobParams) (*pb.JobStatuses, error) {
	var statuses []*pb.JobStatuses_JobStatus

	for name, job := range jobs {
		deliveredShards := int64(0)
		completedShards := int64(0)
		for _, shard := range job.Shards {
			if shard == nil {
				continue
			}
			if shard.Status == 1 {
				deliveredShards++
			} else if shard.Status == 2 {
				completedShards++
			}
		}
		statuses = append(statuses, &pb.JobStatuses_JobStatus{
			Name:            name,
			DeliveredShards: deliveredShards,
			CompletedShards: completedShards,
			TotalShards:     int64(len(job.Shards)),
			Found:           !job.Status,
			Location:        job.Solution,
		})
	}
	return &pb.JobStatuses{
		Statuses: statuses,
	}, nil
}

func Main(cmd *cobra.Command, args []string) {
	jobs = make(map[string]*Job)

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
