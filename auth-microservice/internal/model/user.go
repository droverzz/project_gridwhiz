package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Role      string             `bson:"role"`
	Password  string             `bson:"password"`
	Deleted   bool               `bson:"deleted"`
	CreatedAt time.Time          `bson:"created_at"`
}
