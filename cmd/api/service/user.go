package service

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/go-project-layout/cmd/api/model"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type UserService struct {
	logger log.Logger
}

func NewUserService(logger log.Logger) *UserService {
	return &UserService{
		logger: logger,
	}
}

func (s *UserService) Handle(_ context.Context, users []model.User) error {
	fmt.Println(users)
	return nil
}
