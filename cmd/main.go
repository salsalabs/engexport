package main

import (
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
		config   = app.Flag("tables", `Optional table layout spec.  See "schema.yaml".`).Default("schema.yaml").String()
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
	t, err := engexport.LoadSchema(config)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	var e *engexport.E
	switch args {
	case "groups":
		e = engexport.NewGroups(api, t, *outDir)
	case "donations":
		e = engexport.NewDonation(api, t, *outDir)
	case "supporters all":
		e = engexport.NewAllSupporters(api, t, *outDir)
	case "supporters active":
		e = engexport.NewActiveSupporter(api, t, *outDir)
	case "supporters inactive all":
		e = engexport.NewInactiveSupporter(api, t, *outDir)
	case "supporters inactive donors":
		e = engexport.NewInactiveDonors(api, t, *outDir)
	}
	if e == nil {
		log.Println("Error: you *must* choose a table to export!")
		return
	}

	e.Run(engexport.Threads, *start)
	log.Println("main: done")
	return
}
