package engexport

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/salsalabs/godig"
	"gopkg.in/yaml.v2"
)

const queueSize = 100

//NewDonation instantiates an environment for copying donations to CSV files.
//TODO: Allow a user to iverride these selections with a YAML file.
func NewDonation(api *godig.API, t Tables, dir string) *E {
	c := []string{
		"RESULT IN 0,-1",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%"}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         t.Donation.Fields,
		Headers:        t.Donation.Headers,
		Conditions:     c,
		CsvFilename:    "donations.csv",
		TableName:      "donation(supporter_KEY)supporter",
		CountTableName: "donation",
		OffsetChan:     make(chan int32, queueSize),
		RecordChan:     make(chan R, queueSize),
		DoneChan:       make(chan bool),
	}
	return &e
}

//NewGroups instantiates an environment for copying Gruups and Emails
//t CSV files.
func NewGroups(api *godig.API, t Tables, dir string) *E {
	c := []string{
		"groups.Group_Name IS NOT EMPTY",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%",
		"supporter.Receive_Email>0",
	}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         t.Groups.Fields,
		Headers:        t.Groups.Headers,
		Conditions:     c,
		CsvFilename:    "groups.csv",
		TableName:      "groups(groups_KEY)supporter_groups(supporter_KEY)supporter",
		CountTableName: "supporter_groups",
		OffsetChan:     make(chan int32, queueSize),
		RecordChan:     make(chan R, queueSize),
		DoneChan:       make(chan bool, queueSize),
	}
	return &e
}

//NewSupporter instantiates an environment for copying supporters to  CSV files.
//The default behavior is to save suupporters that have valid email addresses.
//That means that both subscribed and unsubscribed supporrters are written to CSV
//files.  TODO: Allow a user to iverride these selections with a YAML file.
func NewSupporter(api *godig.API, t Tables, dir string) *E {
	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
	}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         t.Supporter.Fields,
		Headers:        t.Supporter.Headers,
		Conditions:     c,
		CsvFilename:    "supporters.csv",
		TableName:      "supporter",
		CountTableName: "supporter",
		OffsetChan:     make(chan int32, queueSize),
		RecordChan:     make(chan R, queueSize),
		DoneChan:       make(chan bool, queueSize),
	}
	log.Printf("R is %+v\n", e.Fields)
	// Just a reminder...
	e.API.Verbose = false
	return &e
}

//NewActiveSupporter instantiates an envionrment to copy active supportes to
//CSV files.  Active supporters have a good email address and have not opted out
//or been opted out (i.e. Receive_Email > 0).
func NewActiveSupporter(api *godig.API, t Tables, dir string) *E {
	e := NewSupporter(api, t, dir)
	c := e.Conditions
	c = append(c, "Receive_Email>0")
	e.Conditions = c
	return e
}

//NewAllSupporters instantiates an environment to copy all supporter to
//CSV files.  Not a good idea for Engage, but useful for other vendors.
func NewAllSupporters(api *godig.API, t Tables, dir string) *E {
	e := NewSupporter(api, t, dir)
	e.Conditions = []string{
		"supporter_KEY>0",
	}
	return e
}

//NewInactiveSupporter instantiates an envionrment to copy inactive supporters to
//CSV files.  Inactive supporters have a good email address but have either opted
// out or been opted out (i.e. Receive_Email < 1).
func NewInactiveSupporter(api *godig.API, t Tables, dir string) *E {
	e := NewSupporter(api, t, dir)
	c := e.Conditions
	c = append(c, "Receive_Email<1")
	e.Conditions = c
	e.CsvFilename = "inactive_" + e.CsvFilename
	return e
}

//NewInactiveDonors instantiates an envionrment to copy inactive supporters with
//donation history to CSV files.  Inctive supporters have a good email address but
//have either opted out or been opted out (i.e. Receive_Email < 1).  Processing uses
//the donation table as a guide to find inactive supporters.
func NewInactiveDonors(api *godig.API, t Tables, dir string) *E {
	e := NewSupporter(api, t, dir)

	//Salsa gets confused with the contents of the "&uinclude=" in selective places.
	//This is one of those times.
	e.DisableInclude = true

	//The table is a join between donations and supporters.
	e.TableName = "supporter(supporter_KEY)donation"

	//The output filename needs to change to protect existing supporter records.
	//Not really.  It just makes accounting easier.
	e.CsvFilename = "inactive_donors.csv"

	//We are counting donations.
	e.CountTableName = "donation"

	//Add a  condition that looks for inactive supporters.
	c := e.Conditions
	c = append(c, "supporter.Receive_Email<1")

	//Add a condition that looks for valid donations.
	c = append(c, "donation.RESULT IN 0,-1")

	e.Conditions = c
	return e
}

//LoadTables accepts a pointer to a filename and loads the contents
//into a Table object.
func LoadTables(fn *string) (Tables, error) {
	var t Tables
	if fn == nil {
		panic(errors.New("a config file is required"))
	}
	b, err := ioutil.ReadFile(*fn)
	if err != nil {
		return t, err
	}
	err = yaml.Unmarshal(b, &t)
	return t, err
}
