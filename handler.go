package be_manja

import (
	"encoding/json"
	"net/http"
	"os"

	model "github.com/MancinginAja/be_manja/model"
	module "github.com/MancinginAja/be_manja/module"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	credential model.Credential
	response   model.Response
	user	   model.User
	password   model.UpdatePassword
)

func SignUpHandler(MONGOCONNSTRINGENV, dbname string, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return module.GCFReturnStruct(response)
	}
	email, err := module.SignUp(conn, collectionname, user)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Berhasil SignUp"
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data": bson.M{
			"email": email,
		},
	}
	return module.GCFReturnStruct(responData)
}

func LogInHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return module.GCFReturnStruct(response)
	}
	user, err := module.LogIn(conn, collectionname, user)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	tokenstring, err := module.Encode(user.ID, user.Email, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		response.Message = "Gagal Encode Token : " + err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	credential.Message = "Selamat Datang " + user.Fullname
	credential.Token = tokenstring
	credential.Status = 200
	responData := bson.M{
		"status":  credential.Status,
		"message": credential.Message,
		"data": bson.M{
			"token": credential.Token,
			"email": user.Email,
		},
	}
	return module.GCFReturnStruct(responData)
}

func GetProfileHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	payload, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	user, err := module.GetUserFromID(payload.Id, conn)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Get Success"
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data": bson.M{
			"_id":          user.ID,
			"nama_lengkap": user.Fullname,
			"email":        user.Email,
			"phonenumber":        user.PhoneNumber,
		},
	}
	return module.GCFReturnStruct(responData)
}

func EditProfileHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = "Gagal Decode Token : " + err.Error()
		return module.GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return module.GCFReturnStruct(response)
	}
	data, err := module.EditProfile(user.Id, conn, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Berhasil mengubah profile" + user.Email
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data":    data,
	}
	return module.GCFReturnStruct(responData)
}

func EditPasswordHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = "Gagal Decode Token : " + err.Error()
		return module.GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return module.GCFReturnStruct(response)
	}
	data, err := module.EditPassword(user.Id, conn, password)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Berhasil mengubah password" + user.Email
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data":    data,
	}
	return module.GCFReturnStruct(responData)
}

func EditEmailHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user_login, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = "Gagal Decode Token : " + err.Error()
		return module.GCFReturnStruct(response)
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = "error parsing application/json: " + err.Error()
		return module.GCFReturnStruct(response)
	}
	data, err := module.EditEmail(user_login.Id, conn, user)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Berhasil mengubah email" + user_login.Email
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data":    data,
	}
	return module.GCFReturnStruct(responData)
}


