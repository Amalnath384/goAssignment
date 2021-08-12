package model

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/go-playground/validator/v10"
	"io"
)

type StudentDetails struct {
	Id     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Token  string        `json:"token"`
	Name   string        `json:"name" validate:"required"`
	Age    int32         `json:"age" validate:"required"`
	Class  int32         `json:"class,omitempty"`
	Marks *Marks `json:"marks" bson:"marks,omitempty"`
}

type Marks struct {
	Golang float64 `json:"golang,omitempty" bson:"golang,omitempty"`
	Java  float64 `json:"java,omitempty" bson:"java,omitempty"`
}


func (p *StudentDetails) DecodeFromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *StudentDetails) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *StudentDetails) EncodeToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

type Credentials struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username"`
	Password string        `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
