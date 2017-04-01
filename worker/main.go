package worker

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/brimstone/dwork/pb"
	"github.com/mattn/anko/vm"
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

func performWork(env *vm.Env, w *pb.WorkUnit) *pb.Results {
	r := &pb.Results{}
	r.Workid = w.Id
	log.Printf("Performing work with %#v\n", w)
	// TODO real work happens here
	time.Sleep(time.Second * 2)
	for i := w.Offset * w.Size; i < (w.Offset+1)*w.Size; i++ {
		log.Println("Checking", i)
		v, err := env.Execute("check(work(iterate(" + strconv.FormatInt(i, 10) + ")))")
		if err != nil {
			log.Fatal("Error iterating ", err)
		}
		if v.Type().Kind() == reflect.Bool && v.Bool() == true {
			fmt.Printf("Found it: %d\n", i)
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
	done <- true
}

func createVm(usercode string) (*vm.Env, error) {
	env := vm.NewEnv()

	env.Define("printf", fmt.Printf)
	env.Define("sprintf", fmt.Sprintf)
	env.Define("sha256", sha256.Sum256)

	_, err := env.Execute(`
	func iterate(x) {
		return x
	}

	func work(x) {
		return x
	}

	func check(x) {
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
