package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// type Comment struct {
// 	guideid       string `json:"guideid" bson:"guideid"`
// 	useridcomment string `json:"guideid" bson:"guideid"`
// 	useremail     string `json:"guideid" bson:"guideid"`
// 	avatar        string `json:"guideid" bson:"guideid"`
// 	author        string `json:"guideid" bson:"guideid"`
// 	datecomment   string `json:"guideid" bson:"guideid"`
// 	message       string `json:"guideid" bson:"guideid"`
// }

func main() {
	fmt.Println("Starting the application...")

	router := mux.NewRouter()

	router.HandleFunc("/user", GetPeopleEndpoint).Methods("GET")

	http.ListenAndServe(":12345", router)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://baac.topwork.asia:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	localguideCollection := client.Database("fighto")
	commentCollection := localguideCollection.Collection("localguide_users")

	cursor, err := commentCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	// var comments []bson.M
	// if err = cursor.All(ctx, &comments); err != nil {
	// 	log.Fatal(err)
	// }
	// for com := range comments {
	// 	fmt.Println(com["guideid"])
	// }
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var com bson.M
		if err = cursor.Decode(&com); err != nil {
			log.Fatal(err)
		}
		// json.NewEncoder(response).Encode(com)
	}

	filterCursor, err := commentCollection.Find(ctx, bson.M{"usertype": "1"})
	if err != nil {
		log.Fatal(err)
	}
	var episodesFiltered []bson.M
	if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodesFiltered)
	json.NewEncoder(response).Encode(episodesFiltered)
}
