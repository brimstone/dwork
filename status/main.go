package status

import (
	"context"
	"fmt"
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

	jobs, err := c.GetAllJobs(context.Background(), &pb.JobParams{})

	fmt.Println("Job      Shards Found? Location")
	for _, job := range jobs.Statuses {
		fmt.Printf("%8.8s %6.0d %t  \"%0.80s\"\n",
			job.Name,
			job.Shards,
			job.Found,
			job.Location,
		)
	}
}
