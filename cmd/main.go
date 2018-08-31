package main

import (
	"fmt"
	"log"

	"github.com/salsalabs/engexport"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app    = kingpin.New("qdd", "Quick and dirty donation scraper.")
		login  = app.Flag("login", "YAML file with login credentials").Required().String()
		outDir = app.Flag("dir", "Directory to use to store results").Default("./data").String()
		start  = app.Flag("start", "start processing at this offset").Default("0").Int32()
		supp   = app.Command("supporters", "process supporters")
		_      = supp.Command("all", "process all supporters")
		_      = supp.Command("active", "process active supporters")
		_      = supp.Command("inactive", "process active supporters")
		_      = app.Command("groups", "process groups")
		_      = app.Command("donations", "process donations")
	)
	args, _ := app.Parse(os.Args[1:])
	api, err := (godig.YAMLAuth(*login))
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	var e *E
	switch args {
	case "groups":
		e = engexport.NewGroups(api, *outDir)
	case "donations":
		e = engexport.NewDonation(api, *outDir)
	case "supporters all":
		e = engexport.NewSupporter(api, *outDir)
	case "supporters active":
		e = engexport.NewActiveSupporter(api, *outDir)
	case "supporters inactive":
		e = engexport.NewInActiveSupporter(api, *outDir)
	}
	if e == nil {
		fmt.Println("Error: you *must* choose a table to engexport!")
		return
	}

	err = e.Run(Threads, *start)
	if err != nil {
		panic(err)
	}
}
