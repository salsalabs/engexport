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

2.  The correct Go directory structure.  Believe it or not, this is _very important_.  Here's a sample.

```text
$(HOME)
  +- go
    +- bin
    +- pkg
    +- src
```
Make sure that `$(HOME)/go/bin` is in your PATH environment variable.

## Installation

```bash
go get github.com/salsalabs/engexport
go install
go build -o ~/go/bin/engexport cmd/main.go
```
When installation is done, you'll have an application named `engexport` that can be executed anywhere on your system.

## Configuration

The `engexport` app uses a YAML file named 'run.yaml' to provide runtime parameters.
The runtime parameters include
* Salsa Classic API host
* Salsa Classic API login credentials
* Schema selection to show which fields need to be exported
* Output directory
* Runtime options

Contents of 'run.yaml'.

|Field|Description|Notes|Default|
| -- | -- | -- | -- |
|host|[API Host](https://help.salsalabs.com/hc/en-us/articles/115000341773-Salsa-Application-Program-Interface-API-#api_host)||None, required|
|email|Email address to log into Salsa Classic||None, required|
|password|Password to log into SalsaClassic||None, required|
|schema||
||"engage"|Use the schema that we use to export to Engage.  See `public/engage_schema.yaml`.|
||"goodbye"|See `public/goodbye_schema.yaml`||
||filename.yaml|Read the schema from `filename.yaml`.||
|dir|Output directory, created if it doesn't already exist||./data|
|tag|Only retrieve records with this tag|||
|start|Number of records to skip before reading|||
|apiVerbose|If true, then see the API calls and their responses.  Really ugly...||false|
|disableInclude|If true, then retrieve all fields on each read.  If false, only retrieve the necessary fields.||false|
|dumpSchema|If true, then write the as-used schema to "./generated_schema.yaml"||false|
|args|YAML list of table options for the tables to drop.  See cmd/main.go for the authoritative list|||
||"supporters all"|All supporters, both subscribed and unsubscribed.||
||"supporters active"|Supporters with a valid-looking email that are opted in||
||"supporters inactive all"|Supporters with invalid-looking emails or are opted out.||
||"supporters inactive donors"|Inactive supporters that donated||
||"supporters only_email"|Supporters with valid-looking emails||
||"supporters no_email"|Supporters without valid-looking emails||
||"groups active"|Supporters and groups for active supporters||
||"groups all"|Supporters and groups for all supporters||
||"groups only_email"|Supporters and groups for supporters with valid emails||
||"donations active"|Donations for active supporters||
||"donations all"|Donations for all supporters||
||"tags"|Tags and supporters.  Really long.  Do this on an old node and systems will unleash the dogs.||
||"actions"|Actions and supporters.  Another opportunity to be chased by the systems hellhounds...||
||"events"|Events an supporters.||
||"contact_history"|Contact history and supporters.||
||"email_statistics"|Email statistics and supporters.||
||"blast_statistics"|Email blast statistics.||

_Remember to remove this file after you're done!_

A sample `run.yaml` file can be found in `public/sample_run.yaml`.  Your best bet will be to create a working directory for the export, then run this command:

engexport --sample-run-yaml`

That will create a file called `sample_run.yaml` that contains these contents.

```yaml
host: org2.salsalabs.com
email: you@your.org
password: really-long-password
schema: engage
dir: ./data2
tag: example
start: 2000
apiVerbose: false
disableInclude: true
args:
        - supporters all
        - donations all
        - groups all
        - events
        - actions
        - blast_statistics
        - email_statistics
```

## Schema

This app uses a schema file to determine which fields to send to the output.  A schema file is a YAML-formatted file containing definitions for each of the files that the app can export.

The YAML file consists of a number of sections.  Each section describes
the rules to use when exporting a Salsa Classic database table.

```yaml
supporter:
    # Supporter rules here.
donation:
    # Donation rules here.
groups:
    # Rules for the `groups` table here.

# etc.
```

The table rules are composed of
* a `fieldmap` that describes the database table, and
* `headers` that tell `engexport` which fields need to go into the CSV file.

The `fieldmap` contains a line for each field in the Salsa Classic table.  Each line contains the output field name, a colon, and the database table field name.

```yaml
supporter:
    fieldmap:
        "SalsalClassicID": "supporter_KEY"
        "email":           "Email"
        "title":           "Title"
        "firstName":       "First_Name"
        "middleName":      "MI"
        "lastName":        "Last_Name"
```

The `headers` is a simple list of fields that go into the CSV file.  The
output fields will appear in the order that you specify them.


```yaml
supporter:
    fieldmap:
        . . .
    headers:
        - "SalsalClassicID"
        - "email"
        - "title"
        - "firstName"
        - "lastName"
```

Here's a sample of a CSV file containing supporter records for the example schema for supporter.  Note that middleName is not in the export because "middleName" is not in the list of headers.
```csv
SalsaClassicID,email, title,firstName,lastName
123456789,bob@johnson.bizi,Mr,Bob,Johnson
123456788,carol@johnson.bizi,Ms,Carol,Johnson
123456787,ted@johnson.bizi,Mr,Theodore,Johnson
123456786,alice@johnson.bizi,Mrs,Alycia,Johnson

```

Some things you need to know.

1.  `supporter`, `donation` and `groups` _must_ be in column 1.
2.  Indents are four spaces.  Don't use tabs.  Using two spaces also works.  Mostly.
3.  Quotations are required for name with a space.
4.  `engexport` automatically adds all supporter custom fields to the schema.  If you don't need all custom fields, then
    1. Set "dumpSchema" to true in `run.yaml`.
    2. Run the app once.
    3. There will be a file named `generated_schema.yaml` in the current directory.
    4. Rename that file to something useful (e.g. "your_org.yaml")
    5. Edit the file and remove all of the unwanted custom fields.
    6. Put the schema's filename ("your_org.yaml") into the "schema" field of `run.yaml`.
    7. Set "dumpSchema" to false in `run.yaml`.
    7. Run again.
4.  The name before the colon is for the target system.
5.  The name after the colon must be a Salsa Classic field name.
6.  If the Classic field name starts `supporter.`, then leave it alone.
7.  Headers will appear on the first line of the CSV file in the order shown. Feel free to change the order as you see fit.

If you create your own version of `schema.yaml`, then put the schema's filename in the `schema` entry in `run.yaml`.

## Execution

```text
usage: engexport
```

### Notes

If the app crashes spectacularly with this message

`panic: invalid response code 555 `

then the number of fields in the `&include=` in the URL going to Salsa has too many items.  The fix is to set `disableInclude` to false in `run.yaml`

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
-   These are observed estimates.  YMMV.
-   See, observe and obey the next section.

### Caution!

**Do not** run really big clients through this application on org, org2, wfc or salsa3.  Just don't.  Arrange for Salsa to get the data for you.  Salsa folks: Submmit a JIRA case to a DBA.  Non-Salsa folks: submit a request for the data via an email to support@salsalabs.com

### TODO
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
