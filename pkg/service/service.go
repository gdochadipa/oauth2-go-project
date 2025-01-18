package service

import (
	"github.com/gdochadipa/oauth2-go-project/internal/repository"
)

type ItemServiceServer struct {
	repository repository.Repository
}

func NewItemService(r repository.Repository) ItemService {
	return &ItemServiceServer{r}
}

type ServiceInterface interface {
	GrantInterface
}

type ServiceServer struct {
	repository repository.Repository
	jwt        JWTInterface
}

func NewGrantService(r repository.Repository, jwt JWTInterface) ServiceInterface {
	return &ServiceServer{repository: r, jwt: jwt}

}
