package engexport

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/dustin/go-humanize"
)

//Drive reads qualified records from Salsa and writes them to a save queue.
func (env *E) Drive(id int) {
	t := env.API.NewTable(env.TableName)
	c := env.Conditions
	checkPrimaryKey := false

	//If there are keys in the schema then we'll need to add an "IN"
	//clause to the API call to filter down to just those keys.
	if len(env.PrimaryKey) != 0 && len(env.PrimaryKeyMatchFills) != 0 && len(env.Keys) != 0 {
		checkPrimaryKey = true
		var keys []string
		for k := range env.Keys {
			keys = append(keys, k)
		}
		k := strings.Join(keys, ",")
		kc := fmt.Sprintf("%v IN %v", env.PrimaryKey, k)
		c = append(c, kc)
	}
	cond := strings.Join(c, "&condition=")

	//Salsa doesn't react well to some include queries in some calls.  Adding
	//the "&include=" can cause errors even though the URL is clearly well-formed.
	if !env.DisableInclude {
		var f []string
		for _, v := range env.Fields {
			if len(v) != 0 {
				f = append(f, v)
			}
		}
		incl := strings.Join(f, ",")
		cond = fmt.Sprintf("%v&include=%v", cond, incl)
	}
	for {
		offset, ok := <-env.OffsetChan
		if !ok {
			break
		}
		var a []map[string]string
		var err error

		if strings.Index(env.TableName, ")") != -1 {
			a, err = t.LeftJoinMap(offset, 500, cond)
		} else {
			if env.Tag != nil {
				a, err = t.ManyMapTagged(offset, 500, cond, *env.Tag)
			} else {
				a, err = t.ManyMap(offset, 500, cond)
			}
		}
		if err != nil {
			panic(err)
		}
		if math.Mod(float64(offset), 10e3) == 0 {
			x := float64(offset)
			m := humanize.FormatFloat("###,###,###", x)
			log.Printf("drive_%02d: %s\n", id, m)
		}
		if len(a) == 0 {
			break
		}
		for _, r := range a {
			// Massage the data if the primary key in the record
			//is one of the primary keys in the schema.
			//
			if checkPrimaryKey {
				pk := strings.Split(env.PrimaryKey, ".")[1]
				k, ok := r[pk]
				if ok {
					n, ok := env.Keys[k]
					if ok {
						r["tag"] = n
					}
					env.RecordChan <- r
				}
			} else {
				env.RecordChan <- r
			}
		}
	}
	env.DoneChan <- true
	log.Printf("drive_%02d: signaled done\n", id)
}
