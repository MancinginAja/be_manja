package be_manja

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"

	model "github.com/MancinginAja/be_manja/model"
	intermoni "github.com/intern-monitoring/backend-intermoni"
)

var imageUrl string

// mongo
func MongoConnect(MongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(MongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}
	return client.Database(dbname)
}

// crud
func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error GetAllDocs %s: %s", col, err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return err
	}
	return docs
}

func InsertOneDoc(db *mongo.Database, col string, doc interface{}) (insertedID primitive.ObjectID, err error) {
	result, err := db.Collection(col).InsertOne(context.Background(), doc)
	if err != nil {
		return insertedID, fmt.Errorf("kesalahan server : insert")
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func UpdateOneDoc(id primitive.ObjectID, db *mongo.Database, col string, doc interface{}) (err error) {
	filter := bson.M{"_id": id}
	result, err := db.Collection(col).UpdateOne(context.Background(), filter, bson.M{"$set": doc})
	if err != nil {
		return fmt.Errorf("error update: %v", err)
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("tidak ada data yang diubah")
		return
	}
	return nil
}

func DeleteOneDoc(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}

// get user login
func GetUserLogin(PASETOPUBLICKEYENV string, r *http.Request) (model.Payload, error) {
	tokenstring := r.Header.Get("Authorization")
	payload, err := Decode(os.Getenv(PASETOPUBLICKEYENV), tokenstring)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// get id
func GetID(r *http.Request) string {
	return r.URL.Query().Get("id")
}

// return struct
func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// validate phonenumber
func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	// Define the regular expression pattern for numeric characters
	numericPattern := `^[0-9]+$`

	// Compile the numeric pattern
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	// Check if the phone number consists only of numeric characters
	if !numericRegexp.MatchString(phoneNumber) {
		return false, nil
	}

	// Define the regular expression pattern for "62" followed by 6 to 12 digits
	pattern := `^62\d{6,13}$`

	// Compile the regular expression
	regexpPattern, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	
	// Test if the phone number matches the pattern
	isValid := regexpPattern.MatchString(phoneNumber)

	return isValid, nil
}

// validate latitude longitude
func CheckLatitudeLongitude(db *mongo.Database, latitude, longitude string) bool {
	collection := db.Collection("billboard")
	filter := bson.M{"latitude": latitude, "longitude": longitude}
	err := collection.FindOne(context.Background(), filter).Decode(&model.FishingSpot{})
	return err == nil
}

// user
// get-user
func GetUserFromID(_id primitive.ObjectID, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"_id": _id}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return doc, fmt.Errorf("no data found for ID %s", _id)
		}
		return doc, fmt.Errorf("error retrieving data for ID %s: %s", _id, err.Error())
	}
	return doc, nil
}

func GetUserFromEmail(email string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"email": email}
	err = collection.FindOne(context.TODO(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("email tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetUserFromPhonenumber(phonenumber string, db *mongo.Database) (doc model.User, err error) {
	collection := db.Collection("user")
	filter := bson.M{"phonenumber": phonenumber}
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("nomor telepon tidak ditemukan")
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

// update-userProfile
func EditProfile(idparam primitive.ObjectID, db *mongo.Database, r *http.Request) (bson.M, error) {
	dataUser, err := GetUserFromID(idparam, db)
	if err != nil {
		return bson.M{}, err
	}
	fullname := r.FormValue("fullname")
	phonenumber := r.FormValue("phonenumber")

	image := r.FormValue("file")

	if fullname == "" || phonenumber == ""{
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}
	if image != "" {
		imageUrl = image
	} else {
		imageUrl, err = intermoni.SaveFileToGithub("erditona", "erditonaushaadam@gmail.com", "image-manja", "manja", r)
		if err != nil {
			return bson.M{}, fmt.Errorf("error save file: %s", err)
		}
		image = imageUrl
	}

	profile := bson.M{
		"fullname":    fullname,
		"email":       dataUser.Email,
		"password":    dataUser.Password,
		"phonenumber": phonenumber,
		"image":       image,
		"salt":        dataUser.Salt,
	}
	err = UpdateOneDoc(idparam, db, "user", profile)
	if err != nil {
		return bson.M{}, err
	}
	data := bson.M{
		"fullname":    fullname,
		"email":       dataUser.Email,	
		"phonenumber": phonenumber,
		"image":       image,
	}

	return data, nil
}

func EditEmail(idparam primitive.ObjectID, db *mongo.Database, insertedDoc model.User) (bson.M, error) {
	dataUser, err := GetUserFromID(idparam, db)
	if err != nil {
		return bson.M{}, err
	}
	if insertedDoc.Email == "" {
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return bson.M{}, fmt.Errorf("email tidak valid")
	}
	existsDoc, _ := GetUserFromEmail(insertedDoc.Email, db)
	if existsDoc.Email == insertedDoc.Email {
		return bson.M{}, fmt.Errorf("email sudah terdaftar")
	}
	user := bson.M{
		"fullname":    dataUser.Fullname,
		"email":       insertedDoc.Email,
		"password":    dataUser.Password,
		"phonenumber": dataUser.PhoneNumber,
		"image":       dataUser.Image,
		"salt":        dataUser.Salt,
	}
	err = UpdateOneDoc(idparam, db, "user", user)
	if err != nil {
		return bson.M{}, err
	}
	data := bson.M{
		"fullname":    dataUser.Fullname,
		"email":       insertedDoc.Email,
		"phonenumber": dataUser.PhoneNumber,
		"image":      dataUser.Image,
	}
	return data, nil
}

func EditPassword(idparam primitive.ObjectID, db *mongo.Database, insertedDoc model.UpdatePassword) (bson.M, error) {
	dataUser, err := GetUserFromID(idparam, db)
	if err != nil {
		return bson.M{}, err
	}
	salt, err := hex.DecodeString(dataUser.Salt)
	if err != nil {
		return bson.M{}, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Oldpassword), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != dataUser.Password {
		return bson.M{}, fmt.Errorf("password lama salah")
	}
	if strings.Contains(insertedDoc.Newpassword, " ") {
		return bson.M{}, fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Newpassword) < 8 {
		return bson.M{}, fmt.Errorf("password terlalu pendek")
	}
	salt = make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return bson.M{}, fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Newpassword), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"fullname":    dataUser.Fullname,
		"email":       dataUser.Email,
		"password":    hex.EncodeToString(hashedPassword),
		"phonenumber": dataUser.PhoneNumber,
		"image":       dataUser.Image,
		"salt":        hex.EncodeToString(salt),
	}
	err = UpdateOneDoc(idparam, db, "user", user)
	if err != nil {
		return bson.M{}, err
	}
	data := bson.M{
		"fullname":    dataUser.Fullname,
		"email":       dataUser.Email,
		"phonenumber": dataUser.PhoneNumber,
		"image":       dataUser.Image,
	}

	return data, nil
}


// user Sign Up
func SignUp(db *mongo.Database, col string, insertedDoc model.User) (string, error) {
	if insertedDoc.Fullname == "" || insertedDoc.Email == "" || insertedDoc.Password == "" || insertedDoc.PhoneNumber == ""{
		return "", fmt.Errorf("mohon untuk melengkapi data")
	}
	if err := checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return "", fmt.Errorf("email tidak valid")
	}
	userExists, _ := GetUserFromEmail(insertedDoc.Email, db)
	if insertedDoc.Email == userExists.Email {
		return "", fmt.Errorf("email sudah terdaftar")
	}
	validatePhoneNumber, _ := ValidatePhoneNumber(insertedDoc.PhoneNumber)
	if !validatePhoneNumber {
		return "", fmt.Errorf("nomor telepon tidak valid")
	}
	PhoneNumberExists, _ := GetUserFromPhonenumber(insertedDoc.PhoneNumber, db)
	if insertedDoc.PhoneNumber == PhoneNumberExists.PhoneNumber {
		return "", fmt.Errorf("nomor telepon sudah terdaftar")
	}
	if strings.Contains(insertedDoc.Password, " ") {
		return "", fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if len(insertedDoc.Password) < 8 {
		return "", fmt.Errorf("password terlalu pendek")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("kesalahan server : salt")
	}
	hashedPassword := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	user := bson.M{
		"fullname": insertedDoc.Fullname,
		"email": insertedDoc.Email,
		"password": hex.EncodeToString(hashedPassword),
		"phonenumber": insertedDoc.PhoneNumber,
		"salt": hex.EncodeToString(salt),
	}
	_, err = InsertOneDoc(db, col, user)
	if err != nil {
		return "", err
	}
	return insertedDoc.Email, nil
}

// user Sign In
func LogIn(db *mongo.Database, col string, insertedDoc model.User) (user model.User, err error) {
	if insertedDoc.Email == "" || insertedDoc.Password == "" {
		return user, fmt.Errorf("mohon untuk melengkapi data")
	}
	if err = checkmail.ValidateFormat(insertedDoc.Email); err != nil {
		return user, fmt.Errorf("email tidak valid")
	}
	existsDoc, err := GetUserFromEmail(insertedDoc.Email, db)
	if err != nil {
		return
	}
	salt, err := hex.DecodeString(existsDoc.Salt)
	if err != nil {
		return user, fmt.Errorf("kesalahan server : salt")
	}
	hash := argon2.IDKey([]byte(insertedDoc.Password), salt, 1, 64*1024, 4, 32)
	if hex.EncodeToString(hash) != existsDoc.Password {
		return user, fmt.Errorf("password salah")
	}
	return existsDoc, nil
}

// Fishing Spot
// post-fishingSpot
func PostFishingSpot(db *mongo.Database, col string,r *http.Request) (bson.M, error) {
	name := r.FormValue("name")
	phonenumber := r.FormValue("phonenumber")
	topfish := r.FormValue("topfish")
	rating := r.FormValue("rating")
	openinghour := r.FormValue("openinghour")
	description := r.FormValue("description")
	address := r.FormValue("address")
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	
	if name == "" || phonenumber == "" || topfish == "" || rating == "" || openinghour == "" || description == "" || address == "" || latitude == "" || longitude == "" {
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}
	if CheckLatitudeLongitude(db, latitude, longitude) {
		return bson.M{}, fmt.Errorf("lokasi sudah terdaftar")
	}
	validatePhoneNumber, _ := ValidatePhoneNumber(phonenumber)
	if !validatePhoneNumber {
		return bson.M{}, fmt.Errorf("nomor telepon tidak valid")
	}

	imageUrl, err := intermoni.SaveFileToGithub("erditona", "erditonaushaadam@gmail.com", "image-manja", "manja", r)
	if err != nil {
		return bson.M{}, fmt.Errorf("error save file: %s", err)
	}

	fishingSpot := bson.M{
		"_id": primitive.NewObjectID(),
		"name": name,
		"phonenumber": phonenumber,
		"topfish": topfish,
		"rating": rating,
		"openinghour": openinghour,
		"description": description,
		"image": imageUrl,
		"address": address,
		"longitude": longitude,
		"latitude": latitude,
	}
	_, err = InsertOneDoc(db, col, fishingSpot)
	if err != nil {
		return bson.M{}, err
	}
	return fishingSpot, nil
}

// get-fishingSpot
func GetFishingSpotById(db *mongo.Database, col string, idparam primitive.ObjectID) (doc model.FishingSpot, err error) {
	collection := db.Collection(col)
	filter := bson.M{"_id": idparam}
	err = collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return doc, fmt.Errorf("data tidak ditemukan untuk ID %s", idparam)
		}
		return doc, fmt.Errorf("kesalahan server")
	}
	return doc, nil
}

func GetAllFishingSpot(db *mongo.Database, col string) (docs []model.FishingSpot, err error) {
	collection := db.Collection(col)
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	err = cursor.All(context.Background(), &docs)
	if err != nil {
		return docs, fmt.Errorf("kesalahan server")
	}
	return docs, nil
}

// put-fishingSpot
func PutFishingSpot(_id primitive.ObjectID, db *mongo.Database, r *http.Request) (bson.M, error) {
	name := r.FormValue("name")
	phonenumber := r.FormValue("phonenumber")
	topfish := r.FormValue("topfish")
	rating := r.FormValue("rating")
	openinghour := r.FormValue("openinghour")
	description := r.FormValue("description")
	address := r.FormValue("address")
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")

	image := r.FormValue("file")

	if name == "" || phonenumber == "" || topfish == "" || rating == "" || openinghour == "" || description == "" || address == "" || latitude == "" || longitude == "" {
		return bson.M{}, fmt.Errorf("mohon untuk melengkapi data")
	}

	if image != "" {
		imageUrl = image
	} else {
		imageUrl, err := intermoni.SaveFileToGithub("erditona", "erditonaushaadam@gmail.com", "image-manja", "manja", r)
		if err != nil {
			return bson.M{}, fmt.Errorf("error save file: %s", err)
		}
		image = imageUrl
	}

	billboard := bson.M{
		"name": name,
		"phonenumber": phonenumber,
		"topfish": topfish,
		"rating": rating,
		"openinghour": openinghour,
		"description": description,
		"image": image,
		"address": address,
		"longitude": longitude,
		"latitude": latitude,
	}
	err := UpdateOneDoc(_id, db, "billboard", billboard)
	if err != nil {
		return bson.M{}, err
	}
	return billboard, nil
}

// delete-fishingSpot
func DeleteFishingSpot(_id primitive.ObjectID, db *mongo.Database, col string) error {
	collection := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s: %s", _id, err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found", _id)
	}

	return nil
}