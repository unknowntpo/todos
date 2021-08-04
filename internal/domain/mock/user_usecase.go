package mock

import (
	"github.com/unknowntpo/todos/internal/domain"
)

type MockUserUsecase struct{}

func NewUserUsecase() domain.UserUsecase {
	return &MockTaskUsecase{}
}
