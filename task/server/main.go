// Original code from https://www.youtube.com/watch?v=6Ol6zeocR28&t=100s

package main

import (
	"context"
	"fmt"
	pb2 "go-study/task/pb/proto"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type taskServer struct {
	pb2.UnimplementedTaskManagerServer
	mu      sync.Mutex
	tasks   []*pb2.Task
	counter int
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb2.CreateTaskRequest) (*pb2.CreateTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("task-%d", s.counter)

	task := &pb2.Task{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}
	s.tasks = append(s.tasks, task)

	log.Printf("created task: %s", id)

	return &pb2.CreateTaskResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func (s *taskServer) ListTasks(req *pb2.ListTasksRequest, stream pb2.TaskManager_ListTasksServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb2.RegisterTaskManagerServer(grpcServer, &taskServer{})

	log.Println("grpc server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
