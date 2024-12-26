package service

import (
	"github.com/gdochadipa/oauth2-go-project/internal/repository"
)

type Service interface {
	ItemService
}

type oauthService struct {
	repository repository.Repository
}

func NewService(r repository.Repository) Service {
	return &oauthService{r}
}
