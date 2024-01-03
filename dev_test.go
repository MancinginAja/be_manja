package be_manja

import (
	"fmt"
	"testing"

	model "github.com/MancinginAja/be_manja/model"
	module "github.com/MancinginAja/be_manja/module"
)

var db = module.MongoConnect("MONGOSTRING", "db_manja")
var collectionname = "user"

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey := module.GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)
}

func TestSignUp(t *testing.T) {
	conn := db
	var user model.User
	user.Fullname = "dito"
	user.Email = "dito@gmail.com"
	user.Password = "dito12345678"
	user.PhoneNumber = "6285718177810"
	email, err  := module.SignUp(conn, collectionname, user)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Berhasil SignUp : ", email)
	}
}

func TestLogIn(t *testing.T) {
	conn := db
	var user model.User
	user.Email = "dito@gmail.com"
	user.Password = "dito12345678"
	user, err := module.LogIn(conn, collectionname, user)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Berhasil LogIn : ", user.Fullname)
	}
}

