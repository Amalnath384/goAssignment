package main

import (
	"project1/assignment/pkg/model"
	db "project1/assignment/pkg/repository"
	."project1/assignment/router"
)

func main() {
	dbHost := "127.0.0.1:27017"
	db.Init(&model.Database{
		Driver:   "mongodb",
		Endpoint: dbHost})
	defer db.Exit()

	Router()

}
