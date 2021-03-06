//Describe shows the structure of a table.
package main

import (
	"fmt"
	"log"
	"strings"

	godig "github.com/salsalabs/godig/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func one(a *godig.API, table string) {
	t := a.NewTable(table)
	target, err := t.Describe()
	if err != nil {
		log.Fatalf("Describe %v on %v\n", err, table)
	}
	fmt.Printf("%v:\n", table)
	fmt.Printf("    fieldmap:\n")
	for _, r := range target {
		r.Label = strings.Replace(r.Label, " ", "_", -1)
		fmt.Printf("        \"%v\": \"%v\"\n", r.Label, r.Name)
	}
	fmt.Printf("    headers:\n")
	for _, r := range target {
		r.Label = strings.Replace(r.Label, " ", "_", -1)
		fmt.Printf("        - \"%v\"\n", r.Label)
	}
	fmt.Printf("\n")

}
func main() {
	cpath := kingpin.Flag("login", "YAML file containing login for Salsa Classic API").PlaceHolder("FILENAME").Required().String()
	kingpin.Parse()

	a, err := godig.YAMLAuth(*cpath)
	if err != nil {
		log.Fatalf("Authentication error: %+v\n", err)
	}
	one(a, "supporter")
	one(a, "donation")
	one(a, "groups")
	one(a, "supporter_groups")
	one(a, "action")
	one(a, "supporter_action")
	one(a, "event")
	one(a, "supporter_action")
	one(a, "contact_history")
	one(a, "supporter_email_statistics")
}
