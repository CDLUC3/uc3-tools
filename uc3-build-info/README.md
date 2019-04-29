# uc3-build-info

Tools for gathering UC3 Jenkins/Maven build information.

## Invocation

Invocation is for the form

```
uc3-build-info <command> [flags] [URL]
```

where `<command>` is one of:

- [`deps`](#deps): List internal Maven dependencies
- [`jobs`](#jobs): List Jenkins jobs
- [`poms`](#poms): List Maven poms

and `URL` is the URL of the Jenkins server. (If not specified, `URL` defaults to `http://builds.cdlib.org/`.)

Note that to access private Git repositories requires a
[GitHub API token](https://github.com/settings/tokens),
which can be passed with the `-t` flag. 

## Shared flags

#### `-t`/`--token TOKEN`

Passes a [GitHub API token](https://github.com/settings/tokens) for accessing private repositories.

#### `-v`/`--verbose`

Verbose output. This both provides more information in the main table (to stdout) and provides
progress information and warning messages to stderr.

## Commands

### `deps`

List internal Maven dependencies

#### Usage

```
uc3-build-info deps [flags]
```

#### Flags

| Short form | Flag           | Description                                           |
| :---       | :---           | :---                                                  |
|            | --all          | list all dependencies                                |
| -a,        | --artifacts    | list Maven artifact dependencies                      |
|            | --expand       | expand tables to rows                                 |
| -f,        | --full-sha     | don't abbreviate SHA hashes in URLs                   |
| -j,        | --jobs         | list Jenkins job dependencies                         |
| -p,        | --poms         | list Maven pom dependencies                           |
| -t,        | --token string | GitHub API token (https://github.com/settings/tokens) |
|            | --tsv          | tab-separated output (default is fixed-width)         |
| -v,        | --verbose      | verbose output                                        |

#### Example



### `jobs`

List Jenkins jobs

#### Usage

```
uc3-build-info jobs [flags]
```

#### Flags

| Short form | Flag           | Description                                           |
| :---       | :---           | :---                                                  |
|            | --api-url      | show Jenkins API URLs                                 |
| -a,        | --artifacts    | show artifacts from last successful build             |
| -b,        | --build        | show info for last successful build                   |
|            | --config-xml   | show Jenkins config.xml URLs                          |
| -f,        | --full-sha     | don't abbreviate SHA hashes in URLs                   |
| -j,        | --job string   | show info only for specified job                      |
| -p,        | --parameters   | show build parameters                                 |
|            | --poms         | show POMs                                             |
| -r,        | --repositories | show repositories                                     |
| -t,        | --token string | GitHub API token (https://github.com/settings/tokens) |
|            | --tsv          | tab-separated output (default is fixed-width)         |
| -v,        | --verbose      | verbose output                                        |

### `poms`

List Maven poms

#### Usage

```
uc3-build-info poms [flags]
```

#### Flags

| Short form | Flag           | Description                                           |
| :---       | :---           | :---                                                  |
| -a,        | --artifacts    | list POM artifacts                                    |
| -d,        | --deps         | list POM dependencies                                 |
| -f,        | --full-sha     | don't abbreviate SHA hashes in URLs                   |
| -j,        | --job string   | show info only for specified job                      |
| -u,        | --pom-urls     | list URL used to retrieve POM file                    |
| -t,        | --token string | GitHub API token (https://github.com/settings/tokens) |
|            | --tsv          | tab-separated output (default is fixed-width)         |
| -v,        | --verbose      | verbose output                                        |
