## Classic-to-Engage Exporter

Go application to export the usual data necessary to convert a client form Classic to Engage.

-   Active supporter data
-   Inactive supporters that are donors
-   Ggroup names and supporter emails.
-   Selected tag and supporter emails.
-   Actions and supporter emails.
-   Events and supporter emails.
-   Blast email statistics.
-   Supporter email statistics.

All output files are CSV.

## Prerequisites

1.  A [recent version of Go](https://golang.org/doc/install) installed.  You should use the googles
    if you are installing on Windows.  Really. Trust me on this.

2.  The correct Go directory structure.  Believe it or not, this is _very_ inportant.  Here's a sample.

```text
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```
Make sure that `$(HOME)/bin` is in your PATH environment variable.

## Installation

```bash
go get github.com/salsalabs/engexport
go install
go build -o engexport cmd/main.go
mv engexport ~/go/bin
```
When installation is done, you'll have an application named `engexport` that can be executed anywhere on your system.

## Setup

Put your Salsa Classic login credentials into a YAML file.

```yaml
host: hostname
email: you@yours.com
password: super-secret-password
```

Where

-   `host` is the [API Host](https://help.salsalabs.com/hc/en-us/articles/115000341773-Salsa-Application-Program-Interface-API-#api_host)
-   `email` is the email address that you use to login
-   `password` is the password that you use to login

_Remember to remove this file after you're done!_

## Configuration

The default behavior for this app is to export Classic data for Engage.
Only the standard supporter fields are exported.  Custom fields are ignored.

You can change this behavior to

-   Change the field names in the CSV file.
-   Change the fields that are exported.
-   Add custom fields.
-   Remove fields that the client doesn't need.
-   Set up an export to another system besides Engage.

The standard field mappings and headings are stored in `schema.yaml`.  YOu can create new file mappings and headings by copying `schema.haml` and then editing the copy.

(You can also edit `schema.yaml`, but that would Not Be A Good Thing.)

The comments in `schema.yaml` are a good guideline for editing.  Here's a sample from `schema.yaml` for the "groups" table.

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
    # Key maps.  A record has a primary key.  If the primary
    # key is in this list of keys, then the record is saved.
    # If this list of keys is empty, or does not exist, then
    # the record is saved anyway.
    keymap:
```

Some things you need to know.

1.  `supporter`, `donation` and `groups` _must_ be in column 1.
2.  Indents are four spaces.  Don't use tabs.  Using two spaces also works.  Mostly.
3.  Quotations are required for name with a space.
4.  The name before the colon is for the target system.
5.  The name after the colon must be a Salsa Classic field name.
6.  If the Classic field name starts `supporter.`, then leave it alone.
7.  Headers will appear on the first line of the CSV file in the order shown. Feel free to change the order as you see fit.

If you create your own version of `schema.yaml`, then you can add it to the command line when you invoke the app.  See below.

## Execution

```text
usage: engexport --login=LOGIN [<flags>] <command> [<args> ...]

Classic-to-Engage exporter.

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN           YAML file with login credentials
  --schema="schema.yaml"  Classic table schema.
  --dir="./data"          Directory to use to store results
  --tag=TAG               Retrieve records tagged with TAG
  --start=0               start processing at this offset
  --apiVerbose            each api call and response is displayed if true
  --disableInclude        do not use &include in URLs

Commands:
```
| command | description |
| --- | --- |
| `help`  [<command>...] | Show help. |
| `supporters all` | process all supporters |
| `supporters active` | process active supporters |
| `supporters only_email` | process supporters that have emails |
| `supporters inactive all` | process all inactive supporters |
| `supporters inactive donors ``| process inactive supporters with donation history|
| `groups active` | process groups for active supporters |
| `groups only_email` | process groups for supporters that have emails only |
| `groups all` | process groups for all supporters, even ones without emails |
| `donations` | process donations |
| `tags` | process tags as groups |
| `actions` | process supporters and actions |
| `events` | process supporters and events |
| `contact_history` | contact history for all supporters |
| `email_statistics` | email statistics for all supporters |


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

### Notes

If the app crashes spectacularly with this message

`panic: invalid response code 555 `

then the number of fields in the `&include=` in the URL going to Salsa has too many items.  The fix is to use `--disableInclude`:

`engexport --disableInclude --login (etc.)`

Doing that adds some overhead since the full record is being returned and not just parts.  The crashing will stop, though, so there's an upside.

## Output

Output goes to a directory of our choosing.  The default is `./data`.  The output
 directory contains one or more files for each of the exports that the app runs.
 The filenames have a sequence number in them.  Each file will contain, at most,
 50,001 records.  There will be one record for the CSV header and up to 50,000
 data records.

-   The CSV files contain information in comma-separated file format.
-   The CSV headers are guaranteed to be in the same order from run to run.
-   The CSV files contain `NNN` in their names.  For example `supporter_010.csv`.
-   The app uses `NNN` so that older
    files are not overwritten. Any time that a new file is needed, the app searches
    for the last file in that series, then increments `NNN` to create a new file.
-   The filenames indicate the kind of data that they contain.

| Name                         | Contents                                       |
| ---------------------------- | ---------------------------------------------- |
| `supporters_NNN.csv`          | supporter records                              |
| `inactive_supporters_NNN.csv` | inactivesupporters                             |
| `inactive_donors_NNN.csv`     | inactive supporters that have donation history |
| `donations_NNN.csv`           | donations                                      |
| `groups_NNN.csv`             | (group, email) duples                          |
| `tag_groups_NNN.csv` | (tag, email) duples |
| `inactive_donors_NNN.csv` | inactive supporters that have donations |
| `supporter_actions_NNN.csv` | supporter emails, action names, action reference names|
| `supporter_events_NNN.csv` | supporter emails, event names, event titles, etc. |
| `contact_history_NNN.csv` | suppoter emails and contact history |
| `supporter_email_statistics_NNN.csv` | supporter emails and supporter email statistics |

### Performance

-   You can expect Salsa to provide records at a rate of about 10,000 reads per minute.
-   A PC or Mac will get really slow overall.  Run this off-peak if you expect to do your day job.
-   I use a small AWS instance in an IDE.  That does the work and lets my little old MacBook have some breathing room.
-   These are observed estimates.  YMMV.

### Caution!

**Do not** run really big clients through this application on org, org2 or salsa3.  Just don't.  Arrange for Salsa to get the data for you.  Salsa folks: Submmit a JIRA case to a DBA.

### TODO

1.  (Optional) Export custom fields separately. Use `supporter_KEY` and `Email` as the identifiers for importing into Engage.
1.  (Wishful thinking) One pass through the database that updates everyting in Engage automatically.

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

    17:12:06   131306 total
    17:13:06   142874 total
    17:14:07   153435 total
    17:15:07   163983 total
    17:16:07   173560 total

Here's a command that finds all of the groups names and
puts them into a file called `group_names.txt`.  Having a
list of group names can be helpful when creating the groups
on the target system.

```bash
cat data/groups*.csv | \
cut -f 1 -d, | \
sort | \
uniq | \
grep -v ^Group$ \
> group_names.txt
```

## Questions?  Comments?

Use the [Issues](https://github.com/salsalabs/exporter/issues) link in the repository.  Don't waste your time bothering the nice folks at Salsalabs Support.

## License

See the LICENSE file in this directory.
