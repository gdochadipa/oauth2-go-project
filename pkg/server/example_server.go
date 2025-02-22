package server

import (
	"context"
	"fmt"
	"net"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
	"github.com/gdochadipa/oauth2-go-project/pkg/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcExampleServer struct {
	service service.ItemService
	pb.UnimplementedItemServiceServer
}

// DeleteItem implements pb.ItemServiceServer.
func (g *grpcExampleServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeletetemResponse, error) {
	err := g.service.DropItem(ctx, uuid.MustParse(req.Id))

	if err != nil {
		return nil, err
	}

	return &pb.DeletetemResponse{
		Id: req.Id,
	}, nil
}

// GetItem implements pb.ItemServiceServer.
func (g *grpcExampleServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := g.service.GetItem(ctx, uuid.MustParse(req.Id))

	if err != nil {
		return nil, err
	}

	return &pb.GetItemResponse{
		Data: &pb.Item{
			Id:          item.Id.String(),
			Name:        item.Name,
			Description: &item.Description,
		},
	}, err
}

// GetItems implements pb.ItemServiceServer.
func (g *grpcExampleServer) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	items, err := g.service.GetAllItems(ctx, req.Skip, req.LimitPage)

	if err != nil {
		return nil, err
	}

	newItems := []*pb.Item{}

	for _, p := range items {
		newItems = append(
			newItems,
			&pb.Item{
				Id:          p.Id.String(),
				Name:        p.Name,
				Description: &p.Description,
			},
		)
	}

	return &pb.GetItemsResponse{Data: newItems, Meta: &pb.PaginationMeta{Skip: 0, LimitPage: 0}}, nil
}

// StoreItem implements pb.ItemServiceServer.
func (g *grpcExampleServer) StoreItem(ctx context.Context, req *pb.StoreItemRequest) (*pb.StoreItemResponse, error) {
	item, err := g.service.PostItem(ctx, req.Name, *req.Description)
	if err != nil {
		return nil, err
	}

	return &pb.StoreItemResponse{Item: &pb.Item{
		Id:          item.Id.String(),
		Name:        item.Name,
		Description: &item.Description,
	}}, err
}

// UpdateItem implements pb.ItemServiceServer.
func (g *grpcExampleServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	item := &entity.Item{
		Id:          uuid.MustParse(req.Id),
		Name:        *req.Name,
		Description: *req.Description,
	}

	result, err := g.service.PutItem(ctx, *item)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateItemResponse{Item: &pb.Item{
		Id:          result.Id.String(),
		Name:        result.Name,
		Description: &result.Description,
	}}, err
}

func ListenExampleGRPC(s service.ItemService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterItemServiceServer(server, &grpcExampleServer{service: s})
	reflection.Register(server)

	return server.Serve(lis)

}
