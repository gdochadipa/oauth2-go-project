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
	service service.ServiceInterface
	pb.UnimplementedOAuthServiceServer
}

func (g *grpcServer) AuthorizeToken(ctx context.Context, req *pb.AuthCodeGrantRequest) (*pb.AuthCodeGrantResponse, error) {

	return nil, nil
}
func (g *grpcServer) ClientCredentGrant(ctx context.Context, req *pb.ClientCredentGrantRequest) (*pb.ClientCredentGrantResponse, error) {
	return nil, nil
}

// this function will trigger for getting auth code, before verified from user
func (g *grpcServer) GenerateAuthCode(ctx context.Context, req *pb.GenerateCodeRequest) (*pb.GenerateCodeResponse, error) {

	authRequest, err := g.service.ValidateAuthorizationRequest(ctx, req)

	if err != nil {
		return nil, err
	}

	// validation user
	// if user not login, will return error
	// i think will send credentials login from Password Token Grant

	authCodeToken, uri, errorAuth := g.service.CompleteAuthorizationRequest(ctx, authRequest)

	if errorAuth != nil {
		return nil, errorAuth
	}

	return &pb.GenerateCodeResponse{
		StateCode: req.State,
		Code:      *authCodeToken,
		Uri:       *uri,
	}, nil
}
func (g *grpcServer) PasswordTokenGrant(ctx context.Context, req *pb.CredentialsGrantRequest) (*pb.CredentialsGrantResponse, error) {
	return nil, nil
}
func (g *grpcServer) RefreshTokenGrant(ctx context.Context, req *pb.RefreshTokenGrantRequest) (*pb.RefreshTokenGrantResponse, error) {
	return nil, nil
}

func ListenGRPC(s service.ServiceInterface, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterOAuthServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)

	return server.Serve(lis)

}
