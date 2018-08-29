# logs-converter

[![pipeline status](https://gitlab.com/oleg.balunenko/logs-converter/badges/master/pipeline.svg)](https://gitlab.com/oleg.balunenko/logs-converter/commits/master)

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/oleg.balunenko/logs-converter)](https://goreportcard.com/report/gitlab.com/oleg.balunenko/logs-converter)

[![coverage report](https://gitlab.com/oleg.balunenko/logs-converter/badges/master/coverage.svg)](https://gitlab.com/oleg.balunenko/logs-converter/commits/master)

The converter will parse files with different log formats and according
on their basis insert MongoDB documents with a monotonous structure.

## How to run it

1. Install Mongo (oficial installation guides: <https://docs.mongodb.com/manual/installation/)>
2. Run mongo

    ```bash
      mongod
    ```

3. Update `config.toml` file in the root of repository with actual parameters and save it (see Configuration)
4. a) - From root of repository run

    ```bash
      go build
    ```
    b) Download latest artifacts [![artifacts](https://img.shields.io/badge/artifacts-download-blue.svg)](https://gitlab.com/oleg.balunenko/logs-converter/-/jobs/artifacts/master/download?job=Build+Application)

5. Run tool

    ```bash
      .logs-converter
    ```

## Configuration

Tool could be configured in 3 ways:
*run with flags
*config file
*einvironment variables

### Flags

```text
  -dropdb
        if true - will drop whole collection before starting to store all logs (default true)
  -loglevel
        LogLevel level: All, Debug, Info, Error, Fatal, Panic, Warn (default Debug)
  -logsfileslist
         (default map[])
  -logsfileslistjson
        JSON with list of all files that need to be looked at and converted
                                                example of JSON:
                                                        {
                                                                "/log1.txt":"first_format",
                                                                "/dir/log2.log":"second_format",
                                                                "/dir2/log3.txt":"first_format"
                                                        }
                                 (default {"testdata/testfile1.log":"second_format","testdata/dir1/testfile2.log":"first_format"})
  -mongocollection
        Mongo DB collection (default logs)
  -mongodb
        Mongo DB name (default myDB)
  -mongopassword
        MongoDB Password
  -mongourl
        Mongo URL (default localhost:27017)
  -mongousername
        MongoDB Username
```

### TOML`config.toml` update following parameters to what you need

***LogLevel** - stdout log level: All, Debug, Info, Error, Fatal, Panic, Warn (default Debug)
***LogsFilesListJSON** - JSON with list of all files that need to be looked at and converted
***MongoURL** - Mongo URL (default localhost:27017)
***MongoDB** - Mongo DB name (default myDB)
***MongoCollection** - Mongo DB collection (default logs)
***MongoUsername** - Mongo DB Username
***MongoPassword** - Mongo DB password
***DropDB** - if true - will drop whole collection before starting to store all logs

example of `config.toml`:

```toml
LogLevel="Debug"
LogsFilesListJSON='{"testdata/testfile1.log":"second_format","testdata/dir1/testfile2.log":"first_format"}'
MongoURL="localhost:27017"
MongoDB="myDB"
MongoCollection="logs"
MongoUsername=""
MongoPassword=""
DropDB=false
```

### environment variables

export following environment variables with your values

example:

```bash

   export LOGSCONVERTER_DROPDB=false
   export LOGSCONVERTER_LOGLEVEL="Info"
   export LOGSCONVERTER_LOGSFILESLISTJSON='{"testdata/testfile1.log":"second_format","testdata/dir1/testfile2.log":"first_format"}'
   export LOGSCONVERTER_MONGOCOLLECTION="logs"
   export LOGSCONVERTER_MONGODB="myDB"
   export LOGSCONVERTER_MONGOURL="localhost:27017"
   export LOGSCONVERTER_MONGOUSERNAME=""
   export LOGSCONVERTER_MONGOPASSWORD=""
```