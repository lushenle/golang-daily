package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/lushenle/golang-daily/remote-server/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileWalkServer
}

func (s *server) CheckFileMD5(ctx context.Context, in *pb.FileCheckRequest) (*pb.FileCheckResponse, error) {
	fileHash, err := checkSum(in.Filenames)
	if err != nil {
		return nil, err
	}

	return &pb.FileCheckResponse{
		FileHash: fileHash,
	}, nil
}

func checkSum(files []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, f); err != nil {
			return nil, err
		}

		m[file] = hex.EncodeToString(hash.Sum(nil))
	}

	return m, nil
}

func (s *server) GetFiles(ctx context.Context, path *pb.Path) (*pb.FileList, error) {
	files, err := fileWalk(path.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.FileList{Files: files}, nil
}

func fileWalk(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func main() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileWalkServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
