package engexport

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
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
		datePattern := regexp.MustCompile("(date|Date|last_donated)")
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
			if !ok {
				s = ""
			}
			if datePattern.MatchString(k) {
				s = date(s)
			} else {
				switch k {
				case "Transaction_Type":
					if s != "Recurring" {
						s = "OneTime"
					}
				case "Receive_Email":
					t := "Unsubscribed"
					x, err := strconv.ParseInt(s, 10, 32)
					if err == nil && x > 1 {
						t = "Subscribed"
					}
					s = t
				}
			}
			a = append(a, s)
		}
		//If the custom fields are all empty then there's no reason to write
		//this record.
		doWrite := false
		for i := 2; i < len(a); i++ {
			if len(a[i]) > 0 {
				doWrite = true
				break
			}
		}
		if doWrite {
			err := w.Write(a)
			count++
			if err != nil {
				panic(err)
			}
			w.Flush()
		}
	}
	if w != nil {
		w.Flush()
	}
	if f != nil {
		f.Close()
	}
}

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
