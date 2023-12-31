package be_manja

import (
	"encoding/json"
	"net/http"
	"os"

	model "github.com/MancinginAja/be_manja/model"
	module "github.com/MancinginAja/be_manja/module"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TambahFishingSpotHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses, email anda : " + user.Email
		return module.GCFReturnStruct(response)
	}
	data, err := module.PostFishingSpot(conn, collectionname, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 201
	response.Message = "Berhasil menambah Fishing Spot"
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data":    data,
	}
	return module.GCFReturnStruct(responData)
}

func GetFishingSpotHandler(MONGOCONNSTRINGENV, dbname string, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	id := module.GetID(r)
	if id == "" {
		data, err := module.GetAllFishingSpot(conn, collectionname)
		if err != nil {
			response.Message = err.Error()
			return module.GCFReturnStruct(response)
		}
		responData := bson.M{
			"status":  200,
			"message": "Get Success",
			"data":    data,
		}
		//
		return module.GCFReturnStruct(responData)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	fishingSpot, err := module.GetFishingSpotById(conn, collectionname, idparam)
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
		"data":    fishingSpot,
	}
	return module.GCFReturnStruct(responData)
}

func EditFishingSpotHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses"
		return module.GCFReturnStruct(response)
	}
	id := module.GetID(r)
	if id == "" {
		response.Message = "Wrong parameter"
		return module.GCFReturnStruct(response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = "Invalid id parameter"
		return module.GCFReturnStruct(response)
	}
	data, err := module.PutFishingSpot(idparam, conn, collectionname, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 200
	response.Message = "Berhasil mengubah Fishing Spot"
	responData := bson.M{
		"status":  response.Status,
		"message": response.Message,
		"data":    data,
	}
	return module.GCFReturnStruct(responData)
}

func DeleteFishingSpotHandler(PASETOPUBLICKEYENV, MONGOCONNSTRINGENV, dbname string, collectionname string, r *http.Request) string {
	conn := module.MongoConnect(MONGOCONNSTRINGENV, dbname)
	response.Status = 400
	//
	user, err := module.GetUserLogin(PASETOPUBLICKEYENV, r)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	if user.Email != "admin@gmail.com" {
		response.Message = "Anda tidak memiliki akses"
		return module.GCFReturnStruct(response)
	}
	id := module.GetID(r)
	if id == "" {
		response.Message = "Wrong parameter"
		return module.GCFReturnStruct(response)
	}
	idparam, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response.Message = "Invalid id parameter"
		return module.GCFReturnStruct(response)
	}
	err = module.DeleteFishingspot(idparam, collectionname, conn)
	if err != nil {
		response.Message = err.Error()
		return module.GCFReturnStruct(response)
	}
	//
	response.Status = 204
	response.Message = "Berhasil menghapus Fishing Spot"
	return module.GCFReturnStruct(response)
}
