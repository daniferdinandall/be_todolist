package todolist

import (
	"fmt"
	"testing"
	"time"

	model "github.com/daniferdinandall/be_todolist/model"
	module "github.com/daniferdinandall/be_todolist/module"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/whatsauth/watoken"
)

var db = module.MongoConnect("MONGOSTRING", "db_todolist")

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)
}

func TestEncodeToken(t *testing.T) {
	privateKey := "ae647abe102886b0cca6eac1a9ab78174d209466eaddffe977040fb0badcae87e00b393053b1e23efa50af4d919a1cfb853fcc4cb0cb5b5cce6fc73088fc1722"
	userid := "659654c7931d81bd72a9c4c6"
	tokenstring, err := watoken.Encode(userid, privateKey)

	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", tokenstring)
	}
}

func TestDecodeToken(t *testing.T) {
	publicKey := "e00b393053b1e23efa50af4d919a1cfb853fcc4cb0cb5b5cce6fc73088fc1722"
	tokenstring := "v4.public.eyJleHAiOiIyMDI0LTAxLTA0VDE2OjIzOjEyKzA3OjAwIiwiaWF0IjoiMjAyNC0wMS0wNFQxNDoyMzoxMiswNzowMCIsImlkIjoiNjU5NjU0Yzc5MzFkODFiZDcyYTljNGM2IiwibmJmIjoiMjAyNC0wMS0wNFQxNDoyMzoxMiswNzowMCJ9lr0vK81FvhZmHt7BDcbTbR-ylZFEEs80CE99NhqGr_JLeOtc5_0La4glOt2JfqdcCftQdwE9hntMDOS5R9eYDg"
	useridstring, err := watoken.Decode(publicKey, tokenstring)

	if err != nil {
		t.Errorf("Error inserting document: %v", err)
	} else {
		fmt.Println("Data berhasil disimpan dengan nama :", useridstring.Id)
	}
}

func TestInsertDoc(t *testing.T) {
	var doc model.Todolist
	doc.Title = "Test2"
	doc.Description = "Test"
	doc.DueDate = time.Now().Unix()
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
		fmt.Println("Welcome bang:", user.ID.Hex())
	}
}

func TestCreateTodolist(t *testing.T) {
	var doc model.Todolist
	doc.Title = "test time2"
	doc.Description = "Test"
	doc.DueDate = time.Now().Unix()
	doc.Priority = 1
	doc.Completed = true

	err := module.CreateTodolist(db, doc)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateDoc(t *testing.T) {

	id, err := primitive.ObjectIDFromHex("659426d7294a5abb63c66959")
	if err != nil {
		t.Error(err)
	}

	var doc model.Todolist
	doc.ID = id

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
	err = module.DeleteTodolist(db, doc.ID)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(doc)
}

func TestGetProfile(t *testing.T) {
	var doc model.User
	id, err := primitive.ObjectIDFromHex("659654c7931d81bd72a9c4c6")
	if err != nil {
		t.Error(err)
	}
	doc.ID = id
	doc, img, err := module.GetProfile(db, doc.ID)
	if err != nil {
		t.Error(err)
	}
	if img != "" {
		fmt.Println(img)
	}
	fmt.Println(doc)
}

func TestUpdateProfile(t *testing.T) {
	var doc model.User
	id, err := primitive.ObjectIDFromHex("659654c7931d81bd72a9c4c6")
	if err != nil {
		t.Error(err)
	}
	doc.ID = id
	doc.Name = "dani ferdinan"
	doc.PhoneNumber = "625156122123"

	var img model.Image
	img.Base64Url = "sasa"

	err = module.UpdateProfile(db, doc, img.Base64Url)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(doc)
}
