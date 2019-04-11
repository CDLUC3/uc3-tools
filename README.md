# mrt-build-info

Tools for gathering Merritt build information.

## Invocation

Invocation is for the form

```
mrt-build-info <command> [flags] [URL]
```

where `<command>` is one of:

- [`jobs`](#jobs): list Jenkins jobs

and `URL` is the URL of the Jenkins server. (If not specified, `URL` defaults to `http://builds.cdlib.org/`.)

## Commands

### `jobs`

List Jenkins jobs.

By default, the `jobs` command simply lists all Jenkins jobs by name. The
flags below can be set to provide more information for each job, taken from
the last successful build.

| Short form | Flag             | Description                         |
| ---        | ---              | ---                                 |
| `-a`       | `--artifacts`    | list artifacts                      |
| `-b`       | `--build`        | show info for last successful build |
| `-r`       | `--repositories` | list repositories                   |

If any of these flags are set, output will be in the form of a tab-separated
table, with header.

The `jobs` command supports the following additional flags:

| Short form | Flag           | Description                    |
| ---        | ---            | ---                            |
| `-h`       | `--help`       | help for jobs                  |
| `-v`       | `--verbose`    | verbose output                 |
| `-l`       | `--log-errors` | log non-fatal errors to stderr |


