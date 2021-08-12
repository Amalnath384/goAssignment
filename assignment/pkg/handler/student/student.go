package student

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project1/assignment/pkg/model"
	"project1/assignment/pkg/repository"
	"time"
)

type StudentService struct {
	repo repository.Repository
}

func NewStudentService() *StudentService {

	return &StudentService{}
}

var jwtKey = []byte("secret_key")


func (p *StudentService) ListStudent(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	res, err := p.repo.ListStudent(ctx)
	if err != nil {
		http.Error(w, "Failed to create in database", http.StatusBadRequest)
		return
	}
	//var Student model.StudentDetails
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

func (p *StudentService) GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	params := mux.Vars(r)
	id := params["id"]
	res, err := repository.Repo.GetStudent(ctx, id)
	if err != nil {
		http.Error(w, "Failed to create in database", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (p *StudentService) CreateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Student model.StudentDetails
	err := Student.DecodeFromJSON(r.Body)
	if err != nil {
		http.Error(w, "Failed to Decode", http.StatusBadRequest)
		return
	}
	if err = Student.Validate(); err != nil {
		http.Error(w, "failed to validate struct", http.StatusBadRequest)
		return
	}
	res := TokenValidation(Student.Token)

	if res == true {
		ctx := context.Background()

		res, err := repository.Repo.CreateStudent(ctx, &Student)
		if err != nil {
			http.Error(w, "Failed to create in database", http.StatusBadRequest)
			return
		}
		fmt.Println(res)

		w.WriteHeader(http.StatusCreated)
		Student.EncodeToJSON(w)
		return

	} else {
		fmt.Fprintln(w, "not validated")
	}

}


func TokenValidation(token string) bool {
	claim := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claim,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("signature is not valid")
			return false
		}
		log.Println("error while parsing request")
		return false
	}

	if !tkn.Valid {
		log.Println("token is not valid")
		return false
	}

	if tkn.Valid {
		return true
	}
	return true
}


func (p *StudentService) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var User model.Credentials
	json.NewDecoder(r.Body).Decode(&User)

	ctx := context.Background()

	res, err := repository.Repo.CreateUser(ctx, &User)
	if err != nil {
		http.Error(w, "Failed to create in database", http.StatusBadRequest)
		return
	}
	fmt.Println(res)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(User)
	return

}

func (p *StudentService) CreateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//

	var credentials model.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	params := mux.Vars(r)
	id := params["id"]
	res, err := repository.Repo.GetUser(ctx, id)

	if credentials.Username == res.Username {
		if credentials.Password == res.Password {
			fmt.Fprintln(w, "successfully logged in")

		} else {
			fmt.Fprintln(w, "wrong password")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Fprintln(w, "wrong username")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &model.Claims{
		Username: res.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("access-control-expose-headers", tokenString)

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	fmt.Fprintln(w, tokenString)

}
