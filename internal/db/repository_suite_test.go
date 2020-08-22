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
	var ctxTest context.Context
	var mongoContainer testcontainers.Container
	var mongoClient *mongo.Client
	var repository *Repository

	BeforeSuite(func() {
		ctxContainer := context.Background()

		req := testcontainers.ContainerRequest{
			Image: "mongo:4.4",
			ExposedPorts: []string{
				"27017/tcp",
			},
			WaitingFor: wait.ForLog("Waiting for connections"),
		}

		var err error
		mongoContainer, err = testcontainers.GenericContainer(ctxContainer, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		ctxTest = context.Background()

		mongoHost, err := mongoContainer.Host(ctxTest)
		Expect(err).ToNot(HaveOccurred())
		mongoPort, err := mongoContainer.MappedPort(ctxTest, "27017")
		Expect(err).ToNot(HaveOccurred())

		mongoClient, err = mongo.Connect(
			ctxTest,
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort),
			),
		)
		Expect(err).ToNot(HaveOccurred())

		err = mongoClient.Ping(ctxTest, nil)
		Expect(err).ToNot(HaveOccurred())

		repository = New(mongoClient)
	})

	When("the dog collection is empty", func() {
		BeforeEach(func() {
			dogCollection := mongoClient.Database(mongoDbName).Collection(mongoDogCollection)

			Expect(dogCollection).ToNot(BeNil(), "Dog collection not found")
			_, err := dogCollection.DeleteMany(ctxTest, bson.D{})
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GetAllDogs", func() {
			It("returns an empty list of dogs", func() {
				dogs, err := repository.GetAllDogs(ctxTest)

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

			Expect(existingDogs).ToNot(HaveLen(0))
			Expect(dogCollection).ToNot(BeNil(), "Dog collection not found")

			_, err := dogCollection.DeleteMany(ctxTest, bson.D{})
			Expect(err).ToNot(HaveOccurred())

			result, err := dogCollection.InsertMany(ctxTest, insertDogs)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.InsertedIDs).To(HaveLen(len(existingDogs)))
		})

		Describe("GetAllDogs", func() {
			It("returns all dogs", func() {
				dogs, err := repository.GetAllDogs(ctxTest)

				Expect(err).ToNot(HaveOccurred())
				Expect(dogs).To(ConsistOf(existingDogs))
			})
		})
	})

	AfterSuite(func() {
		mongoContainer.Terminate(ctxTest)
	})
})
