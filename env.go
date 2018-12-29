package engexport

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const queueSize = 100

//NewEnv instantiances the non-varying part of an enviroonment object.
func NewEnv(p P) *E {
	e := E{
		API:        p.API,
		Schema:     p.T,
		Tag:        p.Tag,
		OutDir:     p.Dir,
		OffsetChan: make(chan int32, queueSize),
		RecordChan: make(chan R, queueSize),
		DoneChan:   make(chan bool),
	}
	return &e

}

//NewDonation instantiates an environment for copying donations to CSV files.
//TODO: Allow a user to iverride these selections with a YAML file.
func NewDonation(p P) *E {
	c := []string{
		"RESULT IN 0,-1",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%"}
	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Donation.Fields
	e.Headers = p.T.Donation.Headers
	e.Keys = p.T.Donation.Keys
	e.CsvFilename = "donations.csv"
	e.TableName = "donation(supporter_KEY)supporter"
	e.CountTableName = "donation"
	e.PrimaryKey = "donation_KEY"
	return e
}

//NewActiveGroups instantiates an environment for copying Groups and Emails
//to CSV files.
func NewActiveGroups(p P) *E {
	c := []string{
		"groups.Group_Name IS NOT EMPTY",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%",
		//Left join munges Receive_Email, see SCT-969.
		//"supporter.Receive_Email>0",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Groups.Fields
	e.Headers = p.T.Groups.Headers
	e.Keys = p.T.Groups.Keys
	e.CsvFilename = "groups.csv"
	e.TableName = "groups(groups_KEY)supporter_groups(supporter_KEY)supporter"
	e.CountTableName = "supporter_groups"
	e.PrimaryKey = "groups_KEY"
	return e
}

//NewAllGroups instantiates an environment for copying Groups and Emails
//to CSV files for all supporters, email or not.
func NewAllGroups(p P) *E {
	c := []string{
		"groups.Group_Name IS NOT EMPTY",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Groups.Fields
	e.Headers = p.T.Groups.Headers
	e.Keys = p.T.Groups.Keys
	e.CsvFilename = "groups.csv"
	e.TableName = "groups(groups_KEY)supporter_groups(supporter_KEY)supporter"
	e.CountTableName = "supporter_groups"
	e.PrimaryKey = "groups_KEY"
	return e
}

//NewTagGroups instantiates an environment for copying tag names and Emails
//to CSV files.  The CSV files will be imported as groups in Engage.
func NewTagGroups(p P) *E {
	c := []string{
		"tag.tag IS NOT EMPTY",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%",
		//Left join munges Receive_Email, see SCT-969.
		//"supporter.Receive_Email>0",
		"tag_data.database_table_KEY=142",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Tag.Fields
	e.Headers = p.T.Tag.Headers
	e.Keys = p.T.Tag.Keys
	e.CsvFilename = "tag_groups.csv"
	e.TableName = "tag(tag_KEY)tag_data(tag_data.table_key=supporter.supporter_KEY)supporter"
	e.CountTableName = "tag_data"
	e.PrimaryKey = "tag.tag_KEY"
	e.PrimaryKeyMatchFills = "tag"
	return e
}

//NewEmailOnlyGroups instantiates an environment for copying Groups and Emails
//to CSV files where the only requirement is that a supporter has an email.
//There is no requirement for being able to deliver to the supporter.
func NewEmailOnlyGroups(p P) *E {
	e := NewActiveGroups(p)
	c := []string{
		"groups.Group_Name IS NOT EMPTY",
		"supporter.Email IS NOT EMPTY",
		"supporter.Email LIKE %@%.%",
	}
	e.Conditions = c
	return e
}

//NewSupporter instantiates an environment for copying supporters to CSV files.
//The default behavior is to save supporters that have valid email addresses.
//That means that both subscribed and unsubscribed supporrters are written to CSV
//files.  TODO: Allow a user to iverride these selections with a YAML file.
func NewSupporter(p P) *E {
	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Supporter.Fields
	e.Headers = p.T.Supporter.Headers
	e.Keys = p.T.Supporter.Keys
	e.CsvFilename = "supporters.csv"
	e.TableName = "supporter"
	e.CountTableName = "supporter"
	e.PrimaryKey = "supporter_KEY"
	return e
}

//NewActiveSupporter instantiates an envionrment to copy active supportes to
//CSV files.  Active supporters have a good email address and have not opted out
//or been opted out (i.e. Receive_Email > 0).
func NewActiveSupporter(p P) *E {
	e := NewSupporter(p)
	c := e.Conditions
	c = append(c, "Receive_Email>0")
	e.Conditions = c
	return e
}

//NewAllSupporters instantiates an environment to copy all supporter to
//CSV files.  Not a good idea for Engage, but useful for other vendors.
func NewAllSupporters(p P) *E {
	e := NewSupporter(p)
	e.Conditions = []string{
		"supporter_KEY>0",
	}
	return e
}

//NewInactiveSupporter instantiates an envionrment to copy inactive supporters to
//CSV files.  Inactive supporters have a good email address but have either opted
// out or been opted out (i.e. Receive_Email < 1).
func NewInactiveSupporter(p P) *E {
	e := NewSupporter(p)
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
func NewInactiveDonors(p P) *E {
	e := NewSupporter(p)

	//Salsa gets confused with the contents of the "&include=" in selective places.
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

//NewAllActions instantiates an environment for copying Emails and Actions
//to CSV files for all supporters, email or not.
func NewAllActions(p P) *E {
	c := []string{
		"action.action_KEY IS NOT EMPTY",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Action.Fields
	e.Headers = p.T.Action.Headers
	e.Keys = p.T.Action.Keys
	e.CsvFilename = "supporter_actions.csv"
	e.TableName = "action(action_KEY)supporter_action(supporter_KEY)supporter"
	e.CountTableName = "supporter_action"
	e.PrimaryKey = "action_KEY"
	return e
}

//NewAllEvents instantiates an environment for copying Emails and Actions
//to CSV files for all supporters, email or not.
func NewAllEvents(p P) *E {
	c := []string{
		"event.event_KEY IS NOT EMPTY",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.Event.Fields
	e.Headers = p.T.Event.Headers
	e.Keys = p.T.Event.Keys
	e.CsvFilename = "supporter_events.csv"
	e.TableName = "event(event_KEY)supporter_event(supporter_KEY)supporter"
	e.CountTableName = "supporter_event"
	e.PrimaryKey = "event_KEY"
	return e
}

//NewContactHistory instantiates an environment for copying Groups and Emails
//to CSV files for all supporters, email or not.
func NewContactHistory(p P) *E {
	c := []string{
		"contact_history.contact_history_KEY>0",
	}

	e := NewEnv(p)
	e.Conditions = c
	e.Fields = p.T.ContactHistory.Fields
	e.Headers = p.T.ContactHistory.Headers
	e.Keys = p.T.ContactHistory.Keys
	e.CsvFilename = "contact_history.csv"
	e.TableName = "campaign_manager(campaign_manager_KEY)contact_history(supporter_KEY)supporter"
	e.CountTableName = "contact_history"
	e.PrimaryKey = "contact_history_KEY"
	return e
}

//LoadSchema accepts a pointer to a YAML filename and loads the contents
//into a Table object.
func LoadSchema(fn *string) (Schema, error) {
	var t Schema
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
