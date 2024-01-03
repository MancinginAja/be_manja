package be_manja

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Fullname    string             `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	Image	    string             `json:"image,omitempty" bson:"image,omitempty"`
	Salt        string             `json:"salt,omitempty" bson:"salt,omitempty"`
}

type UpdatePassword struct {
	Oldpassword string `json:"oldpassword,omitempty" bson:"oldpassword,omitempty"`
	Newpassword string `json:"newpassword,omitempty" bson:"newpassword,omitempty"`
}

type FishingSpot struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	TopFish	    string             `json:"topfish,omitempty" bson:"topfish,omitempty"`
	Rating      string             `json:"rating,omitempty" bson:"rating,omitempty"`
	OpeningHour string             `json:"openinghour,omitempty" bson:"openinghour,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	Longitude   string             `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Latitude    string             `json:"latitude,omitempty" bson:"latitude,omitempty"`
}

type Credential struct {
	Status  int    `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Payload struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	Exp   time.Time          `json:"exp"`
	Iat   time.Time          `json:"iat"`
	Nbf   time.Time          `json:"nbf"`
}