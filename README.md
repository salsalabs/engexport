## Classic-to-Engage Exporter

Go application to export the usual data necessary to convert a client form Classic to Engage.
* Active supporter data
* Inactive supporters that are donors
* Selected group names and supporter emails.
* (Future) Selected tag and supporter emails.
* (Future) Actions and supporter emails.
* (Future) Events and supporter emails.

All output files are CSV.

## Prerequisites

1. A [recent version of Go](https://golang.org/doc/install) installed.  You should use the googles
if you are installing on Windows.  Really. Trust me on this.

1. The correct Go directory structure.  Believe it or not, this is _very_ inportant.  Here's a sample.

```text
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```

## Installation

```bash
go get github.com/salsalabs/engexport
go install
go build -o engexport cmd/main.go
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

## Configuration

The default behavior for this app is to export Classic data for Engage.
Only the standard supporter fields are exported.  Custom fields are ignored.

You can change this behavior to

* Change the field names in the CSV file.
* Change the fields that are exported.
* Add custom fields.
* Remove fields that the client doesn't need.
* Set up an export to another system besides Engage.

The standard field mappings and headings are stored in `schema.yaml`.  YOu can create
new file mappings and headings by copying `schema.haml` and then editing the copy.
(You can also edit `schema.yaml`, but that would Not Be A Good Thing.)

The comments in `schema.yaml` are a good guideline for editing.  Here's a sample
from `schema.yaml` for the "groups" table.

```yaml
groups:
    # Field map.  The column on the left is the target.
    # The list of headers must come of out the left column.
    #
    # The column on the right is the Salsa Classic groups
    # table field name *or* fields from a joined supporter
    # record.  You can use both standard fields and custom
    # fields.  Changing the joined supporter fields is 
    # Not A Good Idea.
    #
    # This file is the default mapping for transferring
    # groups information from Salsa Classic to Engage.
    fieldmap:
        "Group": "Group_Name"
        "Email": "supporter.Email"
    # Headers for the CSV file.  Headers will appear at the
    # top of the CSV file in this order.  Note that headers
    # *must* come from the left column of the field map.
    headers:
        - "Group"
        - "Email"
```

Some things you need to know.

1. `supporter`, `donation` and `groups` *must* be in column 1.
1. Indents are four spaces.  Don't use tabs.  Using two spaces also works.  Mostly.
2. Quotations are required for name with a space.
1. The name before the colon is for the target system.
1. The name after the colon must be a Salsa Classic field name.
1. If the Classic field name starts `supporter.`, then leave it alone.
2. Headers will appear on the first line of the CSV file in the order shown. Feel free to change the order as you see fit.

If you cceate your own version of `schema.yaml`, then you can add it to the command 
line when you invoke the app.  See below.

## Execution

```bash
usage: ./engexport --login=LOGIN [<flags>] <command> [<args> ...]


Classic-to-Engage exporter.

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN           YAML file with login credentials
  --schema="schema.yaml"  Classic table schema.
  --dir="./data"          Directory to use to store results
  --tag=TAG               Retrieve records tagged with TAG
  --start=0               start processing at this offset

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
    process donations for active and inactive supporters

```
### Examples

#### Active supporters

```bash
go run cmd/main.go --login YOUR_YAML_FILE supporters active
```

#### Donations by all active and inactive supporters.

This is the most common
variant that clients ask for.  See "TODO" section for others.

```bash
go run cmd/main.go --login YOUR_YAML_FILE donations
```

#### Group names and emails for all active supporters.

```bash
go run cmd/main.go --login YOUR_YAML_FILE groups
```

#### Inactive supporters.

```bash
go run cmd/main.go --login YOUR_YAML_FILE supporters inactive
```

#### Inactive supporters that have donation history.

```bash
go run cmd/main.go --login YOUR_YAML_FILE supporters inactive donors
```

## Output

Output goes to a directory of our choosing.  The default is `./data`.  The output
 directory contains one or more files for each of the exports that the app runs.
 The filenames have a sequence number in them.  Each file will contain, at most,
 50,001 records.  There will be one record for the CSV header and up to 50,000
 data records.

* The CSV files contain information in comma-separated file format.
* The CSV headers are guaranteed to be in the same order from run to run.
* The CSV files contain `NNN` in their names.  For example `supporter_010.csv`.
* The app uses `NNN` so that older
files are not overwritten. Any time that a new file is needed, the app searches
for the last file in that series, then increments `NNN` to create a new file.
* The filenames indicate the kind of data that they contain.

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
1. Retrieve donations just for active supporters.
2. Retrieve supporters who are not donors.
1. (Optional) Retrieve action names and supporters.
1. (Optional) Retrieve tags and supporters/donations/other.
1. (Optional) Retrieve custom fields separately using `supporter_KEY` and `Email`
as identifiers.
1. (Wishful thinking) One pass through the database that updates everyting in Engage automatically

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
