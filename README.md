## Classic-to-Engage Exporter

Go application to export the usual data necessary to convert a client form Classic to Engage.
* Active supporter data
* Inactive supporters that are donors
* Selected group names and supporter emails.
* (Future) Selected tag and supporter emails.
* (Future) Actions and supporter emails.
* (Future) Events and supporter emails.

All output files are CSV.

## Read this!

This is the "Here to Serve" version of engexport.  The process
that reads and transmogrifies data does a lot of special
work that is only for Here to Serve.  

If you're doing another org and you have this version installed, then go get "master".

```bash
git checkout master
```

## Prerequisites
1. A [recent version of Go](https://golang.org/doc/install) installed.  You should use the googles
if you are installing on Windows.  Really. Trust me on this.

1. The correct Go directory structure.  Believe it or not, this is _very_ inportant.  Here's a sample.
```
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```
## Installation
```
go get github.com/salsalabs/engexport
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
* `host` is the [API Host](https://help.salsalabs.com/hc/en-us/articles/115000341773-Salsa-Application-Program-Interface-API-#api_host)
* `email` is the email address that you use to login
* `password` is the password that you use to login

*Remember to remove this file after you're done!*

## Execution

### Help

```
usage: engexport --login=LOGIN [<flags>] <command> [<args> ...]

Classic-to-Engage exporter.

Flags:
  --help          Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN   YAML file with login credentials
  --dir="./data"  Directory to use to store results
  --start=0       start processing at this offset

Commands:

  help [<command>...]
    Show help.

  supporters all
    process all supporters

  supporters active
    process active supporters

  supporters inactive all
    process all inactive supporters

  supporters inactive donors
    process inactive supporters with donation history

  groups
    process groups for active supporters

  donations
    process donations for all active and inactive supporters
```
### Extraction commands for Here to Serve
```
go run cmd/main.go --login YOUR_YAML_FILE supporters all
go run cmd/main.go --login YOUR_YAML_FILE donations
go run cmd/main.go --login YOUR_YAML_FILE groups
```
### Other commands
Other commands exist but do not apply to Here to Serve.

## Output
Output goes to a directory of your choosing.  The default is `./data`.  The output
 directory contains one or more files for each of the exports that the app runs.
 
 The filenames have a sequence number.  The sequence number is there to avoid confusion in the application.  There is no implied or specified order of the output files.
 
 Each file will contain, at most,
 50,001 records.  There will be one record for the CSV header and up to 50,000 data records.

The CSV headers are guaranteed to be in the same order from run to run.  

The filenames indicate the kind of data that they contain.

|Name|Contents|
| --- | --- |
|`supporter_NNN.csv`|supporter records|
|`inactive_supporter_NNN.csv`|inactivesupporters|
|`inactive_donor_NNN.csv`|inactive supporters that have donation history|
|`donation_NNN.csv`|donations|
|`groups_NNN.csv`|(group, email) duples|

### Performance
* You can expect Salsa to provide records at a rate of about 10,000 per minute.
* A PC or Mac will get really slow overall.  Run this off-peak if you expect to do your day job.
* I use a small AWS instance in an IDE.  That does the work and lets my little old MacBook have some breathing room.
* These are observed estimates.  YMMV.

### TODO

1. One pass through the supporter data that genrerates all of the requied CSV files.
2. One pass through the database that updates everyting in Engage automatically
1. Retrieve donations just for active supporters.
1. (Optional) Retrieve action names and supporters.
1. (Optional) Retrieve tags and supporters/donations/other.
1. (Optional) Retrieve custom fields separately using `supporter_KEY` and `Email`
as identifiers.
1. (Whishful thinking) One pass through the database that updates everyting in Engage automatically

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
