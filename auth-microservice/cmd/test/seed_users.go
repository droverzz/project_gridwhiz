package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func generateRandomEmail(index int) string {
	return fmt.Sprintf("user%d@example.com", index)
}

func generateRandomUsername(index int) string {
	return fmt.Sprintf("user%d", index)
}

func generateUsers(count int, startIndex int) []interface{} {
	var users []interface{}
	password := hashPassword("password123")

	for i := 0; i < count; i++ {
		user := bson.M{
			"username":     generateRandomUsername(startIndex + i),
			"email":        generateRandomEmail(startIndex + i),
			"password":     password,
			"is_active":    true,
			"created_date": time.Now(),
		}
		users = append(users, user)
	}
	return users
}

func countUsers(collection *mongo.Collection) int64 {
	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal("CountDocuments error:", err)
	}
	return count
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("gridwhizdb").Collection("users")

	batchSize := 1000
	total := 100000

	fmt.Println("Seeding users...")
	for i := 0; i < total; i += batchSize {
		batch := generateUsers(batchSize, i)
		_, err := collection.InsertMany(context.TODO(), batch)
		if err != nil {
			log.Fatal("InsertMany error:", err)
		}
		fmt.Printf("Inserted %d users\n", i+batchSize)
	}

	fmt.Println("✅ Seeding complete.")

	userCount := countUsers(collection)
	fmt.Printf("Total users in DB: %d\n", userCount)
}
