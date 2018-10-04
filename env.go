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
			"supporter_KEY":                     "supporter_KEY",
			"Last_Modified":                     "Last_Modified",
			"Date_Created":                      "Date_Created",
			"Title":                             "Title",
			"First_Name":                        "First_Name",
			"MI":                                "MI",
			"Last_Name":                         "Last_Name",
			"Suffix":                            "Suffix",
			"Email":                             "Email",
			"Password":                          "Password",
			"Receive_Email":                     "Receive_Email",
			"Email_Status":                      "Email_Status",
			"Email_Preference":                  "Email_Preference",
			"Soft_Bounce_Count":                 "Soft_Bounce_Count",
			"Hard_Bounce_Count":                 "Hard_Bounce_Count",
			"Last_Bounce":                       "Last_Bounce",
			"Receive_Phone_Blasts":              "Receive_Phone_Blasts",
			"Phone":                             "Phone",
			"Cell_Phone":                        "Cell_Phone",
			"Phone_Provider":                    "Phone_Provider",
			"Work_Phone":                        "Work_Phone",
			"Pager":                             "Pager",
			"Home_Fax":                          "Home_Fax",
			"Work_Fax":                          "Work_Fax",
			"Street":                            "Street",
			"Street_2":                          "Street_2",
			"Street_3":                          "Street_3",
			"City":                              "City",
			"State":                             "State",
			"Zip":                               "Zip",
			"PRIVATE_Zip_Plus_4":                "PRIVATE_Zip_Plus_4",
			"County":                            "County",
			"District":                          "District",
			"Country":                           "Country",
			"Latitude":                          "Latitude",
			"Longitude":                         "Longitude",
			"Organization":                      "Organization",
			"Department":                        "Department",
			"Occupation":                        "Occupation",
			"Instant_Messenger_Service":         "Instant_Messenger_Service",
			"Instant_Messenger_Name":            "Instant_Messenger_Name",
			"Web_Page":                          "Web_Page",
			"Alternative_Email":                 "Alternative_Email",
			"Other_Data_1":                      "Other_Data_1",
			"Other_Data_2":                      "Other_Data_2",
			"Other_Data_3":                      "Other_Data_3",
			"Notes":                             "Notes",
			"Source":                            "Source",
			"Source_Details":                    "Source_Details",
			"Source_Tracking_Code":              "Source_Tracking_Code",
			"Tracking_Code":                     "Tracking_Code",
			"Status":                            "Status",
			"uid":                               "uid",
			"Timezone":                          "Timezone",
			"Language_Code":                     "Language_Code",
			"salesforce_id":                     "salesforce_id",
			"append_run_date":                   "append_run_date",
			"salsa_deleted":                     "salsa_deleted",
			"immigrant_parents_":                "immigrant_parents_",
			"worker_s_rights":                   "worker_s_rights",
			"african_american_":                 "african_american_",
			"number_of_children_under_18":       "number_of_children_under_18",
			"date_became_lpr":                   "date_became_lpr",
			"receit_received":                   "receit_received",
			"biometric_appmt_":                  "biometric_appmt_",
			"reunite_haitian_families":          "reunite_haitian_families",
			"donation_options":                  "donation_options",
			"immigration_reform":                "immigration_reform",
			"date_of_birth_____mm_dd_yy_":       "date_of_birth_____mm_dd_yy_",
			"marital_status":                    "marital_status",
			"gender":                            "gender",
			"stop_deportations":                 "stop_deportations",
			"drivers_licenses":                  "drivers_licenses",
			"minimum_wage":                      "minimum_wage",
			"citizen":                           "citizen",
			"organizationtype":                  "organizationtype",
			"miamideepdive":                     "miamideepdive",
			"country_of_origin":                 "country_of_origin",
			"language":                          "language",
			"childcare":                         "childcare",
			"immigrant_":                        "immigrant_",
			"requires_identifying_as":           "requires_identifying_as",
			"how_did_you_hear_about_us":         "how_did_you_hear_about_us",
			"interested_in_membership":          "interested_in_membership",
			"fc_transportation_needed":          "fc_transportation_needed",
			"mono_language":                     "mono_language",
			"monolingual":                       "monolingual",
			"voterreg":                          "voterreg",
			"miami_dive_experience_attendant_":  "miami_dive_experience_attendant_",
			"what_track_are_you_interested_in_": "what_track_are_you_interested_in_",
			"Ages_your_children":                "Ages_your_children",
			"flic_congress_field_speaker":       "flic_congress_field_speaker",
			"translation_needed__":              "translation_needed__",
			"dietary_needs":                     "dietary_needs",
			"do_you_need_interpretation_during_the_conference_": "do_you_need_interpretation_during_the_conference_",
			"fliccongresspayment":                               "fliccongresspayment",
			"form_completed":                                    "form_completed",
			"fee_waiver_used__accepted_or_denied_":              "fee_waiver_used__accepted_or_denied_",
			"status_of_application":                             "status_of_application",
			"why_":                                              "why_",
			"denied":                                            "denied",
			"interview_date":                                    "interview_date",
			"oath_ceremony":                                     "oath_ceremony",
		}
	}
	h := []string{
		"supporter_KEY",
		"Last_Modified",
		"Date_Created",
		"Title",
		"First_Name",
		"MI",
		"Last_Name",
		"Suffix",
		"Email",
		"Password",
		"Receive_Email",
		"Email_Status",
		"Email_Preference",
		"Soft_Bounce_Count",
		"Hard_Bounce_Count",
		"Last_Bounce",
		"Receive_Phone_Blasts",
		"Phone",
		"Cell_Phone",
		"Phone_Provider",
		"Work_Phone",
		"Pager",
		"Home_Fax",
		"Work_Fax",
		"Street",
		"Street_2",
		"Street_3",
		"City",
		"State",
		"Zip",
		"PRIVATE_Zip_Plus_4",
		"County",
		"District",
		"Country",
		"Latitude",
		"Longitude",
		"Organization",
		"Department",
		"Occupation",
		"Instant_Messenger_Service",
		"Instant_Messenger_Name",
		"Web_Page",
		"Alternative_Email",
		"Other_Data_1",
		"Other_Data_2",
		"Other_Data_3",
		"Notes",
		"Source",
		"Source_Details",
		"Source_Tracking_Code",
		"Tracking_Code",
		"Status",
		"uid",
		"Timezone",
		"Language_Code",
		"salesforce_id",
		"append_run_date",
		"salsa_deleted",
		"immigrant_parents_",
		"worker_s_rights",
		"african_american_",
		"number_of_children_under_18",
		"date_became_lpr",
		"receit_received",
		"biometric_appmt_",
		"receit_received",
		"biometric_appmt_",
		"reunite_haitian_families",
		"donation_options",
		"immigration_reform",
		"date_of_birth_____mm_dd_yy_",
		"marital_status",
		"gender",
		"stop_deportations",
		"drivers_licenses",
		"minimum_wage",
		"citizen",
		"organizationtype",
		"miamideepdive",
		"country_of_origin",
		"language",
		"childcare",
		"immigrant_",
		"requires_identifying_as",
		"how_did_you_hear_about_us",
		"interested_in_membership",
		"fc_transportation_needed",
		"mono_language",
		"monolingual",
		"voterreg",
		"miami_dive_experience_attendant_",
		"what_track_are_you_interested_in_",
		"Ages_your_children",
		"flic_congress_field_speaker",
		"translation_needed__",
		"dietary_needs",
		"do_you_need_interpretation_during_the_conference_",
		"fliccongresspayment",
		"form_completed",
		"fee_waiver_used__accepted_or_denied_",
		"status_of_application",
		"receit_received",
		"why_",
		"biometric_appmt_",
		"denied",
		"interview_date",
		"oath_ceremony",
	}

	c := []string{"supporter_KEY>0"}

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
