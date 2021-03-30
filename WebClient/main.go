package main

import (
	"WebClient/includes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	routing "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	router := routing.New()

	router.GET("/attack", includes.Attack)
	router.GET("/reload", includes.Reload)
	router.GET("/methods", includes.GetMethods)
	//router.GET("/index", includes.Index)
	// Static files
	fs := fasthttp.FSHandler("./html", 0)
	router.GET("/html", fs)

	switch len(os.Args) {
	case 1:
		fmt.Printf("Invalid usage!\n\rUsage: %s {HOST} {PORT} {API Key}\n", os.Args[0])
	case 2:
		ready(true)
		fasthttp.ListenAndServe(":"+os.Args[1], router.Handler)
	case 3:
		ready(true)
		fasthttp.ListenAndServe(os.Args[1]+":"+os.Args[2], router.Handler)

	case 4:
		ready(false)
		fasthttp.ListenAndServe(os.Args[1]+":"+os.Args[2], router.Handler)
	}
}

func ready(setKey bool) {

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error occured while opening: config.json.\n", err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(data), &includes.Config)
	fmt.Println("Webserver ready to serve!")

	if setKey {
		includes.Key = fmt.Sprint(rand.Int())
		fmt.Printf("Your randomly generated API key is: %s\n", includes.Key)
	} else {
		includes.Key = os.Args[3]
		fmt.Printf("Your API key is: %s\n", includes.Key)
	}
}
