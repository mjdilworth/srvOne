package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mjdilworth/srvOne/internal/apploader"
	"github.com/mjdilworth/srvOne/internal/command"
	"github.com/mjdilworth/srvOne/internal/server"
)

var Version = "development"
var CommitId string
var BuildTimestamp string

func main() {

	appFlag := flag.String("flag", "foo", "pass in configuration")
	daemon := flag.Bool("daemon", true, "run as http daemon")
	port := flag.String("port", "8080", "port when running as server daemon")

	flag.Parse()

	fmt.Println("appFlag:\t", *appFlag)
	fmt.Println("daemon:\t", *daemon)
	fmt.Println("Version:\t", Version)
	fmt.Println("CommitId\t", CommitId)
	fmt.Println("BuildTimestamp\t", BuildTimestamp)

	//depending on command line choice
	var app apploader.App
	if *daemon {
		//create the application
		app = server.New(*port) // we want the daemon

	} else {
		app = command.New()
	}
	//create the application

	//load the applications
	al := apploader.New(app)

	//run the application
	if err := al.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//the application has ended
	al.LogMe("Service is ending", "key", 5)

}
