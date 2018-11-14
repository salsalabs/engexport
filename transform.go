package engexport

import (
	"strconv"
	"strings"

	"github.com/salsalabs/godig"
)

//Transform cleans up the value and returns it.  Separated from Save so
//that save can be the same and clients can differ with just the Transform.
func Transform(m string, d R) string {
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
	case "Date_Created":
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
