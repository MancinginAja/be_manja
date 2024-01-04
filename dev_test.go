package be_manja

import (
	"fmt"
	"testing"

	model "github.com/MancinginAja/be_manja/model"
	module "github.com/MancinginAja/be_manja/module"
)

var db = module.MongoConnect("MONGOSTRING", "db_manja")
var collectionnameUser = "user"
var collectionnameFishingspot = "fishingspot"

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
	email, err  := module.SignUp(conn, collectionnameUser, user)
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
	user, err := module.LogIn(conn, collectionnameUser, user)
	// user, err := PASETOPRIVATEKEYENV.Decode(os.Getenv("v4.public.eyJleHAiOiIyMDI0LTAxLTA0VDA5OjMxOjUxWiIsImZ1bGxuYW1lIjoiYWRtaW5AZ21haWwuY29tIiwiaWF0IjoiMjAyNC0wMS0wNFQwNzozMTo1MVoiLCJpZCI6IjY1OTY1ZWNkY2MxOGQxNmNkNGNhNGY4YSIsIm5iZiI6IjIwMjQtMDEtMDRUMDc6MzE6NTFaIn2zI79OV1A342NV3eFbJaA38iQClTp07aLSod43pc5B6ZnyPOZLTNHtwfVqyHFBY5QkFYyshAv5LZ6l6wPVqf8I"), user.Token)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Berhasil LogIn : ", user.Fullname)
		fmt.Print("Berhasil LogIn : " + user.Email)
	}
}

// func TestTambahFishingSpot(t *testing.T) {
//     conn := db
//     var spot model.FishingSpot
//     spot.Name = "Spot 1"
//     spot.Phonenumber = "6285718177810"
//     spot.TopFish = "Ikan 1"
//     spot.Rating = "5"
//     spot.OpeningHour = "08:00 - 17:00"
//     spot.Description = "Deskripsi Spot 1"
//     spot.Image = "https://www.google.com"
//     spot.Address = "Alamat Spot 1"
//     spot.Longitude = "0.000000"
//     spot.Latitude = "0.000000"

//     // Perbaikan #1: Memastikan tipe data return dan argumen yang benar
//     _, err := module.PostFishingSpot(conn, collectionnameFishingspot, &spot)
//     if err != nil {
//         fmt.Println(err)
//     } else {
//         fmt.Println("Berhasil TambahFishingSpot : ")
//     }
// }
