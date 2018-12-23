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
		app        = kingpin.New("engexport", "Classic-to-Engage exporter.")
		login      = app.Flag("login", "YAML file with login credentials").Required().String()
		config     = app.Flag("schema", `Classic table schema.`).Default("schema.yaml").String()
		outDir     = app.Flag("dir", "Directory to use to store results").Default("./data").String()
		tag        = app.Flag("tag", "Retrieve records tagged with TAG").String()
		start      = app.Flag("start", "start processing at this offset").Default("0").Int32()
		apiVerbose = app.Flag("apiVerbose", "each api call and response is displayed if true").Default("false").Bool()
		supp       = app.Command("supporters", "process supporters")
		_          = supp.Command("all", "process all supporters")
		_          = supp.Command("active", "process active supporters")
		_          = supp.Command("only_email", "process supporters that have emails")
		_          = supp.Command("subscribed", "process subscribed emails without regard to emails")
		inactive   = supp.Command("inactive", "process inactive supporters")
		_          = inactive.Command("all", "process all inactive supporters")
		_          = inactive.Command("donors", "process inactive supporters with donation history")
		groups     = app.Command("groups", "process groups")
		_          = groups.Command("active", "process groups for active supporters")
		_          = groups.Command("only_email", "process groups for supporters that have emails only")
		_          = groups.Command("all", "process groups for all supporters, even ones without emails")
		_          = groups.Command("subscribed", "process groups for subscribed supporters")
		donations  = app.Command("donations", "process donations ")
		_          = donations.Command("standard", "process donations for active and inactive supporters")
		_          = donations.Command("subscribed", "process donations for subscribed supporters")
		_          = app.Command("tags", "process tags as groups")
		_          = app.Command("actions", "process supporters and actions")
		_          = app.Command("events", "process supporters and events")
		_          = app.Command("contact_history", "contact history for all supporters")
	)
	args, _ := app.Parse(os.Args[1:])
	if tag != nil && len(*tag) == 0 {
		tag = nil
	}
	api, err := (godig.YAMLAuth(*login))
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	api.Verbose = *apiVerbose
	t, err := engexport.LoadSchema(config)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	p := engexport.P{
		API: api,
		T:   t,
		Tag: tag,
		Dir: *outDir,
	}
	var e *engexport.E
	switch args {
	case "groups active":
		e = engexport.NewActiveGroups(p)
	case "groups all":
		e = engexport.NewAllGroups(p)
	case "groups only_email":
		e = engexport.NewEmailOnlyGroups(p)
	case "groups subscribed":
		e = engexport.NewSubscribedGroups(p)
	case "donations standard":
		e = engexport.NewDonation(p)
	case "donations subscribed":
		e = engexport.NewSubscribedDonation(p)
	case "supporters all":
		e = engexport.NewAllSupporters(p)
	case "supporters active":
		e = engexport.NewActiveSupporter(p)
	case "supporters only_email":
		e = engexport.NewSupporter(p)
	case "supporters subscribed":
		e = engexport.NewSubscribedSupporter(p)
	case "supporters inactive all":
		e = engexport.NewInactiveSupporter(p)
	case "supporters inactive donors":
		e = engexport.NewInactiveDonors(p)
	case "tags":
		e = engexport.NewTagGroups(p)
	case "actions":
		e = engexport.NewAllActions(p)
	case "events":
		e = engexport.NewAllEvents(p)
	case "contact_history":
		e = engexport.NewContactHistory(p)
	}
	if e == nil {
		log.Println("Error: you *must* choose a table to export!")
		return
	}

	e.Run(engexport.Threads, *start)
	log.Println("main: done")
	return
}
