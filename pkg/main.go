package engexport

import (
	"encoding/csv"
	"os"

	godig "github.com/salsalabs/godig/pkg"
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
const Threads = 10

//R is a map between output columns and supporters table columns.
//
//Oddly enough, Go makes a distinction thie type and a map[string]string.
//When in doubt, use this.
//
//TODO: Given this type a better name.
type R map[string]string

//E is a runtime environment for a single application.  It contains
//everything that an application needs to read stuff from Salsa and
//write CSV files. TODO: Given this type a better name.
type E struct {
	API                  *godig.API
	Schema               Schema
	Tag                  *string
	OutDir               string
	Fields               R
	Headers              []string
	Keys                 R
	DisableInclude       bool
	Conditions           []string
	CsvFilename          string
	TableName            string
	CountTableName       string
	PrimaryKey           string
	PrimaryKeyMatchFills string
	OffsetChan           chan int32
	RecordChan           chan R
	DoneChan             chan bool
}

//P passes arguments from main() to the rest of the app.
type P struct {
	API            *godig.API
	T              Schema
	Tag            *string
	Dir            string
	DisableInclude bool
}

//RunConfig is read from a "run.yaml" file.  It can also be used
//to store commandline arguments.
type RunConfig struct {
	Host           string   `yaml:"host"`
	Email          string   `yaml:"email"`
	Password       string   `yaml:"password"`
	Schema         string   `yaml:"schema"`
	OutDir         string   `yaml:"dir"`
	Tag            *string  `yaml:"tag"`
	Start          int32    `yaml:"start"`
	APIVerbose     bool     `yaml:"apiVerbose"`
	DisableInclude bool     `yaml:"disableInclude"`
	DumpSchema     bool     `yaml:"dumpSchema"`
	Args           []string `yaml:"args"`
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
//Note that the app expects to have a YAML file that provides
//field map and heading details.  Note, too, that Keys are
//optional.  If not provided, then YAML returns an empty map.
type Schema struct {
	Supporter struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"supporter"`
	Donation struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"donation"`
	Groups struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"groups"`
	Tag struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"tag"`
	Action struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"action"`
	Event struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"event"`
	ContactHistory struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"contact_history"`
	EmailStatistics struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"supporter_email_statistics"`
	BlastStatistics struct {
		Fields  R        `yaml:"fieldmap"`
		Headers []string `yaml:"headers"`
		Keys    R        `yaml:"keymap"`
	} `yaml:"blast_statistics"`
}
