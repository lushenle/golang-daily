package main

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/lushenle/golang-daily/fileserver/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.UnimplementedFileServerServer
}

func (s *Server) DownloadFile(req *pb.DownloadRequest, stream pb.FileServer_DownloadFileServer) error {
	f, err := os.Open(req.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err = stream.Send(&pb.FileChunk{Chunk: buf[:n]}); err != nil {
			return err
		}
	}
}

func (s *Server) UploadFile(stream pb.FileServer_UploadFileServer) error {
	fd, err := stream.Recv()
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join("/tools/upload", fd.Filename))
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		fd, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Status{
				Code:    201,
				Message: "File uploaded successfully",
			})
		}
		if err != nil {
			return err
		}
		file.Write(fd.Data)
	}
}

func (s *Server) CheckFileMD5(ctx context.Context, in *pb.FileCheckRequest) (*pb.FileCheckResponse, error) {
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

func (s *Server) GetFiles(ctx context.Context, path *pb.Path) (*pb.FileList, error) {
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
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load Server certificates: %v", err)
	}
	// TLS configuration
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	creds := credentials.NewTLS(config)

	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterFileServerServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
