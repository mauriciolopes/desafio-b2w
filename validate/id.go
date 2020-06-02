package validate

import (
	"desafio-b2w/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ID represents ID validation.
type ID struct {
	hex string
}

// NewID returns new ID validation.
func NewID(hex string) ID {
	return ID{hex}
}

// Validate executes validation.
func (id ID) Validate() error {
	_, err := primitive.ObjectIDFromHex(id.hex)
	if err != nil {
		return errors.InputError("ID inválido")
	}
	return nil
}
