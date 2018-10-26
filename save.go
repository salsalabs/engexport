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
		skillsOther := ""
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

			case "supporter_org_worship_state":
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

			case "supporter_employer_hr_contact_name":
				s = humanResourcesContact(d)

			case "other_data_3_supporter":
				s = otherData3Supporter(d)

			case "phone_secondary_type":
				s = phoneSecondaryType(d)

			case "skill_to_offer":
				s, skillsOther, err = skillToOffer(env, d)
				if err != nil {
					fmt.Printf("SkillToOffer %v, %v %v\n", d["supporter_KEY"], d["Email"], err)
					s = ""
				}

			case "skill_to_offer_other":
				s = skillsOther
			}
			a = append(a, s)
		}
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
			if len(v) > 0 {
				s := strings.ToLower(v)
				if s != "na" && s != "n/a" && s != "test" {
					a = append(a, v)
				}
			}
		}
	}
	s := strings.Join(a, " ")
	s = strings.TrimSpace(s)
	return s
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

//skillTags scans for skill tag records for the provided supporter.  It returns
//a somewhat ordered list of skills that a supporter selected.  The function uses
//an internal list of tag_KEYs and their descriptions.  The descriptions are derived
//from the data dictionary.
//
//The rules are...
//
//If "family_being_helped" is not empty, then
//  "skills_to_offer" is the list of tagged skills
//  "skills_to_offer_other" is the list of "skill___.+_type" values
//Else
//  "skills_to_offer_other" is any unmapped skill
func skillToOffer(env *E, d R) (string, string, error) {
	var skills []string
	var other []string

	tagKeys := map[string]string{
		"260050": "healthcare provider",
		"260051": "computer, technology, social media",
		"260141": "accounting, financial services",
		"260142": "microsoft office proficient",
		"260145": "legal, attorney",
		"260146": "professional counseling",
		"260147": "skilled in complex health insurance issues",
		//The rest are not in Engage. They are are stored in "other".
		"260144": "provided licensed child care",
		"260047": "I have cared for someone with a life-threatening illness",
		"260048": "I have had a life-threatening illness",
		"260049": "other",
	}

	isSkill := map[string]bool{
		"260050": true,
		"260051": true,
		"260141": true,
		"260142": true,
		"260145": true,
		"260146": true,
		"260147": true,
		//The rest are not in Engage, are stored in "other".
		"260144": false,
		"260047": false,
		"260048": false,
		"260049": false,
	}

	keyOrder := []string{
		"260050",
		"260051",
		"260141",
		"260142",
		"260145",
		"260146",
		"260147",
		"260144",
		"260047",
		"260048",
		"260049",
	}

	fields := []string{
		"skill___health_care_provider_type",
		"skill___computer_internet_type",
		"skill___microsoft_office_type",
		"skill___cpa_finance_type",
		"skill___attorney_type",
		"skill___counseling_type",
		"skill___other_type",
	}

	t := env.API.NewTable("tag_data")
	inString := strings.Join(keyOrder, ",")

	conditions := []string{
		"database_table_KEY=142",
		fmt.Sprintf("table_KEY=%v", d["supporter_KEY"]),
		fmt.Sprintf("tag_KEY IN %v", inString),
	}
	crit := strings.Join(conditions, "&condition=")
	crit = crit + "&include=tag_KEY"

	//Only reading once because there are so few keys and
	//even fewer potential matches.
	a, err := t.ManyMap(0, 500, crit)
	if err != nil {
		return "", "", err
	}
	//Do for all matching tag_data records
	for _, r := range a {
		//retrieve the tag_KEY
		k, ok := r["tag_KEY"]
		if ok {
			//retrieve the official skill description
			s, ok := tagKeys[k]
			if ok {
				//see if this is an Engage skill
				t, _ := isSkill[k]
				if t {
					skills = append(skills, s)
				} else {
					other = append(other, s)
				}
			}
		} else {
			//No tag KEY.  Whine...
			fmt.Printf("%v %v, no tag_KEY in %d records of tag data.\n%v\n\n", d["supporter_KEY"], d["Email"], len(a), a)
		}
	}
	// append all "skill___*" fields to the other array.
	for _, s := range fields {
		v, ok := d[s]
		if ok {
			v = strings.TrimSpace(v)
			if len(v) > 0 {
				other = append(other, v)
			}
		}
	}
	s := strings.Join(skills, ", ")
	x := strings.Join(other, ", ")
	return s, x, nil
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
