package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DBName is name of the database.
	DBName = "glottery"
	// Name of the collection.
	notesCollection = "notes"
	//URI = "mongodb://<user>:<password>@<host>/<name>"
	URI = "mongodb://127.0.0.1:27017"
)

// Note is structure of the record.
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

func main() {
	// Base context.
	ctx := context.Background()

	// Options to the database.
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := client.Database(DBName)
	fmt.Println(db.Name()) // output: glottery

	coll := db.Collection(notesCollection)
	fmt.Println(coll.Name()) // output: notes

	// Insert One Document.
	note := Note{}

	// An ID for MongoDB.
	note.ID = primitive.NewObjectID()
	note.Title = "First note"
	note.Body = "Some spam text"
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	result, err := coll.InsertOne(ctx, note)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ID of the inserted document.
	objectID := result.InsertedID.(primitive.ObjectID)
	fmt.Println(objectID)

	// Insert Many Documents.
	notes := []interface{}{}

	note2 := Note{}
	note2.ID = primitive.NewObjectID()
	note2.Title = "Second note"
	note2.Body = "Some spam text"
	note2.CreatedAt = time.Now()
	note2.UpdatedAt = time.Now()

	note3 := Note{
		ID:        primitive.NewObjectID(),
		Title:     "Third note",
		Body:      "Some spam text",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	notes = append(notes, note2, note3)

	results, err := coll.InsertMany(ctx, notes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(results.InsertedIDs)

	// Update a document.

	// Parsing a string ID to ObjectID from MongoDB.
	objID, err := primitive.ObjectIDFromHex("5fcb95a4a7b657b6d0579569")
	if err != nil {
		fmt.Println(err)
		return
	}

	resultUpdate, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"body":       "Some updated text",
				"updated_at": time.Now(),
			},
		},
	)

	fmt.Println(resultUpdate.ModifiedCount) // output: 1

	// Delete one document.

	objID, err = primitive.ObjectIDFromHex("5fcb95a4a7b657b6d057956b")
	if err != nil {
		fmt.Println(err)
		return
	}

	resultDelete, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resultDelete.DeletedCount) // output: 1

	// Find documents.

	// Find one document.
	objID, err = primitive.ObjectIDFromHex("5fcb95695e67cf9bab5c1519")
	if err != nil {
		fmt.Println(err)
		return
	}

	findResult := coll.FindOne(ctx, bson.M{"_id": objID})
	if err := findResult.Err(); err != nil {
		fmt.Println(err)
		return
	}

	n := Note{}
	err = findResult.Decode(&n)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(n.Body)

	//Find many documents
	notesResult := []Note{}
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through the returned cursor.
	for cursor.Next(ctx) {
		cursor.Decode(&n)
		notesResult = append(notesResult, n)
	}

	for _, el := range notesResult {
		fmt.Println(el.Title)
	}

}
