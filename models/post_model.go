package models

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
)

type PostgisGeometry struct {
	orb.Geometry
	SRID int
}

func (g *PostgisGeometry) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var data []byte
	var err error

	switch value := value.(type) {
	case []uint8:
		str := string(value)

		data, err = hex.DecodeString(str)
		if err != nil {
			return err
		}

		g.Geometry, g.SRID, err = ewkb.Unmarshal(data)

		return err
	default:
		return fmt.Errorf("expected string but got %T", value)
	}
}

func (g *PostgisGeometry) Value() (driver.Value, error) {
	if g.Geometry == nil {
		return nil, nil
	}

	d := ewkb.MustMarshalToHex(g.Geometry, g.SRID)

	return d, nil
}

type Post struct {
	Base
	Title                   string          `json:"title"`
	Description             string          `json:"description"`
	RequiredVolunteersCount int             `json:"requiredVolunteersCount"`
	RequiredSkills          []string        `json:"requiredSkills"`
	Media                   []*PostMedia    `bun:"rel:has-many,join:id=post_id" json:"media"`
	Latitude                float64         `json:"latitude"`
	Longitude               float64         `json:"longitude"`
	Location                PostgisGeometry `bun:"type:location" json:"-"`
	UserID                  uuid.UUID       `bun:"type:uuid" json:"userID"`
	User                    *User           `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	CategoryID              uuid.UUID       `bun:"type:uuid" json:"categoryID"`
	Category                *Category       `bun:"rel:belongs-to,join:category_id=id" json:"category"`
}
