package service

import (
	"context"
)

func (s *Service) CreateUser(ctx context.Context, email string, name string, password string, usertype string) error {

	// userModel := &model.User{
	// 	Email:    &email,
	// 	Name:     &name,
	// 	Password: &password,
	// 	UserType: usertype,
	// }
}
