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
	Role        string             `json:"role,omitempty" bson:"role,omitempty"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
}

type Todolist struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userid,omitempty" json:"userid,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     string             `bson:"duedate,omitempty" json:"duedate,omitempty"`
	Priority    int                `json:"priority,omitempty" bson:"priority,omitempty"`
	Completed   bool               `json:"completed,omitempty" bson:"completed,omitempty"`
}

type Payload struct {
	ID    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}
