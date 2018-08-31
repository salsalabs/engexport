package engexport

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
)

//Open creates a new CSV output file.  If the current one is still open, then it's closed.
func (env *E) Open(id int, f *os.File, w *csv.Writer) (*os.File, *csv.Writer, error) {
	if w != nil {
		w.Flush()
		w = nil
	}
	if f != nil {
		f.Close()
		f = nil
	}

	e := path.Ext(env.CsvFilename)
	b := strings.Replace(env.CsvFilename, e, "", -1)
	fn := fmt.Sprintf("%s_%03d%s", b, id, e)
	fn = path.Join(env.OutDir, fn)
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
	fmt.Printf("open: %v\n", fn)
	return f, w, err
}
