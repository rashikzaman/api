package application

import (
	"github.com/uptrace/bun"
)

type Application struct {
	DB *bun.DB
}
