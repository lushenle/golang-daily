package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lushenle/golang-daily/remote-server/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewFileWalkClient(conn)

	// Set the directory path here
	path := "/tools/upload"
	response, err := c.GetFiles(context.Background(), &pb.Path{Filename: path})
	if err != nil {
		log.Fatalf("could not get files: %v", err)
	}

	log.Println("Files: ", response.Files)

	fileCheck, err := c.CheckFileMD5(context.Background(), &pb.FileCheckRequest{Filenames: response.Files})
	if err != nil {
		log.Fatalf("could not get files md5sum: %v", err)
	}
	fmt.Println(fileCheck.FileHash)
}
