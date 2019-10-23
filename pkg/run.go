package engexport

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Run starts all of the parts of a processor as goroutines.  It then waits
//for the goroutines to complete.
func (env *E) Run(Threads int, start int32) {
	var wg sync.WaitGroup
	go (func(wg *sync.WaitGroup, env *E, c int) {
		wg.Add(1)
		defer wg.Done()
		env.WaitFor(c)
	})(&wg, env, Threads)

	go (func(wg *sync.WaitGroup, env *E) {
		wg.Add(1)
		defer wg.Done()
		env.Save()
	})(&wg, env)

	for id := 1; id <= Threads; id++ {
		go (func(wg *sync.WaitGroup, env *E, id int) {
			wg.Add(1)
			defer wg.Done()
			env.Drive(id)
		})(&wg, env, id)
	}

	//KLUDGE: Salsa's API does not have a way to count for a LeftJoin.  We'll
	//use the whole "count" table as a guide.  Drivers will get zero
	//records at some point.  That causes a graceful shutdown.

	t := env.API.NewTable(env.CountTableName)
	cond := ""
	switch env.CountTableName {
	case "donation":
	case "supporter_groups":
	case "supporter_action":
	case "supporter_event":
	case "tag_data":
		break
	default:
		cond = strings.Join(env.Conditions, "&condition=")
	}
	x, err := t.Count(cond)
	if err != nil {
		panic(err)
	}
	m, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		panic(err)
	}
	m = m - int64(start)
	r := float64(m)
	log.Printf("run: Using %10.0f records from %v\n", r, env.CountTableName)
	go (func(wg *sync.WaitGroup, env *E, start int32) {
		wg.Add(1)
		defer wg.Done()
		for i := int64(start); i < m+499; i += 500 {
			env.OffsetChan <- int32(i)
		}
		close(env.OffsetChan)
	})(&wg, env, start)

	time.Sleep(5 * time.Second)
	wg.Wait()
	log.Println("run done")
}
