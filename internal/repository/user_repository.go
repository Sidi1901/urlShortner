package repository

import (
	"fmt"
	"time"
	"errors"
	"context"
	"github.com/jmoiron/sqlx"

)

type UserRepository interface {
	CreateUser()
	GetUser()
	UpdateUser()
	DeleteUser()
}

