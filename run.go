package engexport

import (
	"fmt"
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
	//use the whole donations table as a guide.  Drivers will get zero
	//records at some point.  That causes a graceful shutdown.

	t := env.API.NewTable(env.CountTableName)
	cond := ""
	switch env.CountTableName {
	case "donation":
	case "supporter_groups":
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
	fmt.Printf("run: Using %d records from %v\n", m, env.CountTableName)
	for i := int64(start); i < m+499; i += 500 {
		fmt.Printf("pushing %v\n", i)
		env.OffsetChan <- int32(i)
	}

	fmt.Println("run waiting")
	time.Sleep(20 * time.Second)
	close(env.OffsetChan)
	wg.Wait()
	fmt.Println("run done")
}
