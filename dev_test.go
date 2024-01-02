package todolist

import (
	"testing"

	model "github.com/daniferdinandall/be_todolist/model"
	module "github.com/daniferdinandall/be_todolist/module"
)

var db = module.MongoConnect("MONGOSTRING", "db_todolist")

func TestInsertDoc(t *testing.T) {
	var doc model.Todolist
	doc.Title = "Test"
	doc.Description = "Test"
	doc.DueDate = "Test"
	doc.Priority = 1
	doc.Completed = false

	module.InsertDoc(db, "todolist", doc)
}
