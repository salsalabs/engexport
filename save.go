package engexport

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//Save waits for records to arrive on a queue and saves them to a CSV file.  CSV files
//are created as needed and are replaces when they get full.
func (env *E) Save() {
	count := RecordsPerFile
	var f *os.File
	var w *csv.Writer
	var err error

	for {
		d, ok := <-env.RecordChan
		if !ok {
			fmt.Println("save done")
			break
		}
		if count >= RecordsPerFile {
			count = 0
			f, w, err = env.Open(f, w)
			if err != nil {
				panic(err)
			}
		}

		var a []string
		for _, k := range env.Headers {
			m := env.Fields[k]
			//KLUDGE:  Salsa wants to see supporter.supporter_KEY/supporter.Email
			// in the conditions and included fields.  However, the data is stored
			// simply as "supporter_KEY" or "Email"...
			i := strings.Index(m, ".")
			if i != -1 {
				m = strings.Split(m, ".")[1]
			}

			s, ok := d[m]
			if ok {
				s = strings.TrimSpace(s)
			} else {
				s = ""
			}
			//Transform fields as needed.  This includes making pretty dates,
			//setting the Engage transaction type and putting Engage text into
			//Receive_Email.
			switch k {

			case "State":
				s = strings.ToUpper(s)

			case "Transaction_Date":
				s = date(s)

			case "Transaction_Type":
				if s != "Recurring" {
					s = "OneTime"
				}

			case "Receive_Email":
				t := "Unsubscribed"
				x, err := strconv.ParseInt(s, 10, 32)
				if err == nil && x > 0 {
					t = "Subscribed"
				}
				s = t

			case "friend_of_a_friend_name_supporter":
				s = friendOfAFriend(d)

			case "human_resources_contact":
				s = humanResourcesContact(d)

			case "other_data_3_supporter":
				s = otherData3Supporter(d)

			case "phone_secondary_type":
				s = phoneSecondaryType(d)

			case "skill_to_offer":
				s = skillToOffer(d)

			}
			a = append(a, s)
		}
		//fmt.Printf("Save: a=%v\n", a)
		err := w.Write(a)
		count++
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
	if w != nil {
		w.Flush()
	}
	if f != nil {
		f.Close()
	}
}

//friend_of_friend does special formatting to transform Classic custom fields
//into a single Engage field.
func friendOfAFriend(d R) string {
	keys := []string{
		"friend_of_a_friend___first_name",
		"friend_of_a_friend___last_name",
		"friend_of_friend_name",
	}
	return catenateValues(d, keys)
}

//human_resources_contact does special formatting to transform Classic custom
//fields into a single Engage field.
func humanResourcesContact(d R) string {
	keys := []string{
		"human_resources_contact___first_name",
		"human_resources_contact___last_name",
		"human_resources_contact_name",
	}
	return catenateValues(d, keys)
}

//catenateValues accepts a record and a list of keys.  The values for the keys
//are appended and returned.
func catenateValues(d R, keys []string) string {
	var a []string
	for _, k := range keys {
		v, ok := d[k]
		if ok {
			v = strings.TrimSpace(v)
			s := strings.ToLower(v)
			if s == "n/a" || s == "test" {
				v = ""
			}
			if len(v) > 0 {
				a = append(a, v)
			}
		}
	}
	return strings.Join(a, " ")
}

//other_data_3_supporter stores data from two diverse spots into the
//Other Data 3 field inEngage.
func otherData3Supporter(d R) string {
	keys := []string{
		"Other_Data_3",
		"ncoa_codes",
	}
	return catenateValues(d, keys)
}

//phone_secondary populates the secondary phone type with "Work" as needed.
func phoneSecondaryType(d R) string {
	v, ok := d["phone___secondary"]
	s := ""
	if ok {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			s = "Work"
		}
	}
	return s
}

//skill_to_offer accepts a record and a list of keys.  Each key is interpreted
//as a numeric value.  The numeric value is appended to the returned value.
//SIDE-EFFECT: Field "skill_to_offer_other" has any non-empty values appended.
func skillToOffer(d R) string {
	keys := map[string]string{
		"skill___health_care_provider_type": "1",
		"skill___computer_internet_type":    "2",
		"skill___microsoft_office_type":     "2",
		"skill___cpa_finance_type":          "3",
		"skill___attorney_type":             "4",
		"skill___counseling_type":           "5",
	}
	var a []string
	for k, v := range keys {
		x, ok := d[k]
		if ok {
			x = strings.TrimSpace(x)
			if len(x) > 0 {
				a = append(a, v)
				s := d["skill_to_offer_other"]
				if len(s) > 0 {
					s = s + ", "
				}
				d["skill_to_offer_other"] = s + x
			}
		}
	}
	return strings.Join(a, ", ")
}

//date formates a Classic date from the database (ick) to an Engage date.
func date(s string) string {
	// Date_Created, Transaction_Date, etc.  Convert dates from MySQL to Engage.
	p := strings.Split(s, " ")
	if len(p) >= 7 {
		//Pull out the timezone.
		p = append(p[0:5], p[6])
		x := strings.Join(p, " ")
		t, err := time.Parse(ParseFmt, x)
		if err != nil {
			fmt.Printf("Warning: parsing %v returned %v\n", s, err)
		} else {
			s = t.Format(LayoutFmt)
		}
	}
	return s
}
