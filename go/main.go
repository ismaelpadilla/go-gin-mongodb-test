package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Stuff struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title,omitempty"`
	Body  string             `bson:"body,omitempty"`

	// json keys returns how the json object is marshalled and unmarshalled
	// Body  string             `bson:"body,omitempty" json:"text"`
}

var client *mongo.Client
var collection *mongo.Collection
var ctx = context.TODO()

func mongodbInit() {
	// mongodb stuff
	var err error

	// default mongo db url
	mongoURL := "mongodb://root:root@localhost:27017/"

	// read mongo db url from environment
	if envURL, isSet := os.LookupEnv("MONGODB_URL"); isSet {
		mongoURL = envURL
		log.Print("mongo db url read from environment")
	} else {
		log.Print("using default mongo db url")
	}

	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("test").Collection("posts")
}

func main() {
	mongodbInit()
	defer client.Disconnect(ctx)

	configureRoutes()
}

func configureRoutes() {
	router := gin.Default()

	router.GET("/test", TestHandler)
	router.POST("/stuff", SaveStuff)
	router.GET("/stuff", GetAllStuff)
	router.GET("/stuff/:id", GetStuffById)
    router.DELETE("/stuff/:id", DeleteStuffById)

	router.Run()
}

// return all objects
func GetAllStuff(c *gin.Context) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Print(err)
	}

	var allStuffDecoded []Stuff
	if err = cursor.All(ctx, &allStuffDecoded); err != nil {
		log.Print(err)
	}

	c.JSON(http.StatusOK, allStuffDecoded)
}

// return object by id
func GetStuffById(c *gin.Context) {
	// get param from url
	id := c.Param("id")

	// build primitive to use in query filter
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		return
	}

	// query db
	result := collection.FindOne(ctx, bson.D{{"_id", idPrimitive}})

	if result.Err() == mongo.ErrNoDocuments {
		c.Status(http.StatusNotFound)
		return
	}

	// decode into object of desired type
	var decoded Stuff
	err = result.Decode(&decoded)

	if err != nil {
		log.Print(err)
	}

	// return object
	c.JSON(http.StatusOK, decoded)
}

func SaveStuff(c *gin.Context) {
	// decode body into object
	var decoded Stuff
	c.BindJSON(&decoded)

	// set object id
	decoded.ID = primitive.NewObjectID()

	// insert object
	result, err := collection.InsertOne(ctx, decoded)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// get inserted object id
	id := result.InsertedID.(primitive.ObjectID)
	idBytes, err := id.MarshalJSON()
	if err != nil {
		log.Print(err)
		return
	}

	idString := strings.Trim(string(idBytes), "\"")

	// primitive.NewObjectID().String()
	// set location header
	location := c.Request.Host + "/stuff/" + idString
    c.Header("location", location)

	// w.WriteHeader(http.StatusCreated)

	// return object
    c.JSON(http.StatusCreated, decoded)
}

// delete object by id
func DeleteStuffById(c *gin.Context) {
	// get param from url
	id := c.Param("id")

	// build primitive to use in query filter
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// query db
	result, err := collection.DeleteOne(ctx, bson.D{{"_id", idPrimitive}})
	if err != nil {
		log.Print(err)
		return
	}

	// set header to not found if nothing was deleted
	if result.DeletedCount == 0 {
		c.Status(http.StatusNotFound)
	} else {
		c.Status(http.StatusOK)
	}
}

// test function
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
