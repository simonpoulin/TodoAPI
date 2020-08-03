package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Status  bool               `json:"isComplete" bson:"isComplete"`
	Content string             `json:"content" bson:"content" validate:"min=1"`
}
