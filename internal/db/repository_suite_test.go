package db

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "db")
}

func genDog(id int64, name string, owner string, location string) Dog {
	return Dog{
		Id:       id,
		Name:     name,
		Owner:    owner,
		Location: location,
	}
}

var _ = Describe("db", func() {
	var ctx context.Context
	var mongoContainer testcontainers.Container
	var mongoClient *mongo.Client
	var repository *Repository

	// Overkill but fun for now
	BeforeEach(func() {
		ctx = context.Background()
		req := testcontainers.ContainerRequest{
			Image: "mongo:4.4",
			ExposedPorts: []string{
				"27017/tcp",
			},
			WaitingFor: wait.ForLog("aiting for connections"),
		}

		var err error
		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		Expect(err).ToNot(HaveOccurred())
		mongoHost, err := mongoContainer.Host(ctx)
		Expect(err).ToNot(HaveOccurred())
		mongoPort, err := mongoContainer.MappedPort(ctx, "27017")
		mongoClient, err = mongo.Connect(
			ctx,
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort),
			),
		)

		Expect(err).ToNot(HaveOccurred())

		err = mongoClient.Ping(ctx, nil)

		Expect(err).ToNot(HaveOccurred())

		repository = New(mongoClient)
	})

	When("the dog collection is empty", func() {
		BeforeEach(func() {
			dogCollection := mongoClient.Database(mongoDbName).Collection(mongoDogCollection)

			Expect(dogCollection).ToNot(BeNil(), "Dog collection not found")
			_, err := dogCollection.DeleteMany(ctx, bson.D{})
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GetAllDogs", func() {
			It("returns an empty list of dogs", func() {
				dogs, err := repository.GetAllDogs(ctx)

				Expect(err).ToNot(HaveOccurred())

				Expect(dogs).To(HaveLen(0))
			})
		})
	})

	When("the dog collection has some dogs", func() {
		var existingDogs []Dog

		BeforeEach(func() {
			existingDogs = []Dog{
				genDog(1234, "Genji", "Ginkgo", "Tokyo"),
				genDog(44444, "Rex", "Ginkgo", "New York"),
			}

			insertDogs := []interface{}{}

			for _, dog := range existingDogs {
				insertDogs = append(insertDogs, dog)
			}

			dogCollection := mongoClient.Database(mongoDbName).Collection(mongoDogCollection)

			Expect(dogCollection).ToNot(BeNil(), "Dog collection not found")
			_, err := dogCollection.DeleteMany(ctx, bson.D{})
			Expect(err).ToNot(HaveOccurred())
			result, err := dogCollection.InsertMany(ctx, insertDogs)
			Expect(err).ToNot(HaveOccurred())

			Expect(result.InsertedIDs).To(HaveLen(len(existingDogs)))
		})

		Describe("GetAllDogs", func() {
			It("returns all dogs", func() {
				dogs, err := repository.GetAllDogs(ctx)

				Expect(err).ToNot(HaveOccurred())
				Expect(dogs).To(HaveLen(len(existingDogs)))
				Expect(dogs).To(ConsistOf(existingDogs))
			})
		})
	})

	AfterEach(func() {
		mongoContainer.Terminate(ctx)
	})
})
