# uc3-system-info

A tool for generating UC3 system info reports

## Input

An inventory file (see [tests/testdata/uc3-inventory.txt](tests/testdata/uc3-inventory.txt) for an example).

## Usage

```
Usage:
  uc3-system-info hosts <FILE> [flags]

Flags:
  -f, --format string    output format (tsv, csv, md) (default "tsv")
  -s, --service string   filter to specified service
      --header           include header
      --footer           include footer
  -h, --help             help for hosts
```

Examples:

- Generate a full report in default TSV format, without header or footer:

  ```
  uc3-system-info hosts uc3-inventory.txt
  ```

  Outputs:

  ```
  dash    dev             uc3-dash2-dev   uc3-dash2-dev.cdlib.org dash-aws-dev.cdlib.org,dash-dev.cdlib.org,dash-ucla-dev.cdlib.org,dash2-crossref-dev.cdlib.org,dash2-dev.cdlib.org,datashare-dev.cdlib.org,oneshare-aws-dev.cdlib.org,oneshare-dev.cdlib.org,oneshare2-dev.cdlib.org,uc3-datashare-dev.cdlib.org
  dash    dev     solr    uc3-dash2solr-dev       uc3-dash2solr-dev.cdlib.org     dash2solr-dev.cdlib.org
  dash    stg             uc3-dash2-stg   uc3-dash2-stg.cdlib.org dash-stg.cdlib.org,dash-ucla-stg.cdlib.org,dash2-aws-stg.cdlib.org,dash2-crossref-stg.cdlib.org,dash2-stg.cdlib.org,datashare-stg.cdlib.org,oneshare-stg.cdlib.org,oneshare2-stg.cdlib.org,uc3-datashare-stg.cdlib.org

  (...)

  Generated 2019-02-13T16:01:42-08:00
  ```

- Generate a full report in Markdown format with header and footer:

  ```
  uc3-system-info hosts uc3-inventory.txt --format md --header --footer
  ```

  Outputs:

  | Service | Environment | Subsystem | Name | FQDN | CNAMEs |
  | :--- | :--- | :--- | :--- | :--- | :--- |
  | dash | dev |  | uc3-dash2-dev | uc3-dash2-dev.cdlib.org | dash-aws-dev.cdlib.org<br/>dash-dev.cdlib.org<br/>dash-ucla-dev.cdlib.org<br/>dash2-crossref-dev.cdlib.org<br/>dash2-dev.cdlib.org<br/>datashare-dev.cdlib.org<br/>oneshare-aws-dev.cdlib.org<br/>oneshare-dev.cdlib.org<br/>oneshare2-dev.cdlib.org<br/>uc3-datashare-dev.cdlib.org |
  | dash | dev | solr | uc3-dash2solr-dev | uc3-dash2solr-dev.cdlib.org | dash2solr-dev.cdlib.org |
  | dash | stg |  | uc3-dash2-stg | uc3-dash2-stg.cdlib.org | dash-stg.cdlib.org<br/>dash-ucla-stg.cdlib.org<br/>dash2-aws-stg.cdlib.org<br/>dash2-crossref-stg.cdlib.org<br/>dash2-stg.cdlib.org<br/>datashare-stg.cdlib.org<br/>oneshare-stg.cdlib.org<br/>oneshare2-stg.cdlib.org<br/>uc3-datashare-stg.cdlib.org |

  (etc.)

- Generate a comma-separated report on Dash hosts only, with header:

  ```
  uc3-system-info hosts uc3-inventory.txt --format csv --service dash --header
  ```

  Outputs:

  ```
  Environment,Subsystem,Name,FQDN,CNAMEs
  dev,,uc3-dash2-dev,uc3-dash2-dev.cdlib.org,dash-aws-dev.cdlib.org;dash-dev.cdlib.org;dash-ucla-dev.cdlib.org;dash2-crossref-dev.cdlib.org;dash2-dev.cdlib.org;datashare-dev.cdlib.org;oneshare-aws-dev.cdlib.org;oneshare-dev.cdlib.org;oneshare2-dev.cdlib.org;uc3-datashare-dev.cdlib.org
  dev,solr,uc3-dash2solr-dev,uc3-dash2solr-dev.cdlib.org,dash2solr-dev.cdlib.org
  stg,,uc3-dash2-stg,uc3-dash2-stg.cdlib.org,dash-stg.cdlib.org;dash-ucla-stg.cdlib.org;dash2-aws-stg.cdlib.org;dash2-crossref-stg.cdlib.org;dash2-stg.cdlib.org;datashare-stg.cdlib.org;oneshare-stg.cdlib.org;oneshare2-stg.cdlib.org;uc3-datashare-stg.cdlib.org
  ```

  (etc.)


## Building

The `uc3-system-info` project can be built and installed simply with `go
build` and `go install`, but it also supports [Mage](https://magefile.org).

To install the latest version of Mage:

1. visit their [releases page](https://github.com/magefile/mage/releases),
   download the appropriate binary, and place it in your `$PATH`, or
2. from _outside_ this project directory (`go get` behaves differently when
   run in the context of a module project), execute the following:

   ```
   go get -u -d github.com/magefile/mage \
   && cd $GOPATH/src/github.com/magefile/mage \
   && go run bootstrap.go
   ```

#### Mage tasks:

| Tasks        | Purpose                                                          |
| :---         | :---                                                             |
| `build`      | builds a binary for the current platform                         |
| `buildAll`   | builds a binary for each target platform                         |
| `buildLinux` | builds a linux-amd64 binary (the most common cross-compile case) |
| `clean`      | removes compiled binaries from the current working directory     |
| `install`    | installs in $GOPATH/bin                                          |
| `platforms`  | lists target platforms for buildAll                              |

Note that `mage build` is a thin wrapper around `go build` and supports the
same environment variables, e.g. `$GOOS` and `$GOARCH`.
