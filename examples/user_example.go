package main

import (
	"fmt"
	"github.com/slok/go-copy/copy"
	"math"
	"os"
)

func main() {
	//Prepare the neccesary data
	appToken := os.Getenv("APP_TOKEN")
	appSecret := os.Getenv("APP_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// Create the client
	client, err := copy.NewDefaultClient(appToken, appSecret, accessToken, accessSecret)
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not create the client, review the auth params")
		os.Exit(-1)
	}

	//Create the service (in this case for a user)
	userService := copy.NewUserService(client)

	//Play with the lib :)
	user, err := userService.Get()
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not retrieve the user")
		os.Exit(-1)
	}

	byteToMegabyte := math.Pow(1024, 2)
	fmt.Printf("User: %v %v\n", user.FirstName, user.LastName)
	fmt.Printf("Email: %v\n", user.Email)
	fmt.Printf("Stored(MB): %G of %G\n",
		float64(user.Storage.Used)/byteToMegabyte, float64(user.Storage.Quota)/byteToMegabyte)

	// We are going to change the name
	fmt.Println("Inser name: ")
	fmt.Scan(&(user.FirstName))
	fmt.Println("Inser surname: ")
	fmt.Scan(&(user.LastName))

	err = userService.Update(user)
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not update the user")
		os.Exit(-1)
	}

	// Get again the user
	user, err = userService.Get()
	if err != nil {
		fmt.Fprint(os.Stderr, "Could not retrieve the user")
		os.Exit(-1)
	}
	fmt.Printf("User: %v %v\n", user.FirstName, user.LastName)

}
