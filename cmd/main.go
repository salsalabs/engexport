package main

import (
	"fmt"
	"log"
	"os"

	"github.com/salsalabs/engexport"
	"github.com/salsalabs/godig"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app      = kingpin.New("engexport", "Classic-to-Engage exporter.")
		login    = app.Flag("login", "YAML file with login credentials").Required().String()
		outDir   = app.Flag("dir", "Directory to use to store results").Default("./data").String()
		start    = app.Flag("start", "start processing at this offset").Default("0").Int32()
		supp     = app.Command("supporters", "process supporters")
		_        = supp.Command("all", "process all supporters")
		_        = supp.Command("active", "process active supporters")
		inactive = supp.Command("inactive", "process inactive supporters")
		_        = inactive.Command("all", "process all inactive supporters")
		_        = inactive.Command("donors", "process inactive supporters with donation history")
		_        = app.Command("groups", "process groups for active supporters")
		_        = app.Command("donations", "process donations for active and inactive supporters")
	)
	args, _ := app.Parse(os.Args[1:])
	api, err := (godig.YAMLAuth(*login))
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	var e *engexport.E
	switch args {
	case "groups":
		e = engexport.NewGroups(api, *outDir)
	case "donations":
		e = engexport.NewDonation(api, *outDir)
	case "supporters all":
		e = engexport.NewSupporter(api, *outDir)
	case "supporters active":
		e = engexport.NewActiveSupporter(api, *outDir)
	case "supporters inactive all":
		e = engexport.NewInactiveSupporter(api, *outDir)
	case "supporters inactive donors":
		e = engexport.NewInactiveDonors(api, *outDir)
	}
	if e == nil {
		fmt.Println("Error: you *must* choose a table to export!")
		return
	}

	e.Run(engexport.Threads, *start)
	fmt.Println("main: done")
	return
}
