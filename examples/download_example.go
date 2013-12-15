package main

import (
	"flag"
	"github.com/slok/go-copy/copy"
	"io/ioutil"
	"os"
)

// Set our cmd params
var downloadPath = flag.String("if", "", "Download path")
var writePath = flag.String("of", "", "File write path")

func main() {

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	// Take all the necessary data
	appToken := os.Getenv("APP_TOKEN")
	appSecret := os.Getenv("APP_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	client, _ := copy.NewDefaultClient(appToken, appSecret, accessToken, accessSecret)
	fs := copy.NewFileService(client)

	r, _ := fs.GetFile(*downloadPath)
	fileBytes, _ := ioutil.ReadAll(r)

	err := ioutil.WriteFile(*writePath, fileBytes, 0644)
	if err != nil {
		panic(err)
	}
}
