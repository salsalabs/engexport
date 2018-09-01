package engexport

import (
	"fmt"
)

//WaitFor is responsible for getting a number of "done" notifications, then closing
//the save queue.
func (env *E) WaitFor(c int) {
	for {
		_, ok := <-env.DoneChan
		if !ok {
			break
		}
		c--
		if c == 0 {
			fmt.Println("waitFor done")
			break
		}
	}
	close(env.RecordChan)
}
