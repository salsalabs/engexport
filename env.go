package engexport

import (
	"sort"

	"github.com/salsalabs/godig"
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
func NewSupporter(api *godig.API, dir string) *E {
	f := R{
		"supporter_KEY":                  "supporter_KEY",
		"Email":                          "Email",
		"person_prefix":                  "Title",
		"person_firstname":               "First_Name",
		"person_middlename":              "MI",
		"person_lastname":                "Last_Name",
		"Suffix":                         "Suffix",
		"Home_Phone":                     "Phone",
		"Cell_Phone":                     "Cell_Phone",
		"Work_Phone":                     "Work_Phone",
		"Phone":                          "Phone",
		"Receive_Email":                  "Receive_Email",
		"Address_Line_1":                 "Street",
		"Address_Line_2":                 "Street_2",
		"City":                           "City",
		"State":                          "State",
		"Zip_Code":                       "Zip",
		"Country":                        "Country",
		"Timezone":                       "Timezone",
		"Language_Code":                  "Language_Code",
		"alt_email_supporter":            "Alternative_Email",
		"other_data_1_supporter":         "Other_Data_1",
		"other_data_2_supporter":         "Other_Data_2",
		"other_data_3_supporter":         "Other_Data_3",
		"source_supporter":               "Source",
		"source_supporter_details":       "Source_Details",
		"source_supporter_tracking_code": "Source_Tracking_Code",
		"supporter_tracking_code":        "Tracking_Code",

		"how_would_you_like_to_help_other": "contacting_on_behalf_of___other",

		"phone_number_type_2":                     "phone___secondary_type",
		"pref_method_contact":                     "prefered_method_of_contact",
		"Etapestry_envelope_salutation_supporter": "etapestry___envelope_salutation",
		"Etapestry_long_salutation_supporter":     "etapestry___long_salutation",
		"Etapestry_persona_type_supporter":        "etapestry__persona_type",
		"supporter_employer_name":                 "employer_s_name",

		"employer_s_phone_number": "supporter_employer_hr_contact_phone",
		//Concatenate
		"supporter_employer_hr_contact_name": "human_resources_contact___first_name",
		"supporter_employer_hr_contact_name": "human_resources_contact___last_name",

		"sub_how_do_you_know_this_person": "how_do_you_know_this_person_",
		//Concatenate
		"friend_of_a_friend_name_supporter": "friend_of_a_friend___first_name",
		"friend_of_a_friend_name_supporter": "friend_of_a_friend___last_name",
		"friend_of_a_friend_name_supporter": "friend_of_friend_name",

		"family_being_helped":             "family_being_helped",
		"primary_phone_number":            "primary_phone",
		"supporter_org_worship_name":      "name_of_organization_or_house_of_worship",
		"supporter_org_worship_phone":     "organization_or_house_of_worship___phone_number",
		"supporter_org_worship_address_1": "organization_or_house_of_worship___address",
		"supporter_org_worship_city":      "organization_or_house_of_worship___city",
		"supporter_org_worship_state":     "organization_or_house_of_worship___state",
		"supporter_org_worship_zip":       "organization_or_house_of_worship___postal_code",
		"account_type_supporter":          "account_type",
		"other_data_3_supporter":          "ncoa_codes",

		// (note: exclude values of 'n/a' and 'test' -- See Rows 74 & 75)",
		"human_resources_contact_name": "supporter_employer_hr_contact_name",

		"gender":         "gender_supporter",
		"marital_status": "martial_status_supporter",
		"people_helped":  "family_being_helped_other_1",
		//Note logic in fill
		"phone___secondary":      "phone_secondary + phone_secondary_type = 'Work'",
		"phone___secondary_type": "Work",

		"n/a":          "solicitor",
		"account_name": "account_name_supporter",
		//"n/a":                                      "business_or_organization___name",
		//"n/a":                                      "business_or_organization___address",
		//"n/a":                                      "business_or_organization___city",
		//"n/a":                                      "business_or_organization___state",
		//"n/a":                                      "business_or_organization___postal_code",
		//"n/a":                                      "business_or_organization___phone_number",
		//"n/a":                                      "business_or_organization___website",
		//"n/a":                                      "business_or_organization___fax_number",
		"who_would_you_like_to_help____comments": "care_community_comments",
		"primary_phone_type":                     "phone_number_type_1",
		"how_would_you_like_to_help":             "how_would_you_like_to_help",
		"contacting_on_behalf_of":                "supporter_contact_reason",
		//"Supporter Custom - skill_to_offer  - Radio selection to indicate a specific skill they can offer to the organization:"

		//" 1 = healthcare provider
		//" 2 = computer, technology, social media
		//" 3 = accounting, financial services
		//" 4 = legal, attorney
		//" 5 = professional counseling
		//" 6 = skilled in complex health insurance issues
		//" 7 = licensed care provider
		//" 8 = provided licensed child care
		//" 9 = I have cared for someone with a life-threatening illness
		//" 10 = I have had a life-threatening illness
		// 11 = other
		"skill_to_offer (see row 109)":                        "skill___health_care_provider_type", //1
		"skill_to_offer (see row 109)":                        "skill___attorney_type",             //4
		"skill_to_offer (see row 109)":                        "skill___cpa_finance_type",          //3
		"skill_to_offer (see row 109)":                        "skill___computer_internet_type",    //2
		"skill_to_offer (see row 109 and merge with row 112)": "skill___microsoft_office_type",     //2
		"skill_to_offer (see row 109)":                        "skill___counseling_type",           //5
		"skill_to_offer (see row 109)":                        "skill___other_type",                //6, 7, 8, 9,, 10, 11
	}

	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
	}

	e := E{
		API:            api,
		OutDir:         dir,
		Fields:         f,
		Headers:        headers(f),
		Conditions:     c,
		CsvFilename:    "supporters.csv",
		TableName:      "supporter",
		CountTableName: "supporter",
		OffsetChan:     make(chan int32, queueSize),
		RecordChan:     make(chan R, queueSize),
		DoneChan:       make(chan bool, queueSize),
	}
	// Just a reminder...
	e.API.Verbose = false
	return &e
}

//NewActiveSupporter instantiates an envionrment to copy active supportes to
//CSV files.  Active supporters have a good email address and have not opted out
//or been opted out (i.e. Receive_Email > 0).
func NewActiveSupporter(api *godig.API, dir string) *E {
	e := NewSupporter(api, dir)
	c := e.Conditions
	c = append(c, "Receive_Email>0")
	e.Conditions = c
	return e
}

//NewInactiveSupporter instantiates an envionrment to copy inactive supporters to
//CSV files.  Inactive supporters have a good email address but have either opted
// out or been opted out (i.e. Receive_Email < 1).
func NewInactiveSupporter(api *godig.API, dir string) *E {
	e := NewSupporter(api, dir)
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
func NewInactiveDonors(api *godig.API, dir string) *E {
	e := NewSupporter(api, dir)

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
