package engexport

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

//Run starts all of the parts of a processor as goroutines.  It then waits
//for the goroutines to complete.
func (env *E) Run(Threads int, start int32) error {

	var wg sync.WaitGroup
	go (func(wg *sync.WaitGroup, env *E, c int) {
		wg.Add(1)
		env.WaitFor(c)
		wg.Done()
	})(&wg, env, Threads)

	go (func(wg *sync.WaitGroup, env *E) {
		wg.Add(1)
		err := env.Save()
		wg.Done()
		if err != nil {
			panic(err)
		}
	})(&wg, env)

	for id := 1; id <= Threads; id++ {
		go (func(wg *sync.WaitGroup, env *E, id int) {
			wg.Add(1)
			err := env.Drive(id)
			wg.Done()
			if err != nil {
				panic(err)
			}
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
	fmt.Printf("run: %d records for %v\n", m, env.TableName)
	m = m - int64(start)
	for i := int64(start); i < m; i += 500 {
		env.OffsetChan <- int32(i)
	}
	close(env.OffsetChan)
	return nil

	time.Sleep(5 * time.Second)
	fmt.Println("main waiting")
	wg.Wait()
	fmt.Println("main done")
	return nil
}
