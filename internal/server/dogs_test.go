package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/evertras/sample-go-app/internal/db"
)

type mockDogGetter struct {
	dogs         []db.Dog
	pendingError error
}

// Generates a dog; we use this so we make sure our tests are using robust data,
// since if the schema changes this signature will change and all calls will
// then fail until we properly fix our tests.  It will get verbose if Dog gets
// to be big, but this is a test so we're ok with that.  Any new fields should
// just be appended to the end of the parameter list.
func genDog(id int64, name string, owner string, location string) db.Dog {
	return db.Dog{
		Id:       id,
		Name:     name,
		Owner:    owner,
		Location: location,
	}
}

func (m *mockDogGetter) GetAllDogs(ctx context.Context) ([]db.Dog, error) {
	return m.dogs, m.pendingError
}

var _ = Describe("fromDbDog", func() {
	It("fills in all expected fields", func() {
		var id int64 = 123
		name := "Rex"
		owner := "Ginkgo"
		location := "New York"

		dbDog := genDog(id, name, owner, location)
		apiDog := fromDbDog(dbDog)

		Expect(apiDog.Id).To(Equal(id))
		Expect(apiDog.Name).To(Equal(name))
		Expect(apiDog.Location).To(Equal(location))
	})
})

var _ = Describe("handlerGetAllDogs", func() {
	When("a request is made that accepts application/json", func() {
		var r *http.Request
		var w *httptest.ResponseRecorder

		BeforeEach(func() {
			r = httptest.NewRequest("GET", "/dogs", nil)
			w = httptest.NewRecorder()

			r.Header.Set("Accept", "application/json")
		})

		When("dogs are returned by the repository", func() {
			var handler http.HandlerFunc
			var getter *mockDogGetter
			var dbDogs []db.Dog

			BeforeEach(func() {
				dbDogs = []db.Dog{
					genDog(
						103,
						"TestDog",
						"Ginkgo",
						"Everywhere",
					),
					genDog(
						33813,
						"Some other dog",
						"Gomega",
						"Paris",
					),
				}

				getter = &mockDogGetter{
					dogs:         dbDogs,
					pendingError: nil,
				}

				handler = handlerGetAllDogs(getter)
			})

			It("returns the dogs as a JSON array", func() {
				handler(w, r)

				Expect(w.Result().StatusCode).To(Equal(200))
				Expect(len(w.Body.Bytes())).NotTo(Equal(0), "Body came back empty")

				var parsedDogs []Dog
				err := json.Unmarshal(w.Body.Bytes(), &parsedDogs)

				Expect(err).NotTo(HaveOccurred(), "Failed to unmarshal response")
				Expect(len(parsedDogs)).To(Equal(len(dbDogs)), "Returned wrong number of dogs")

				for i, dbDog := range dbDogs {
					Expect(parsedDogs[i]).To(Equal(fromDbDog(dbDog)))
				}
			})
		})
	})
})
