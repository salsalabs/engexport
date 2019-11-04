package engexport

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

//Open creates a new CSV output file.  If the current one is still open, then it's closed.
func (env *E) Open(f *os.File, w *csv.Writer) (*os.File, *csv.Writer, error) {
	if w != nil {
		w.Flush()
		w = nil
	}
	if f != nil {
		f.Close()
		f = nil
	}

	searching := true
	id := 1
	var fn string
	for searching {
		e := path.Ext(env.CsvFilename)
		b := strings.Replace(env.CsvFilename, e, "", -1)
		fn = fmt.Sprintf("%s_%03d%s", b, id, e)
		fn = path.Join(env.OutDir, fn)
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			searching = false
		} else {
			id++
		}
	}
	err := os.MkdirAll(path.Dir(fn), os.ModePerm)
	if err != nil {
		return f, w, err
	}
	f, err = os.Create(fn)
	if err != nil {
		return f, w, err
	}
	w = csv.NewWriter(f)
	err = w.Write(env.Headers)
	log.Printf("open: %v\n", fn)
	return f, w, err
}
