package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/evertras/sample-go-app/internal/db"
)

// Dog matches our external API spec, but may not match our internal database
type Dog struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func fromDbDog(dbDog db.Dog) Dog {
	return Dog{
		Id:       dbDog.Id,
		Name:     dbDog.Name,
		Location: dbDog.Location,
	}
}

// DogGetter knows how to get dogs
type DogGetter interface {
	GetAllDogs(ctx context.Context) ([]db.Dog, error)
}

func handlerGetAllDogs(getter DogGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbDogs, err := getter.GetAllDogs(r.Context())

		if err != nil {
			w.WriteHeader(500)
			return
		}

		responseDogs := []Dog{}

		for _, dbDog := range dbDogs {
			responseDogs = append(
				responseDogs,
				fromDbDog(dbDog),
			)
		}

		marshaled, err := json.Marshal(responseDogs)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Write(marshaled)
	}
}
