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
4. From root of repository run
    ```bash
      go build
    ```
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
        if true - will drop whole collection before starting to store all logs
  -loglevel
        LogLevel level: All, Debug, Info, Error, Fatal, Panic, Warn (default Debug)
  -logsfileslistjson
        JSON with list of all files that need to be looked at and converted
                                                example of JSON:
                                                        {
                                                                "/log1.txt":"first_format",
                                                                "/dir/log2.log":"second_format",
                                                                "/dir2/log3.txt":"first_format"
                                                        }
                                                Common JSON schema:
                                                      {
                                                            "log_file_path":"log_format"
                                                      }
  -mongocollection
        Mongo DB collection (default logs)
  -mongodb
        Mongo DB name (default myDB)
  -mongourl
        Mongo URL (default localhost:27017)
```

### TOML`config.toml` update following parameters to what you need

***LogLevel** - stdout log level: All, Debug, Info, Error, Fatal, Panic, Warn (default Debug)
***LogsFilesListJSON** - JSON with list of all files that need to be looked at and converted
***MongoURL** - Mongo URL (default localhost:27017)
***MongoDB** - Mongo DB name (default myDB)
***MongoCollection** - Mongo DB collection (default logs)
***DropDB** - if true - will drop whole collection before starting to store all logs

example of `config.toml`:

```toml
LogLevel="Debug"
LogsFilesListJSON='{"testdata/testfile1.log":"second_format","testdata/dir1/testfile2.log":"first_format"}'
MongoURL="localhost:27017"
MongoDB="myDB"
MongoCollection="logs"
DropDB=true
```

### environment variables

export following environment variables with your values

example:

```bash
   export KAFKADUMP_DROPDB=false
   export KAFKADUMP_LOGLEVEL="Info"
   export KAFKADUMP_LOGSFILESLISTJSON='{"testdata/testfile1.log":"second_format","testdata/dir1/testfile2.log":"first_format"}'
   export KAFKADUMP_MONGOCOLLECTION="logs"
   export KAFKADUMP_MONGODB="myDB"
   export KAFKADUMP_MONGOURL="localhost:27017"
```