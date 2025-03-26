package application

import (
	"rashikzaman/api/config"

	"github.com/uptrace/bun"
)

type Application struct {
	DB     *bun.DB
	Config config.Config
}
