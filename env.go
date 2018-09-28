package engexport

import (
	"io/ioutil"
	"log"
	"sort"

	"github.com/salsalabs/godig"
	"gopkg.in/yaml.v2"
)

const queueSize = 100

//NewDonation instantiates an environment for copying donations to CSV files.
//TODO: Allow a user to iverride these selections with a YAML file.
func NewDonation(api *godig.API, dir string) *E {
	f := R{
		"supporter_KEY":    "supporter.supporter_KEY",
		"Email":            "supporter.Email",
		"donation_KEY":     "donation_KEY",
		"Transaction_Date": "Transaction_Date",
		"Amount":           "amount",
		"Transaction_Type": "Transaction_Type",
		"RESULT":           "RESULT",
	}
	c := []string{
		"RESULT IN 0,-1",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%"}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         f,
		Headers:        headers(f),
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
func NewGroups(api *godig.API, dir string) *E {
	f := R{
		"Group": "Group_Name",
		"Email": "supporter.Email",
	}
	c := []string{
		"groups.Group_Name IS NOT EMPTY",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%",
		"supporter.Receive_Email>0",
	}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         f,
		Headers:        []string{"Group", "Email"},
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
func NewSupporter(api *godig.API, dir string, fn *string) *E {
	var f R
	if fn != nil {
		type field struct {
			fields R
		}
		b, err := ioutil.ReadFile(*fn)
		var x field
		err = yaml.Unmarshal(b, &x)
		if err != nil {
			panic(err)
		}
		f = x.fields
	} else {
		f = R{
			"email":           "Email",
			"title":           "Title",
			"firstName":       "First_Name",
			"middleName":      "MI",
			"lastName":        "Last_Name",
			"suffix":          "Suffix",
			"status":          "Receive_Email",
			"addressLine1":    "Street",
			"addressLine2":    "Street_2",
			"city":            "City",
			"state":           "State",
			"country":         "Country",
			"postalCode":      "Zip",
			"homePhone":       "Phone",
			"cellPhone":       "Cell_Phone",
			"workPhone":       "Work_Phone",
			"timezone":        "Timezone",
			"languageCode":    "Language_Code",
			"exernalSystemId": "supporter_KEY",
			"Organization":    "Organization",
			"Department":      "Department",
			"Occupation":      "Occupation",
		}
	}
	h := []string{
		"externalSystemId",
		"email",
		"status",
		"title",
		"firstName",
		"middleName",
		"lastName",
		"suffix",
		"addressLine1",
		"addressLine2",
		"city",
		"state",
		"country",
		"postalCode",
		"homePhone",
		"cellPhone",
		"workPhone",
		"timezone",
		"languageCode",
		"Organization",
		"Department",
		"Occupation",
	}

	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
	}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         f,
		Headers:        h,
		Conditions:     c,
		CsvFilename:    "supporters.csv",
		TableName:      "supporter",
		CountTableName: "supporter",
		OffsetChan:     make(chan int32, queueSize),
		RecordChan:     make(chan R, queueSize),
		DoneChan:       make(chan bool, queueSize),
	}
	log.Printf("R is %+v\n", f)
	// Just a reminder...
	e.API.Verbose = false
	return &e
}

//NewActiveSupporter instantiates an envionrment to copy active supportes to
//CSV files.  Active supporters have a good email address and have not opted out
//or been opted out (i.e. Receive_Email > 0).
func NewActiveSupporter(api *godig.API, dir string, fn *string) *E {
	e := NewSupporter(api, dir, fn)
	c := e.Conditions
	c = append(c, "Receive_Email>0")
	e.Conditions = c
	return e
}

//NewInactiveSupporter instantiates an envionrment to copy inactive supporters to
//CSV files.  Inactive supporters have a good email address but have either opted
// out or been opted out (i.e. Receive_Email < 1).
func NewInactiveSupporter(api *godig.API, dir string, fn *string) *E {
	e := NewSupporter(api, dir, fn)
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
func NewInactiveDonors(api *godig.API, dir string, fn *string) *E {
	e := NewSupporter(api, dir, fn)

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

//headers accepts a field map and returns the keys.  Used as the header
//line in a CSV file.
func headers(f R) []string {
	var h []string
	for k := range f {
		h = append(h, k)
	}
	sort.Strings(h)
	return h
}
