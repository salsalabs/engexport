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
		"RESULT IN 0,-1"}

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
		DisableInclude: true,
	}
	return &e
}

//NewSupporter instantiates an environment for copying supporters to  CSV files.
//The default behavior is to save suupporters that have valid email addresses.
//That means that both subscribed and unsubscribed supporrters are written to CSV
//files.  TODO: Allow a user to override these selections with a YAML file.
//
//Here to Serve version notes.
//
//The Here to Serve data dictionary mapped Salsa fields (on the right) to Engage
//fields (on the left).  Note that this is not a hard-and-fast rule, as we'll see
//in a moment.
// Some of the mappings involved filtering and catenation.
// The Engage fields that need filtering and catenation have empty Salsa field
//names (""). The values for the matching Engage fields are generated in
//"save.go" using functions written for Here to Serve.
//
// The list of Engage fields can contain Salsa field names.  I had to do that
//so that there was a way to easily retrieve the Salsa values needed for
//filtering and catenation.
//
//Because of the Salsa-fields-in-the-Engage-fields column situation, the
//headers are hard-coded.  The headers are used to import into Engage and will
//be valid and/oor useful Engage field names.
//
//If this all sounds really kludgy to you, then you are correct.  It is *very*
//kludgy.
//One last thing.   This layout closely follows the Here to Serve data
//dictionary.  That means that there will be Engage custom fields mixed in
//with Engage standard fields.  We'll shuffle things around if the Conversion
//folks need them to be better organized.
func NewSupporter(api *godig.API, dir string) *E {
	f := R{
		"supporter_KEY":                    "supporter_KEY",
		"person_prefix":                    "Title",
		"person_firstname":                 "First_Name",
		"person_middlename":                "MI",
		"person_lastname":                  "Last_Name",
		"Email":                            "Email",
		"Receive_Email":                    "Receive_Email",
		"Home_Phone":                       "Phone",
		"Cell_Phone":                       "Cell_Phone",
		"Work_Phone":                       "Work_Phone",
		"Address_Line_1":                   "Street",
		"Address_Line_2":                   "Street_2",
		"City":                             "City",
		"State":                            "State",
		"Zip_Code":                         "Zip",
		"Country":                          "Country",
		"alt_email_supporter":              "Alternative_Email",
		"other_data_1_supporter":           "Other_Data_1",
		"other_data_2_supporter":           "Other_Data_2",
		"other_data_3_supporter":           "",
		"Other_Data_3":                     "Other_Data_3",
		"ncoa_codes":                       "ncoa_codes",
		"source_supporter":                 "Source",
		"source_supporter_details":         "Source_Details",
		"source_supporter_tracking_code":   "Source_Tracking_Code",
		"supporter_tracking_code":          "Tracking_Code",
		"how_would_you_like_to_help_other": "contacting_on_behalf_of___other",
		"phone_number_type_2":              "phone___secondary_type",
		"pref_method_contact":              "prefered_method_of_contact",
		"Etapestry_envelope_salutation_supporter": "etapestry___envelope_salutation",
		"Etapestry_long_salutation_supporter":     "etapestry___long_salutation",
		"Etapestry_persona_type_supporter":        "eTapestry__persona_type",
		"company":                                 "company",
		"job_title":                               "job_title",
		"spouse_prefered_phone":                   "phone___spouse_prefered_phone",
		"spouse_secondary_phone":                  "phone___spouse_secondary_phone",
		"supporter_employer_name":                 "employer_s_name",
		"supporter_employer_phone":                "employer_s_phone_number",
		"supporter_employer_hr_contact_name":      "",
		"human_resources_contact___first_name":    "human_resources_contact___first_name",
		"human_resources_contact___last_name":     "human_resources_contact___last_name",
		"human_resources_contact_name":            "human_resources_contact_name",
		"sub_how_do_you_know_this_person":         "how_do_you_know_this_person_",
		"friend_of_a_friend_name_supporter":       "",
		"friend_of_a_friend___first_name":         "friend_of_a_friend___first_name",
		"friend_of_a_friend___last_name":          "friend_of_a_friend___last_name",
		"friend_of_friend_name":                   "friend_of_friend_name",
		"family_being_helped":                     "family_being_helped",
		"primary_phone_number":                    "primary_phone",
		"supporter_org_worship_name":              "name_of_organization_or_house_of_worship",
		"supporter_org_worship_phone":             "organization_or_house_of_worship___phone_number",
		"supporter_org_worship_address_1":         "organization_or_house_of_worship___address",
		"supporter_org_worship_city":              "organization_or_house_of_worship___city",
		"supporter_org_worship_state":             "organization_or_house_of_worship___state",
		"supporter_org_worship_zip":               "organization_or_house_of_worship___postal_code",
		"account_type_supporter":                  "account_type",
		"gender_supporter":                        "gender",
		"martial_status_supporter":                "marital_status",
		"family_being_helped_other_1":             "people_helped",
		"phone_secondary":                         "phone___secondary",
		"phone_secondary_type":                    "",
		"account_name":                            "account_name_supporter",
		"care_community_comments":                 "who_would_you_like_to_help____comments",
		"phone_number_type_1":                     "primary_phone_type",
		"how_would_you_like_to_help":              "how_would_you_like_to_help",
		"supporter_contact_reason":                "contacting_on_behalf_of",
		"volunteer_skill_to_offer":                "",
		"skill___health_care_provider_type":       "skill___health_care_provider_type", //1
		"skill___computer_internet_type":          "skill___computer_internet_type",    //2
		"skill___microsoft_office_type":           "skill___microsoft_office_type",     //2
		"skill___cpa_finance_type":                "skill___cpa_finance_type",          //3
		"skill___attorney_type":                   "skill___attorney_type",             //4
		"skill___counseling_type":                 "skill___counseling_type",           //5
		"skill___other_type":                      "skill___other_type",                //6, 7, 8, 9, 10, 11
		"volunteer_skill_to_offer_other":          "",
	}

	c := []string{
		"Email IS NOT EMPTY",
		"Email LIKE %@%.%",
	}
	h := []string{
		"supporter_KEY",
		"Email",
		"person_prefix",
		"person_firstname",
		"person_middlename",
		"person_lastname",
		"Suffix",
		"Home_Phone",
		"Cell_Phone",
		"Work_Phone",
		"Phone",
		"Receive_Email",
		"Address_Line_1",
		"Address_Line_2",
		"City",
		"State",
		"Zip_Code",
		"Country",
		"Timezone",
		"Language_Code",
		"alt_email_supporter",
		"other_data_1_supporter",
		"other_data_2_supporter",
		"other_data_3_supporter",
		"source_supporter",
		"source_supporter_details",
		"source_supporter_tracking_code",
		"supporter_tracking_code",
		"how_would_you_like_to_help_other",
		"phone_number_type_1",
		"phone_number_type_2",
		"pref_method_contact",
		"Etapestry_envelope_salutation_supporter",
		"Etapestry_long_salutation_supporter",
		"Etapestry_persona_type_supporter",
		"supporter_employer_name",
		"employer_s_phone_number",
		"company",
		"job_title",
		"spouse_prefered_phone",
		"spouse_secondary_phone",
		"supporter_employer_hr_contact_name",
		"sub_how_do_you_know_this_person",
		"friend_of_a_friend_name_supporter",
		"family_being_helped",
		"primary_phone_number",
		"supporter_org_worship_name",
		"supporter_org_worship_phone",
		"supporter_org_worship_address_1",
		"supporter_org_worship_city",
		"supporter_org_worship_state",
		"supporter_org_worship_zip",
		"account_type_supporter",
		"human_resources_contact_name",
		"gender",
		"marital_status",
		"people_helped",
		"phone_secondary",
		"phone_secondary_type",
		"account_name",
		"care_community_comments",
		"primary_phone_type",
		"how_would_you_like_to_help",
		"contacting_on_behalf_of",
		"volunteer_skill_to_offer",
		"volunteer_skill_to_offer_other",
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
