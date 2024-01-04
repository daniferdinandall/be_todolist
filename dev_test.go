package todolist

import (
	"fmt"
	"testing"

	model "github.com/daniferdinandall/be_todolist/model"
	module "github.com/daniferdinandall/be_todolist/module"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db = module.MongoConnect("MONGOSTRING", "db_todolist")

func TestInsertDoc(t *testing.T) {
	var doc model.Todolist
	doc.Title = "Test2"
	doc.Description = "Test"
	doc.DueDate = "Test"
	doc.Priority = 1
	doc.Completed = true

	module.InsertDoc(db, "todolist", doc)
}

func TestSignUp(t *testing.T) {
	var doc model.User
	doc.Name = "dani ferdinan"
	doc.Email = "dani@mail.com"
	doc.Password = "dani1234"
	doc.PhoneNumber = "625156122123"

	err := module.SignUp(db, doc)
	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", doc.Name)
	}
}

func TestSignIn(t *testing.T) {
	var doc model.User
	doc.Email = "dani@mail.com"
	doc.Password = "dani1234"
	user, Status, err := module.SignIn(db, doc)
	fmt.Println("Status :", Status)
	if err != nil {
		t.Errorf("Error getting document: %v", err)
	} else {
		fmt.Println("Welcome bang:", user)
	}
}

func TestCreateTodolist(t *testing.T) {
	var doc model.Todolist
	doc.Title = "Test3"
	doc.Description = "Test"
	doc.DueDate = "Test"
	doc.Priority = 1
	doc.Completed = true

	err := module.CreateTodolist(db, doc)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDoc(t *testing.T) {

	id, err := primitive.ObjectIDFromHex("65964839c2d27eb6a1456f9a")
	if err != nil {
		t.Error(err)
	}

	var doc model.Todolist
	doc.ID = id
	doc.Title = "TestUpdate"
	doc.Description = "Test"
	doc.UserID = "659647e862c130d11ec816c5"

	err = module.UpdateTodolist(db, doc)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(doc)
}

func TestGetAllDocByUserID(t *testing.T) {
	var docs []model.Todolist
	var doc model.Todolist
	doc.UserID = "659647e862c130d11ec816c5"
	docs, err := module.GetAllTodolistByUserID(db, doc)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(docs)
}

func TestGetDocByID(t *testing.T) {
	var doc model.Todolist
	id, err := primitive.ObjectIDFromHex("659647e862c130d11ec816c5")
	if err != nil {
		t.Error(err)
	}
	doc, err = module.GetTodolistByID(db, id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(doc)
}

func TestDeleteDoc(t *testing.T) {
	var doc model.Todolist
	id, err := primitive.ObjectIDFromHex("659647e862c130d11ec816c5")
	if err != nil {
		t.Error(err)
	}
	doc.ID = id
	err = module.DeleteTodolist(db, doc)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(doc)
}
