package engexport

import (
	"encoding/csv"
	"log"
	"os"
)

//Save waits for records to arrive on a queue and saves them to a CSV file.  CSV files
//are created as needed and are replaces when they get full.
func (env *E) Save() {
	count := RecordsPerFile
	var f *os.File
	var w *csv.Writer
	var err error

	for {
		d, ok := <-env.RecordChan
		if !ok {
			log.Println("save done")
			break
		}
		if count >= RecordsPerFile {
			count = 0
			f, w, err = env.Open(f, w)
			if err != nil {
				panic(err)
			}
		}
		var a []string
		for _, k := range env.Headers {
			m := env.Fields[k]
			s := Transform(m, d)
			a = append(a, s)
		}
		err := w.Write(a)
		if err != nil {
			panic(err)
		}
		count++
		w.Flush()
	}
	if w != nil {
		w.Flush()
	}
	if f != nil {
		f.Close()
	}
}
