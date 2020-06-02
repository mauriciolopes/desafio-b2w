package models

import (
	"desafio-b2w/validate"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Planet represents a planet.
type Planet struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name"`
	Climate         string             `json:"climate"`
	Terrain         string             `json:"terrain"`
	FilmAppearances uint16             `json:"filmAppearances"`
}

// Validate executes validation.
func (p Planet) Validate() error {
	return validate.Exec(
		validate.NewRequiredString(p.Name, "Nome não informado"),
		validate.NewRequiredString(p.Climate, "Clima não informado"),
		validate.NewRequiredString(p.Terrain, "Terreno não informado"),
	)
}
