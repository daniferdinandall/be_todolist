package todolist

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	model "github.com/daniferdinandall/be_todolist/model"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(MongoString, dbname string) *mongo.Database {
	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://admin:admin@projectexp.pa7k8.gcp.mongodb.net"))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	// Define the regular expression pattern for numeric characters
	numericPattern := `^[0-9]+$`

	// Compile the numeric pattern
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	// Check if the phone number consists only of numeric characters
	if !numericRegexp.MatchString(phoneNumber) {
		return false, nil
	}

	// Define the regular expression pattern for "62" followed by 6 to 12 digits
	pattern := `^62\d{6,13}$`

	// Compile the regular expression
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	// Test if the phone number matches the pattern
	isValid := regexpPattern.MatchString(phoneNumber)

	return isValid, nil
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func GetDoc(db *mongo.Database, col string, doc interface{}, filter bson.M) interface{} {
	collection := db.Collection(col)
	err := collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	return doc
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func InsertDoc(db *mongo.Database, col string, doc interface{}) {
	collection := db.Collection(col)
	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		fmt.Println("Error InsertDoc in colection", col, ":", err)
	}
}

func GetUserFromPhoneNumber(phoneNumber string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"phoneNumber": phoneNumber}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// Main Functions

func SignUp(db *mongo.Database, insertedDoc model.User) error {
	var col = "user"
	objectId := primitive.NewObjectID()

	if insertedDoc.Name == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.PhoneNumber == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	valid, _ := ValidatePhoneNumber(insertedDoc.PhoneNumber)
	if !valid {
		return fmt.Errorf("nomor telepon tidak valid")
	}

	// phoneNumberExists, _ := GetUserFromPhoneNumber(insertedDoc.PhoneNumber, db)
	// if insertedDoc.PhoneNumber == phoneNumberExists.PhoneNumber {
	// 	return fmt.Errorf("nomor telepon sudah terdaftar")
	// }

	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}
	userExists := GetDoc(db, col, model.User{}, bson.M{"email": insertedDoc.Email})

	if insertedDoc.Email == userExists {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}

	hash, _ := HashPassword(insertedDoc.Password)

	var doc model.User
	doc.ID = objectId
	doc.Email = insertedDoc.Email
	doc.Password = hash
	doc.Name = insertedDoc.Name
	doc.PhoneNumber = insertedDoc.PhoneNumber
	doc.Role = "user"

	InsertDoc(db, col, doc)

	return nil
}

func SignIn(db *mongo.Database, col string, insertedDoc model.User) (user model.User, Status bool, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, false, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, false, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	if !CheckPasswordHash(insertedDoc.Password, existsDoc.Password) {
		return user, false, fmt.Errorf("password salah")
	}

	return existsDoc, true, nil
}

func CreateTodolist(db *mongo.Database, col string, insertedDoc model.Todolist) (doc model.Todolist, err error) {
	if insertedDoc.Title == "" || insertedDoc.Description == "" || insertedDoc.DueDate == "" || insertedDoc.Priority == 0 {
		return doc, fmt.Errorf("mohon untuk melengkapi data")
	}
	objectId := primitive.NewObjectID()

	doc.ID = objectId
	doc.UserID = insertedDoc.UserID
	doc.Title = insertedDoc.Title
	doc.Description = insertedDoc.Description
	doc.DueDate = insertedDoc.DueDate
	doc.Priority = insertedDoc.Priority
	doc.Completed = false

	InsertDoc(db, col, doc)

	return doc, nil
}

func GetTodolist(db *mongo.Database, col string, filter bson.M) (docs []model.Todolist, err error) {
	collection := db.Collection(col)
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs, nil
}

func GetTodolistByID(db *mongo.Database, col string, filter bson.M) (doc model.Todolist, err error) {
	collection := db.Collection(col)
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	return doc, nil
}

func UpdateTodolist(db *mongo.Database, col string, filter bson.M, update bson.M) (doc model.Todolist, err error) {
	collection := db.Collection(col)
	err = collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	return doc, nil
}

func DeleteTodolist(db *mongo.Database, col string, filter bson.M) (doc model.Todolist, err error) {
	collection := db.Collection(col)
	err = collection.FindOneAndDelete(context.Background(), filter).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	return doc, nil
}
