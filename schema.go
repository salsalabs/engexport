package engexport

import (
	"encoding/csv"
	"github.com/salsalabs/godig"
	"os"
)

//ParseFmt is used to parse a Salsa database mesasge.  NOote that the only
//way to do that is to remove the hour offset before parsing.  See `date()`.
const ParseFmt = "Mon Jan 2 2006 15:04:05 (MST)"

//LayoutFmt is used to format a time so that Engage will recognize it.
const LayoutFmt = "2006-Jan-02T15:04:05"

//RecordsPerFile defines the maximum number of data recors for each CSV
//file.  The actual count is RecordsPerFile + 1 for the header.
const RecordsPerFile = 50000

//Threads is the number of Drive threads.  That's the one that is the
//"slowest".  WaitFor is responsible for shutting things down wihen
//this number of Drives had pushed a "done" message.
const Threads = 5

//R is a record from the database.  It's a basic map of strings pointing to
//strings.  Oddly enough, Go makes a distrinction between the two.  When in
//doubt, use this. TODO: Given this type a better name.
type R map[string]string

//E is a runtime environment for a single application.  It contains
//everything that an application needs to read stuff from Salsa and
//write CSV files.  Note that E is used to iinstantiate the Procwessor
//interface. TODO: Given this type a better name.
type E struct {
	API            *godig.API
	OutDir         string
	Fields         R
	Headers        []string
	DisableInclude bool
	Conditions     []string
	CsvFilename    string
	TableName      string
	CountTableName string
	OffsetChan     chan int32
	RecordChan     chan R
	DoneChan       chan bool
}

//Processor defines the tools that an engexport processor needs.
type Processor interface {
	Drive(id int) error
	Save() error
	WaitFor(count int) error
	Run(Threads int, start int32)
	Open(id int, f *os.File, w *csv.Writer) (*os.File, *csv.Writer, error)
}
