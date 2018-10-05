package engexport

import (
	"encoding/csv"
	"os"

	"github.com/salsalabs/godig"
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
const Threads = 25

//R is a map between output columns and supporters table columns.
//
//Oddly enough, Go makes a distrinction thie type and a map[string]string.
//When indoubt, use this.
//
//TODO: Given this type a better name.
type R map[string]string

//E is a runtime environment for a single application.  It contains
//everything that an application needs to read stuff from Salsa and
//write CSV files. TODO: Given this type a better name.
type E struct {
	API            *godig.API
	Tag            *string
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

//P passes arguments from main() to the rest of the app.
type P struct {
	API *godig.API
	T   Schema
	Tag *string
	Dir string
}

//Processor defines the tools that an engexport processor needs.
type Processor interface {
	Drive(id int)
	Save()
	WaitFor(count int)
	Run(Threads int, start int32)
	Open(f *os.File, w *csv.Writer) (*os.File, *csv.Writer, error)
}

//Schema defines the mapping of Classic fields to CSV output.
//Note that the app expects to have a YAML table to provide
//field map and heading details.
type Schema struct {
	Supporter struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
	} `yaml:"supporter"`
	Donation struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
	} `yaml:"donation"`
	Groups struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
	} `yaml:"groups"`
}
