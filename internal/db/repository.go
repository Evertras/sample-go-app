package db

import (
	"context"
	"errors"
)

type Repository struct {
}

func (r *Repository) GetAllDogs(ctx context.Context) ([]Dog, error) {
	return nil, errors.New("not implemented")
}

func (r *Repository) GetDog(ctx context.Context) (Dog, error) {
	return Dog{}, errors.New("not implemented")
}
