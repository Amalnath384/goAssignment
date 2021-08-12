package repository

import (
	"context"
	"fmt"
	"project1/assignment/pkg/model"
	"project1/assignment/pkg/repository/driver/mongo"
)

type Repository interface {

	CreateUser(ctx context.Context, User *model.Credentials) (*model.Credentials, error)
	GetUser(ctx context.Context, id string) (*model.Credentials, error)
	CreateStudent(ctx context.Context, Student *model.StudentDetails) (*model.StudentDetails, error)
	GetStudent(ctx context.Context, id string) (*model.StudentDetails, error)
	ListStudent(ctx context.Context) ([]*model.StudentDetails, error)
	Close()
}

var Repo Repository

func Init(db *model.Database) {
	switch db.Driver {
	case "etcd":
		fmt.Printf("etcd is not implemented right now!")
		return
	case "mongodb":
		Repo = mongo.Init(db.Endpoint)
		return
	default:
		fmt.Printf("Can't find database driver %s!\n", db.Driver)
	}
}

func Exit() {
	Repo.Close()
}
