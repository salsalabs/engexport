package engexport

import (
	"fmt"
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
	incl := strings.Join(f, ",")
	cond = fmt.Sprintf("%v&include=%v", cond, incl)

	for {
		offset, ok := <-env.OffsetChan
		if !ok {
			fmt.Printf("drive_%02d done", id)
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
		if len(a) == 0 {
			break
		}
		fmt.Printf("drive_%02d: %6d returned %d\n", id, offset, len(a))
		for _, r := range a {
			env.RecordChan <- r
		}
	}
	fmt.Printf("drive_%02d: finished\n", id)
	env.DoneChan <- true
	return nil
}
