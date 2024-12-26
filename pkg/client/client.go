package client

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ItemServiceClient
}

/*
*

ini hubungan antara gRPC dengan server, yaa sebenarnya ga penting
*/
func NewClient(url string) (*Client, error) {
	// conn, err := grpc.Dial(url, grpc.WithInsecure())
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewItemServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) AddItem(ctx context.Context, name string, desc string) (*entity.Item, error) {
	r, err := c.service.StoreItem(ctx, &pb.StoreItemRequest{Name: name, Description: &desc})

	if err != nil {
		return nil, err
	}

	return &entity.Item{
		Id:          uuid.MustParse(r.Item.Id),
		Name:        r.Item.Name,
		Description: *r.Item.Description,
	}, nil
}

func (c *Client) GetAllItems(ctx context.Context, skip int32, limitPage int32) ([]entity.Item, error) {
	r, err := c.service.GetItems(ctx, &pb.GetItemsRequest{Skip: uint32(skip), LimitPage: uint32(limitPage)})

	if err != nil {
		return nil, err
	}

	items := []entity.Item{}

	for _, item := range r.Data {
		items = append(items, entity.Item{
			Id:          uuid.MustParse(item.Id),
			Name:        item.Name,
			Description: *item.Description,
		})
	}
	return items, nil

}
