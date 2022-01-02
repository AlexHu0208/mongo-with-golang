package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	usersCollection := client.Database("testing").Collection("users")

	// Insert

	// user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	// insertOne(usersCollection, user)

	// users := []interface{}{
	// 	bson.D{{"fullName", "User 2"}, {"age", 25}},
	// 	bson.D{{"fullName", "User 3"}, {"age", 20}},
	// 	bson.D{{"fullName", "User 4"}, {"age", 28}},
	// }
	// insertMany(usersCollection, users)

	// Query
	// find(usersCollection)

	// Update
	// update(usersCollection)
	// updateOneAndMany(usersCollection)
	// replaceOne(usersCollection)

	// Delete
	delete(usersCollection)
}

func insertOne(collection *mongo.Collection, user primitive.D) {
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(nil)
	}
	fmt.Println(result.InsertedID)
}

func insertMany(collection *mongo.Collection, users []interface{}) {
	results, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		panic(err)
	}
	fmt.Println(results.InsertedIDs...)
}

func find(collection *mongo.Collection) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println("displaying all results from the search query")
	for _, result := range results {
		fmt.Println(result)
	}

	var result bson.M
	if err = collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("displaying the first result from the search filter")
	fmt.Println(result)
}

func update(collection *mongo.Collection) {
	user := bson.D{{"fullName", "User 5"}, {"age", 22}}
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}

	update := bson.D{
		{"$set",
			bson.D{
				{"fullName", "User V"},
			},
		},
		{"$inc",
			bson.D{
				{"age", 1},
			},
		},
	}

	result, err := collection.UpdateByID(context.TODO(), insertResult.InsertedID, update)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of documents updated: ", result.ModifiedCount)
}

func updateOneAndMany(collection *mongo.Collection) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}

	update := bson.D{
		{"$set",
			bson.D{
				{"age", 40},
			},
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of documents updated: ", result.ModifiedCount)

	results, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of documents updated: ", results.ModifiedCount)
}

func replaceOne(collection *mongo.Collection) {
	filter := bson.D{{"fullName", "User 1"}}

	replacement := bson.D{
		{"firstName", "Alex"},
		{"lastName", "Hu"},
		{"age", 30},
		{"emailAddress", "alexhu@email.com"},
	}

	result, err := collection.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of documents update: ", result.ModifiedCount)
}

func delete(collection *mongo.Collection) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"age", bson.D{{"$gt", 25}}},
				},
			},
		},
	}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("deleting the first result from the search filter")
	fmt.Println("Number of documents deleted:", result.DeletedCount)

	results, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("deleting every result from the search filter")
	fmt.Println("Number of documents deleted:", results.DeletedCount)
}
