package engexport

import (
	"regexp"
	"strconv"
	"strings"

	godig "github.com/salsalabs/godig/pkg"
)

//Transform cleans up the value and returns it.  Separated from Save so
//that save can be the same and clients can differ with just the Transform.
func Transform(m string, d R) string {
	cleanLines := regexp.MustCompile("[\\n\\r\\t]+")
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
	case "Contact_Date":
		s = godig.EngageDate(s)
	case "Contact_Due_Date":
		s = godig.EngageDate(s)
	case "Date_Created":
		s = godig.EngageDate(s)
	case "Email":
		s = strings.TrimSpace(strings.ToLower(s))
	case "End":
		s = godig.EngageDate(s)
	case "First_Email_Time":
		s = godig.EngageDate(s)
	case "last_click":
		s = godig.EngageDate(s)
	case "Last_Email_Time":
		s = godig.EngageDate(s)
	case "Last_Modified":
		s = godig.EngageDate(s)
	case "last_open":
		s = godig.EngageDate(s)
	case "Receive_Email":
		t := "Unsubscribed"
		x, err := strconv.ParseInt(s, 10, 32)
		if err == nil && x > 0 {
			t = "Subscribed"
		}
		s = t
	case "Start":
		s = godig.EngageDate(s)
	case "State":
		s = strings.ToUpper(s)
	case "Transaction_Date":
		s = godig.EngageDate(s)
	case "Transaction_Type":
		if s != "Recurring" {
			s = "OneTime"
		}
	}
	// Convert tabs to spaces. Remove leading/trailing spaces.
	// Remove any quotation marks.
	// Append the cleaned-up value to the output.
	s = strings.Replace(s, "\"", "", -1)
	s = cleanLines.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}
