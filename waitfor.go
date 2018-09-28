package engexport

import "log"

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
			log.Println("waitFor done")
			break
		}
	}
	close(env.RecordChan)
}
