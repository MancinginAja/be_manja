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

// func TestLogIn(t *testing.T) {
// 	conn := db
// 	var user model.User
// 	user.Email = "dito@gmail.com"
// 	user.Password = "dito12345678"
// 	user, _ = module.LogIn(conn, collectionnameUser, user)
// 	token := "v4.public.eyJleHAiOiIyMDI0LTAxLTA0VDEwOjQ2OjQwWiIsImZ1bGxuYW1lIjoiYWRtaW5AZ21haWwuY29tIiwiaWF0IjoiMjAyNC0wMS0wNFQwODo0Njo0MFoiLCJpZCI6IjY1OTY1ZWNkY2MxOGQxNmNkNGNhNGY4YSIsIm5iZiI6IjIwMjQtMDEtMDRUMDg6NDY6NDBaIn2ABnDsInI3fCmAyBB6FUCzcTxhgQxUyz1SnNl5vl_AzTqvIRe193Hh2t4M9I8IPPn7LczZhODquCwa57tqOCEB"
// 	tokenstring, err := module.Decode("2df976efd75eacf2a37b1ce184deec8d3b72cb78f7881ed5e7a02d97351c2aef", token)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Berhasil LogIn : ", user.Fullname)
// 		fmt.Print("Berhasil LogIn : " + tokenstring.Id.Hex())
// 	}
// }

func TestLogInn(t *testing.T) {
	conn := db
	var user model.User
	user.Email = "admin@gmail.com"
	user.Password = "admin12345678"
	user, _ = module.LogIn(conn, collectionnameUser, user)
	tokenstring, err := module.Encode(user.ID, user.Email, "33186fcfc13ba9946bf200cf6c7808e6ebfc605140f65809e06648985b08ebda2df976efd75eacf2a37b1ce184deec8d3b72cb78f7881ed5e7a02d97351c2aef")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Berhasil LogIn : ", user.Fullname)
		fmt.Print("Berhasil LogIn : " + tokenstring)
	}
}

func TestToken (*testing.T) {
	token := "v4.public.eyJleHAiOiIyMDI0LTAxLTA0VDExOjI1OjU0WiIsImZ1bGxuYW1lIjoiYWRtaW5AZ21haWwuY29tIiwiaWF0IjoiMjAyNC0wMS0wNFQwOToyNTo1NFoiLCJpZCI6IjY1OTY1ZWNkY2MxOGQxNmNkNGNhNGY4YSIsIm5iZiI6IjIwMjQtMDEtMDRUMDk6MjU6NTRaIn22kA21UMcQv-6lNrkBu88rV3XGGgToTBqulQui3HrZcYb_Go-qyCBdzje7Qg3Omj-hI5lXRRFj1afCzeMdyG0B"
	tokenstring, err := module.Decode("2df976efd75eacf2a37b1ce184deec8d3b72cb78f7881ed5e7a02d97351c2aef", token)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print("Id Token : " + tokenstring.Id.Hex())
		fmt.Print("Email Token : " + tokenstring.Email)
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
