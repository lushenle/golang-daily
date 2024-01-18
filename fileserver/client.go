package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/lushenle/golang-daily/fileserver/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Set up the credentials for the connection.
	creds, err := credentials.NewClientTLSFromFile("Server.crt", "")
	if err != nil {
		log.Fatalf("fail to create credentials %v", err)
	}

	// Initialize a gRPC client connection
	conn, err := grpc.Dial("localhost:4040", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServerClient(conn)

	// Open file
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer file.Close()

	// Create data stream
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("Could not upload file: %v", err)
	}

	// Send the filename
	err = stream.Send(&pb.FileData{Filename: file.Name()})
	if err != nil {
		log.Fatal(err)
	}

	// Define the data block size and send the file data block by block
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Could not read chunk to buffer: %v", err)
		}

		err = stream.Send(&pb.FileData{
			Data: buffer[:n],
		})
		if err != nil {
			log.Fatalf("Could not send chunk to Server: %v", err)
		}
	}

	// Close the stream and get the response
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Could not receive response: %v", err)
	}

	log.Printf("Status: %v, Message: %v\n", res.GetCode(), res.GetMessage())

	// Set the directory path here
	path := "/tools/upload"
	response, err := client.GetFiles(context.Background(), &pb.Path{Filename: path})
	if err != nil {
		log.Fatalf("could not get files: %v", err)
	}

	log.Println("Files:", response.Files)

	fileCheck, err := client.CheckFileMD5(context.Background(), &pb.FileCheckRequest{Filenames: response.Files})
	if err != nil {
		log.Fatalf("could not get files md5sum: %v", err)
	}
	log.Println("FileHash:", fileCheck.FileHash)

	// Download file
	downloadFile, err := client.DownloadFile(context.Background(), &pb.DownloadRequest{Filename: filepath.Join(path, "test.txt")})
	if err != nil {
		log.Fatalf("could not download file: %v", err)
	}

	f, err := os.Create("myDownload.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		chunk, err := downloadFile.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		_, err = f.Write(chunk.Chunk)
		if err != nil {
			panic(err)
		}
	}

	log.Println("File downloaded successfully")
}
