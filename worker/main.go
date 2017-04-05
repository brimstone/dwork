package worker

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/brimstone/dwork/pb"
	"github.com/mattn/anko/vm"
	"github.com/spf13/cobra"
)

var c pb.DworkClient

func Main(cmd *cobra.Command, args []string) {
	address, err := cmd.Flags().GetString("server")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewDworkClient(conn)

	done := make(chan bool)
	syscall.Setpriority(syscall.PRIO_PROCESS, 0, 19)
	c := runtime.NumCPU()
	for i := 0; i < c; i++ {
		go createWorker(done)
	}

	<-done
}

func performWork(env *vm.Env, w *pb.WorkUnit) *pb.Results {
	r := &pb.Results{}
	r.WorkID = w.ID
	r.JobID = w.JobID
	log.Printf("Performing work on job %s shard %d\n", w.JobID, w.ID)
	// real work happens here
	for i := w.Offset * w.Size; i < (w.Offset+1)*w.Size; i++ {
		v, err := env.Execute("work(" + strconv.FormatInt(i, 10) + ")")
		if err != nil {
			log.Fatal("Error iterating ", err)
		}
		if v.Type().Kind() != reflect.Bool || v.Bool() != false {
			fmt.Printf("Found it: %d %s\n", i, v)
			r.Location = i
			r.Found = true
			break
		}
	}
	return r
}

func createWorker(done chan bool) {
	waitBackoff := 1
	for {
		workUnit, err := c.GiveWork(context.Background(), &pb.WorkerID{})
		if err != nil {
			fmt.Printf("Error retrieving work: %#v Waiting for %ds for more work\n", err.Error(), waitBackoff)
			time.Sleep(time.Second * time.Duration(waitBackoff))
			waitBackoff = waitBackoff * 2
			if waitBackoff > 64 {
				waitBackoff = 64
			}
			continue
		}
		waitBackoff = 1
		env, err := createVm(workUnit.Code)
		if err != nil {
			log.Println("Error loading usercode into worker vm")
			continue
		}
		results := performWork(env, workUnit)
		success, err := c.ReceiveResults(context.Background(), results)
		if !success.Success {
			log.Println("why fail?")
		}
	}
	// Unreachable code
	//done <- true
}

func createVm(usercode string) (*vm.Env, error) {
	env := vm.NewEnv()

	env.Define("printf", fmt.Printf)
	env.Define("sprintf", fmt.Sprintf)
	env.Define("sha256", sha256.Sum256)

	_, err := env.Execute(`
	func work(x) {
		return true
	}
`)
	if err != nil {
		return nil, fmt.Errorf("Error loading base code: %s", err)
	}

	// load user code
	_, err = env.Execute(usercode)
	if err != nil {
		return nil, fmt.Errorf("Error loading user code %s", err)
	}

	return env, nil
}
