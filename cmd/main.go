package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rakyll/statik/fs"
	engexport "github.com/salsalabs/engexport/pkg"
	_ "github.com/salsalabs/engexport/statik"
	godig "github.com/salsalabs/godig/pkg"
)

//customFields adds custom fields to the supporter section of a schema.
func customFields(api *godig.API, schema *engexport.Schema) error {
	t := api.Supporter()
	a, err := t.Describe()
	if err != nil {
		return err
	}
	for _, f := range a {
		if f.IsCustom && len(f.DisplayName) != 0 {
			schema.Supporter.Fields[f.DisplayName] = f.Name
			schema.Supporter.Headers = append(schema.Supporter.Headers, f.DisplayName)
		}
	}
	return nil
}

//dumpSchema writes the generated schema to a disk file.
func dumpSchema(run *engexport.RunConfig, schema engexport.Schema) {
	if run.DumpSchema {
		b, _ := json.MarshalIndent(schema, "", "  ")
		err := ioutil.WriteFile("generated_schema.yaml", b, os.ModePerm)
		if err != nil {
			log.Printf("Warning: unable to write generated schema to ./generated_schema.yaml")
		} else {
			log.Println("Modified schema written to ./generated_schema.yaml")
		}
	}
}

//authenticate uses a RunConfig to authenticate with Salsa Classic's API.
func authenticate(run *engexport.RunConfig) (a *godig.API, err error) {
	c := godig.CredData{
		Host:     run.Host,
		Email:    run.Email,
		Password: run.Password,
	}
	a = godig.NewAPI()
	err = a.Authenticate(c)
	if err == nil {
		a.Verbose = run.APIVerbose
	}
	return a, err
}

//parseRunYAML parses the contents of "run.yaml" (if it exists) and returns
//a RunConfig object.
func parseRunYAML() (run *engexport.RunConfig, err error) {
	if engexport.FileExists("run.yaml") {
		run, err = engexport.LoadRun("run.yaml")
	}
	return run, err
}

//schemaReader parses the "schema" argument in a RunConfig.  It returns a Reader
//to either built-in schemas or to an external schema file.
func selectSchema(run *engexport.RunConfig) (r io.Reader, err error) {
	//The schema can be "engage", "goodbye" or a schema filename.
	//We retrieve a file from the statik object if the schema is
	//engage or goodbye.
	statikFS, err := fs.New()
	if err != nil {
		err = fmt.Errorf("Unable to open statik file system, error is %v", err)
		return r, err
	}
	switch run.Schema {
	case "engage":
		r, err = statikFS.Open("/engage_schema.yaml")
	case "goodbye":
		r, err = statikFS.Open("/goodbye_schema.yaml")
	default:
		r, err = os.Open(run.Schema)
	}
	if err != nil {
		err = fmt.Errorf("Unable to find schema for '%v', %v", run.Schema, err)
	}
	return r, err
}

//main is the starting point for this app.
func main() {
	run, err := parseRunYAML()
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	api, err := authenticate(run)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	r, err := selectSchema(run)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	t, err := engexport.LoadSchema(r)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}

	err = customFields(api, t)
	if err != nil {
		log.Fatalf("Main: %v\n", err)
	}
	dumpSchema(run, *t)
	p := engexport.P{
		API:            api,
		T:              *t,
		Tag:            run.Tag,
		Dir:            run.OutDir,
		DisableInclude: run.DisableInclude,
	}

	for _, arg := range run.Args {
		fmt.Printf("EngExport processing '%v'\n", arg)
		var e *engexport.E
		switch arg {
		case "supporters all":
			e = engexport.NewAllSupporters(p)
		case "supporters active":
			e = engexport.NewActiveSupporter(p)
		case "supporters inactive all":
			e = engexport.NewInactiveSupporter(p)
		case "supporters inactive donors":
			e = engexport.NewInactiveDonors(p)
		case "supporters only_email":
			e = engexport.NewSupporter(p)
		case "supporters no_email":
			e = engexport.NewNoEmailSupporter(p)

		case "groups active":
			e = engexport.NewActiveGroups(p)
		case "groups all":
			e = engexport.NewAllGroups(p)
		case "groups only_email":
			e = engexport.NewEmailOnlyGroups(p)

		case "donations active":
			e = engexport.NewDonation(p)
		case "donations all":
			e = engexport.NewAllDonations(p)

		case "tags":
			e = engexport.NewTagGroups(p)
		case "actions":
			e = engexport.NewAllActions(p)
		case "events":
			e = engexport.NewAllEvents(p)
		case "contact_history":
			e = engexport.NewContactHistory(p)
		case "email_statistics":
			e = engexport.NewEmailStatistics(p)
		case "blast_statistics":
			e = engexport.NewBlastStatistics(p)
		}
		if e == nil {
			log.Println("Error: you *must* choose a table to export!")
			return
		}
		e.Run(engexport.Threads, run.Start)
	}
	log.Println("main: done")
	return
}
