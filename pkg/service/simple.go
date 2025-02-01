package service

import (
	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
)

type service struct {
	pb.UnimplementedItemServiceServer
	items map[string]*pb.Item
}
