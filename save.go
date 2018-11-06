package engexport

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/salsalabs/godig"
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
			log.Println("save done")
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
<<<<<<< HEAD
			//KLUDGE:  Salsa wants to see supporter.supporter_KEY/supporter.Email
			// in the conditions and included fields.  However, the data is stored
			// simply as "supporter_KEY" or "Email"...
			i := strings.Index(m, ".")
			if i != -1 {
				m = strings.Split(m, ".")[1]
			}

			s, ok := d[m]
			if !ok {
				s = ""
			}
			//Transform fields as needed.  This includes making pretty dates,
			//setting the Engage transaction type and putting Engage text into
			//Receive_Email.
			switch m {
			case "friend_of_a_friend___first_name":
				//merge this and friend_of_a_friend___last_name
				//into friend_of_a_friend_name
			case "State":
				s = strings.ToUpper(s)
			case "Transaction_Date":
				s = godig.EngageDate(s)
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
			}

			// Convert tabs to spaces. Remove leading/trailing spaces.
			// Append the cleaned-up value to the output.
			s = strings.Replace(s, "\t", " ", -1)
			s = strings.TrimSpace(s)
=======
			s := transform(m, d)
>>>>>>> add_tags
			a = append(a, s)
		}
		err := w.Write(a)
		if err != nil {
			panic(err)
		}
		count++
		w.Flush()
	}
	if w != nil {
		w.Flush()
	}
	if f != nil {
		f.Close()
	}
}
<<<<<<< HEAD
=======

//transform cleans up the value and returns it.
func transform(m string, d R) string {
	//KLUDGE:  Salsa wants to see supporter.supporter_KEY/supporter.Email
	// in the conditions and included fields.  However, the data is stored
	// simply as "supporter_KEY" or "Email"...
	i := strings.Index(m, ".")
	if i != -1 {
		m = strings.Split(m, ".")[1]
	}
	s, ok := d[m]
	if !ok {
		s = ""
	}
	//Transform fields as needed.  This includes making pretty dates,
	//setting the Engage transaction type and putting Engage text into
	//Receive_Email.
	switch m {
	case "State":
		s = strings.ToUpper(s)
	case "Transaction_Date":
		s = godig.EngageDate(s)
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
	}
	// Convert tabs to spaces. Remove leading/trailing spaces.
	// Remove any quotation marks.
	// Append the cleaned-up value to the output.
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "\"", "", -1)
	s = strings.TrimSpace(s)

	return s
}
>>>>>>> add_tags
