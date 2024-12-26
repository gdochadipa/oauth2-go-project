package server

import (
	"context"
	"fmt"
	"net"

	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
	"github.com/gdochadipa/oauth2-go-project/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service service.Service
}

// DeleteItem implements pb.ItemServiceServer.
func (g *grpcServer) DeleteItem(context.Context, *pb.DeleteItemRequest) (*pb.DeletetemResponse, error) {
	panic("unimplemented")
}

// GetItem implements pb.ItemServiceServer.
func (g *grpcServer) GetItem(context.Context, *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	panic("unimplemented")
}

// GetItems implements pb.ItemServiceServer.
func (g *grpcServer) GetItems(context.Context, *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	panic("unimplemented")
}

// StoreItem implements pb.ItemServiceServer.
func (g *grpcServer) StoreItem(context.Context, *pb.StoreItemRequest) (*pb.StoreItemResponse, error) {
	panic("unimplemented")
}

// UpdateItem implements pb.ItemServiceServer.
func (g *grpcServer) UpdateItem(context.Context, *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedItemServiceServer implements pb.ItemServiceServer.
func (g *grpcServer) mustEmbedUnimplementedItemServiceServer() {
	panic("unimplemented")
}

func ListenGRPC(s service.Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterItemServiceServer(server, &grpcServer{s})
	reflection.Register(server)

	return server.Serve(lis)

}
