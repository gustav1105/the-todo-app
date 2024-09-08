package main

import (
    "log"
    "net"

    "the-todo-app/internal/db"
    "the-todo-app/internal/service"
    pb "the-todo-app/proto"
    "google.golang.org/grpc"
)

func main() {
    dbConn, err := db.InitDB()
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }
    defer dbConn.Close()

    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen on port 50051: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterTodoServiceServer(grpcServer, service.NewTodoService(dbConn))

    log.Println("gRPC server running on port 50051")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }
}
