package service

import (
	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
)

type service struct {
	pb.UnimplementedItemServiceServer
	items map[string]*pb.Item
}

// func (s *service) GetItemDetail(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
// 	id := req.GetId()

// 	s.items[id] =
// }
