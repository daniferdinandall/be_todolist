package todolist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

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

func GetID(r *http.Request) string {
	return r.URL.Query().Get("id")
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

func GetDoc(db *mongo.Database, doc interface{}, filter bson.M) interface{} {
	col := "todolist"
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

	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}

	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == existsDoc.Email {
		return fmt.Errorf("email sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return fmt.Errorf("password terlalu pendek")
	}

	hash, _ := HashPassword(insertedDoc.Password)

	insertedDoc.ID = objectId
	insertedDoc.Password = hash

	collection := db.Collection(col)
	_, err := collection.InsertOne(context.Background(), insertedDoc)
	if err != nil {
		fmt.Println("Error InsertDoc in colection", col, ":", err)
	}
	return nil
}

func SignIn(db *mongo.Database, insertedDoc model.User) (user model.User, Status bool, err error) {

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

func CreateTodolist(db *mongo.Database, doc model.Todolist) (err error) {
	if doc.Title == "" || doc.Description == "" || doc.DueDate == 0 || doc.Priority == 0 || doc.UserID == "" {
		return fmt.Errorf("mohon untuk melengkapi data")
	}
	col := "todolist"
	collection := db.Collection(col)
	// Get the current time
	currentTime := time.Now()

	// Convert time to Unix timestamp
	unixTimestamp := currentTime.Unix()

	doc.CreatedAt = unixTimestamp
	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		fmt.Println("Error InsertDoc in colection", col, ":", err)
	}
	return nil
}

func GetTodolistByID(db *mongo.Database, _id primitive.ObjectID) (doc model.Todolist, err error) {
	col := "todolist"
	filter := bson.M{"_id": _id}
	collection := db.Collection(col)
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	return doc, nil
}

func GetAllTodolistByUserID(db *mongo.Database, doc model.Todolist) (docs []model.Todolist, err error) {
	filter := bson.M{"userid": doc.UserID}
	col := "todolist"
	collection := db.Collection(col)
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	return docs, nil
}

func UpdateTodolist(db *mongo.Database, doc model.Todolist) (err error) {
	filter := bson.M{"_id": doc.ID}
	col := "todolist"
	collection := db.Collection(col)
	result, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}
	return nil
}

func DeleteTodolist(db *mongo.Database, _id primitive.ObjectID) error {
	collection := db.Collection("todolist")
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}
	return nil
}

func UpdateProfile(db *mongo.Database, doc model.User) (err error) {
	col := "user"
	filter := bson.M{"_id": doc.ID}
	collection := db.Collection(col)
	result, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified id")
		return
	}

	return nil
}

func GetProfile(db *mongo.Database, _id primitive.ObjectID) (doc model.User, err error) {
	col := "user"
	filter := bson.M{"_id": _id}
	collection := db.Collection(col)
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		fmt.Println("Error GetDoc in colection", col, ":", err)
	}
	doc = model.User{
		ID:          doc.ID,
		Name:        doc.Name,
		Email:       doc.Email,
		PhoneNumber: doc.PhoneNumber,
		Base64Url:   doc.Base64Url,
	}

	return doc, nil
}
