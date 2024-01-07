package todolist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	Base64Url   string             `bson:"base64url,omitempty" json:"base64url,omitempty"`
}

type Todolist struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      string             `bson:"userid,omitempty" json:"userid,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     int64              `bson:"duedate,omitempty" json:"duedate,omitempty"`
	Priority    int                `json:"priority,omitempty" bson:"priority,omitempty"`
	Completed   bool               `json:"completed,omitempty" bson:"completed,omitempty"`
	CreatedAt   int64              `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

type Payload struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type TodolistResponse struct {
	Status  bool       `json:"status" bson:"status"`
	Message string     `json:"message,omitempty" bson:"message,omitempty"`
	Data    []Todolist `json:"data" bson:"data"`
}

type ProfileResponse struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Data    User   `json:"data" bson:"data"`
}
