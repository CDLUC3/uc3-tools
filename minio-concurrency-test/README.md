# minio-concurrency-test

Test script for Minio concurrency issues.

## Usage

### Generating test data

The `gendata.py` script generates test data in the form of a stream of random bytes.

```
usage: gendata.py [-h] [-d DIR] [-n NUM] [-s SIZE] [-f]

optional arguments:
  -h, --help            show this help message and exit
  -d DIR, --dir DIR     directory for generated files (default ./data)
  -n NUM, --num NUM     number of files to generate (default 1)
  -s SIZE, --size SIZE  size in KiB (default 1)
  -f, --force           overwrite existing files (default false)
```

Example:

```
$ gendata.py -n 5
2019-05-07 16:21:38,497   writing 5 files of size 1024 bytes to /Users/me/data
2019-05-07 16:21:38,497   wrote 1024 of 1024 bytes to data/file-0.bin
2019-05-07 16:21:38,497   wrote 1024 of 1024 bytes to data/file-1.bin
2019-05-07 16:21:38,498   wrote 1024 of 1024 bytes to data/file-2.bin
2019-05-07 16:21:38,498   wrote 1024 of 1024 bytes to data/file-3.bin
2019-05-07 16:21:38,498   wrote 1024 of 1024 bytes to data/file-4.bin
```

### Uploading and verifying files

The `upload.py` script uploads and verifies files. The verification can be run locally (default)
or (to detect timing/consistency issues) on another server via `ssh`.

Authentication uses the AWS environment, or the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

```
usage: upload.py [-h] [-n] [-e ENDPOINT] [-p PREFIX] [-v] [--verify] [--ssh-to SSH_TO] [--remote-md5sum REMOTE_MD5SUM] 
                 -b BUCKET file

positional arguments:
  file                  file or directory to upload

optional arguments:
  -h, --help                           show this help message and exit
  -n, --dryrun                         dry run (don't upload anything)
  -e ENDPOINT, --endpoint ENDPOINT     AWS/Minio endpoint
  -p PREFIX, --prefix PREFIX           prefix for uploaded files (defaults to timestamp)
  -v, --verbose                        verbose output
  --verify                             verify after upload
  --ssh-to SSH_TO                      ssh target for remote verification (default is local)
  --remote-md5sum REMOTE_MD5SUM        alternate md5sum binary for remote verification (default is md5sum)

required arguments:
  -b BUCKET, --bucket BUCKET S3/Minio bucket
```

#### Examples

##### Single file without verification

Uploading a single file to the default endpoint, at a timestamp-based prefix, without verification:

```
$ upload.py -b my-bucket file-00.bin
2019-05-07 15:56:58,464   uploading 1 files
2019-05-07 15:56:58,464   aws s3 cp --only-show-errors file-00.bin s3://my-bucket/2019-05-07T22:56:58.460768/file-00.bin
```

##### Multiple files

Uploading a directory of files to the default endpoint, at a timestamp-based prefix, without verification:

```
$ upload.py -e https://minio.example.edu/ -b my-bucket data
2019-05-07 16:16:40,649   uploading 4 files
2019-05-07 16:16:40,650   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors data/file-00.bin s3://my-bucket/2019-05-07T23:16:40.046575/data/file-00.bin
2019-05-07 16:16:42,280   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors data/file-01.bin s3://my-bucket/2019-05-07T23:16:40.046575/data/file-01.bin
2019-05-07 16:16:43,898   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors data/file-02.bin s3://my-bucket/2019-05-07T23:16:40.046575/data/file-02.bin
2019-05-07 16:16:45,484   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors data/file-03.bin s3://my-bucket/2019-05-07T23:16:40.046575/data/file-03.bin
2019-05-07 16:16:45,484   uploaded 4 files successfully
```

##### Single file with local verification

Uploading a single file to a specified endpoint with local verification:

```
$ upload.py -e https://minio.example.edu/ -b my-bucket --verify file-00.bin
2019-05-07 15:57:49,612   uploading 1 files
2019-05-07 15:57:49,612   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors file-00.bin s3://my-bucket/2019-05-07T22:57:49.608543/file-00.bin
2019-05-07 15:57:50,164   t=$(mktemp -d) && aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors s3://my-bucket/2019-05-07T22:57:49.608543/file-00.bin $t/file-00.bin && /usr/local/bin/gmd5sum -b $t/file-00.bin
2019-05-07 15:57:50,749   a15f5c8e4a42c50a0157dc4368156227 */var/folders/8s/45cxv60949735s2r_5sjh5gh0000gn/T/tmp.CtgigVH2/file-00.bin
2019-05-07 15:57:50,749   uploaded 1 files successfully
```

Note that on an OS X system, `file-00.bin` will use `gmd5sum` from Homebrew [GNU coreutils](https://formulae.brew.sh/formula/coreutils). (It's
not sophisticated enough to use the BSD `md5` command, which has different options and output syntax.) This might also work with MacPorts or
other BSD-based ports packages, but it hasn't been tested.

Note also that the temporary file, in this case `/var/folders/8s/45cxv60949735s2r_5sjh5gh0000gn/T/tmp.CtgigVH2/file-00.bin`,
is not deleted and will need to be cleaned up.

##### Single file with remote verification

Uploading a single file to a specified endpoint with remote verification:

```
$ upload.py -e https://minio.example.edu/ -b my-bucket --verify --ssh-to myaccount@otherhost.example.edu file-00.bin
2019-05-07 16:03:21,126   uploading 1 files
2019-05-07 16:03:21,127   aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors file-00.bin s3://my-bucket/2019-05-07T23:03:21.121731/file-00.bin
2019-05-07 16:03:21,706   ssh myaccount@otherhost.example.edu 't=$(mktemp -d) && AWS_ACCESS_KEY_ID=•••••••••••••••••• AWS_SECRET_ACCESS_KEY=•••••••••••••••••••••••••••••••••••••••• aws --endpoint https://minio.example.edu/ s3 cp --only-show-errors s3://my-bucket/2019-05-07T23:03:21.121731/file-00.bin $t/file-00.bin && md5sum -b $t/file-00.bin'
2019-05-07 16:03:23,250   a15f5c8e4a42c50a0157dc4368156227 */tmp/tmp.SkSF1HQXXj/file-00.bin
2019-05-07 16:03:23,250   uploaded 1 file successfully
```

Note that the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` variables, if present, are forwarded to the remote server
(but not written to stdout).