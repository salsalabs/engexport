## Classic Exporter

Go application to export the usual data necessary to convert a client form Classic to Engage.
* Active supporter data
* Inactive supporters that are donors
* Selected group names and supporter emails.
* (Future) Selected tag and supporter emails.
* (Future) Actions and supporter emails.
* (Future) Events and supporter emails.
* 

All output files are CSV.

## Prerequisites
1. A [recent version of Go](https://golang.org/doc/install) installed.  You should use the googles
if you are installing on Windows.  Really. Trust me on this.

1. The correct directory structure.  Here's a sample.
```
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```
## Installation
```
go get github.com/salsalabs/exporter
go install
```

## Setup
Put your Salsa Classic login credentials into a YAML file.
```yaml
host: hostname 
email: you@yours.com
password: super-secret-password
```
Where 
* `hostname` is the [API Host](https://help.salsalabs.com/hc/en-us/articles/115000341773-Salsa-Application-Program-Interface-API-#api_host)
* `email` is the email address that you use to login.
* `password` is the password that you use to login.

*Remember to remove this file after you're done!*

## Execution
* Active supporters
```
go run cmd/active/main.go --login YOUR_YAML_FILE supporters active
```
* Donations by all supporters with valid email addresses.
```
go run cmd/active/main.go --login YOUR_YAML_FILE donations
```
* Group names and emails for all active supporters.
```
go run cmd/active/main.go --login YOUR_YAML_FILE groups
```
* Inactive supporters.
```
go run cmd/inactive_and_donors/main.go --login YOUR_YAML_FILE supporters inactive
```

## Output
Output goes to a directory of our choosing.  The default is `./data`.  The output directory contains a single directory for each of the exports that the app runs.  Here's an example of the directory structure.

```
./data -- directory of data
  +- supporters -- directory of active supporters
      +- supporters_001.csv
      +- supporters_002.csv
      +- . . .
  +- inactive-donors -- directory of inactive supporters that made a donation at some point
      +- supporters_001.csv
      +- supporters_002.csv
      +- . . .
  +- groups
      +- groups_001.csv
      +- groups_002.csv
      +- . . .
  +- donations
      +- donations_001.csv
      +- donations_002.csv
      +- ---
```

The CSV files contain the correct information in comma-separated file format.  The CSV headers are guaranteed to
be in the same order from run to run.

The app stores up to 50,000 data records into each file.

The app skips existing files in numerical order.  For example `groups_002.csv` is created if `groups_001.csv` does 
not exxist.

### Performance
You can expect Salsa to provide records at a rate of about 10,000 per minute.

A PC or Mac will get really slow overall.  Run this off-peak if you expect to do your day job.

I use a small AWS instance in an IDE.  That does the work and let's my little old MacBook have some breathing room.

These are observed estimates.  YMMV.

### TODO

1. Retrieve inactive supporters that have donation history.
1. (Optional) Retrieve action names and supporters.
1. (Optional) Retrieve tags and supporters/donations/other.

### Tools
Here's a command that you can use on Linux to see the number of records every minute.

```bash
while true
do
  echo -n `date +%T`
  wc -l data/*.csv | grep total
  sleep 60
done
```

Run that in a separate window.  You should see output like this.
```
17:12:06   131306 total
17:13:06   142874 total
17:14:07   153435 total
17:15:07   163983 total
17:16:07   173560 total
```
## Questions?  Comments?
Use the [Issues](https://github.com/salsalabs/exporter/issues) link in the repository.  Don't waste your time bothering 
the nice folks at Salsalabs Support.

## License
See the LICENSE file in this directory.
