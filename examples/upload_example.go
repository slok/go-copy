package main

import (
	"flag"
	"fmt"
	"github.com/slok/go-copy/copy"
	"os"
)

// Set our cmd params
var filePath = flag.String("if", "", "Local file path")
var uploadPath = flag.String("of", "", "Upload file path")

func main() {

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		//panic("File doesn't exist")
		fmt.Println("File doesn't exist")
		os.Exit(-1)
	}

	// Take all the necessary data
	appToken := os.Getenv("APP_TOKEN")
	appSecret := os.Getenv("APP_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// Create the client
	client, _ := copy.NewDefaultClient(appToken, appSecret, accessToken, accessSecret)
	fs := copy.NewFileService(client)
	fmt.Println(fs.UploadFile(*filePath, *uploadPath, true))
	os.Exit(0)
}
