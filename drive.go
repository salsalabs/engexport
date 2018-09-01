package engexport

import (
	"fmt"
	"math"
	"strings"
)

//Drive reads qualified records from Salsa and writes them to a save queue.
func (env *E) Drive(id int) error {
	t := env.API.NewTable(env.TableName)
	c := env.Conditions
	cond := strings.Join(c, "&condition=")
	var f []string
	for _, v := range env.Fields {
		f = append(f, v)
	}
	//Salsa doesn't react will to some include queries in some calls.  Adding
	//the "&include=" can cause errors even though the URL is clearly well-formed.
	if !env.DisableInclude {
		incl := strings.Join(f, ",")
		cond = fmt.Sprintf("%v&include=%v", cond, incl)
	}

	fmt.Printf("drive_%02d: begin\n", id)
	for {
		offset, ok := <-env.OffsetChan
		if !ok {
			fmt.Printf("drive_%02d: end of queue\n", id)
			break
		}
		var a []map[string]string
		var err error
		if strings.Index(env.TableName, ")") != -1 {
			a, err = t.LeftJoinMap(offset, 500, cond)
		} else {
			a, err = t.ManyMap(offset, 500, cond)
		}
		if err != nil {
			return err
		}
		if math.Mod(float64(offset), 10e3) == 0 {
			fmt.Printf("drive_%02d: %6d\n", id, offset)
		}
		if len(a) == 0 {
			fmt.Printf("drive_%02d: end of data\n", id)
			break
		}
		for _, r := range a {
			env.RecordChan <- r
		}
	}
	env.DoneChan <- true
	return nil
}
