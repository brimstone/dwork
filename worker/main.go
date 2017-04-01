package worker

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/brimstone/dwork/pb"
	"github.com/spf13/cobra"
)

var address = "localhost:9000"

var c pb.DworkClient

func Main(cmd *cobra.Command, args []string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewDworkClient(conn)

	done := make(chan bool)
	createWorker(done)

	<-done
}

func performWork(w *pb.WorkUnit) *pb.Results {
	r := &pb.Results{}
	r.Workid = w.Id
	log.Printf("Performing work with %#v\n", w)
	// TODO real work happens here
	time.Sleep(time.Second * 2)
	for i := w.Offset * w.Size; i < (w.Offset+1)*w.Size; i++ {
		log.Println("Checking", i)
		if i == w.Search {
			r.Found = true
			break
		}
	}
	return r
}

func createWorker(done chan bool) {
	for {
		workUnit, err := c.GiveWork(context.Background(), &pb.WorkerID{})
		if err != nil && err.Error() == "No work" {
			break
		}
		results := performWork(workUnit)
		success, err := c.ReceiveResults(context.Background(), results)
		if !success.Success {
			log.Println("why fail?")
		}
	}
	done <- true
}
