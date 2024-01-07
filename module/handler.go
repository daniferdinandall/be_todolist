package todolist

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	model "github.com/daniferdinandall/be_todolist/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func GCFHandlerSignup(MONGOCONNSTRINGENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	err = SignUp(conn, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "Halo " + dataUser.Name
	return GCFReturnStruct(Response)
}

func GCFHandlerSignin(MONGOCONNSTRINGENV, PASETOPRIVATEKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.Credential
	Response.Status = false
	var dataUser model.User
	err := json.NewDecoder(r.Body).Decode(&dataUser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
		return GCFReturnStruct(Response)
	}
	user, status, err := SignIn(conn, dataUser)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	// Response.Message = "Halo " + dataUser.Name
	tokenstring, err := watoken.Encode(dataUser.ID.Hex(), os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Response.Message = "Gagal Encode Token : " + err.Error()
	} else {
		Response.Message = "Selamat Datang " + user.Email + " di Todolist" + strconv.FormatBool(status)
		Response.Token = tokenstring
	}
	return GCFReturnStruct(Response)
}

func GCFCreateTodolist(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.TodolistResponse
	Response.Status = false
	var dataTodolist model.Todolist

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	err = CreateTodolist(conn, dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "success"
	return GCFReturnStruct(Response)
}

func GCFUpdateTodolist(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.TodolistResponse
	Response.Status = false
	var dataTodolist model.Todolist

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	err = UpdateTodolist(conn, dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "success"
	return GCFReturnStruct(Response)
}

func GCFGetAllTodolist(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.TodolistResponse
	Response.Status = false
	var dataTodolist model.Todolist

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err := json.NewDecoder(r.Body).Decode(&dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	docs, err := GetAllTodolistByUserID(conn, dataTodolist)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "success"
	Response.Data = docs
	return GCFReturnStruct(Response)
}

func GCFGetTodolistByID(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.TodolistResponse
	Response.Status = false

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// Get Id
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}

	doc, err := GetTodolistByID(conn, idparam)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "success"
	Response.Data = append(Response.Data, doc)
	return GCFReturnStruct(Response)
}

func GCFDeleteTodolistByID(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.TodolistResponse
	Response.Status = false

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// Get Id
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}

	err = DeleteTodolist(conn, idparam)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}
	Response.Status = true
	Response.Message = "success"
	return GCFReturnStruct(Response)
}

func GCFGetProfile(MONGOCONNSTRINGENV, PASETOPUBLICKEYENV string, r *http.Request) string {
	conn := MongoConnect(MONGOCONNSTRINGENV, "db_todolist")
	var Response model.ProfileResponse
	Response.Status = false

	// get token from header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		Response.Message = "error parsing application/json1:" + token
		return GCFReturnStruct(Response)
	}

	// decode token
	_, err1 := watoken.Decode(os.Getenv(PASETOPUBLICKEYENV), token)

	if err1 != nil {
		Response.Message = "error parsing application/json2: " + err1.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	// Get Id
	id := GetID(r)
	if id == "" {
		Response.Message = "Wrong parameter"
		return GCFReturnStruct(Response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid id parameter"
		return GCFReturnStruct(Response)
	}

	doc, img, err := GetProfile(conn, idparam)
	if err != nil {
		Response.Message = "error parsing application/json4: " + err.Error()
		return GCFReturnStruct(Response)
	}

	if img != "" {
		Response.Image = img
	}

	Response.Status = true
	Response.Message = "success"
	Response.Data = doc
	return GCFReturnStruct(Response)
}
