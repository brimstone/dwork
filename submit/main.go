package submit

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/brimstone/dwork/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

	name, err := cmd.Flags().GetString("name")
	// TODO check that name is set
	filename, err := cmd.Flags().GetString("file")
	// TODO check that file is set
	code, err := ioutil.ReadFile(filename)
	// TODO check that there's stuff here

	job := &pb.Job{
		Name: name,
		Code: string(code),
	}
	success, err := c.SubmitJob(context.Background(), job)

	if !success.Success {
		log.Fatal("Error submitting job")
	}

}
