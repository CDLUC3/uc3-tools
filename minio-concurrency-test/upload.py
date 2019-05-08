#!/usr/bin/env python

import argparse
import logging
import os
import re
import subprocess
import sys

from datetime import datetime

# ############################################################
# Constants

auth_env_vars = ['AWS_ACCESS_KEY_ID', 'AWS_SECRET_ACCESS_KEY']
md5sum_re = re.compile('^([0-9a-f]+) \*(.+)$')

# ############################################################
# Logging

logging.basicConfig(level=logging.INFO, format='%(asctime)-25s %(message)s')
log = logging.getLogger(__file__)


# ############################################################
# Helper functions

def which(cmd_name):
    args = ['which', cmd_name]
    output = check_output(args)
    return output.strip()


def get_md5_cmd():
    try:
        return which('md5sum')
    except subprocess.CalledProcessError:
        return which('gmd5sum')


def md5sum_args(md5sum, filepath):
    return [md5sum, '-b', filepath]


def md5_of(filepath):
    args = md5sum_args(md5_cmd, filepath)
    try:
        result = check_output(args).strip()
        return result.split()[0]
    except Exception as e:
        cmd_str = ' '.join(args)
        raise ValueError("`%s` failed: %s" % (cmd_str, e))


def check_output(args):
    process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    output, _ = process.communicate()
    rc = process.returncode
    if rc != 0:
        raise subprocess.CalledProcessError(rc, ' '.join(args), output)
    return output


# ############################################################
# Globals

md5_cmd = get_md5_cmd()


# ############################################################
# Uploader

class Uploader:
    def __init__(self, bucket, prefix, filepath, endpoint, dryrun, verbose, verify, ssh_to, remote_md5sum):
        self.bucket = bucket
        self.prefix = prefix
        self.filepath = filepath
        self.endpoint = endpoint
        self.dryrun = dryrun
        self.verbose = verbose
        self.verify = verify
        self.ssh_to = ssh_to
        self.remote_md5sum = remote_md5sum

        self.expected_md5 = md5_of(self.filepath)

    def _base_cp_args(self):
        args = ['aws']
        if self.endpoint is not None:
            args.extend(('--endpoint', self.endpoint))
        args.extend(('s3', 'cp'))
        if self.dryrun:
            args.append('--dryrun')
        if not self.verbose:
            args.append('--only-show-errors')
        return args

    def _s3_cp(self, src, dst):
        args = self._s3_cp_args(src, dst)
        cmd_str = ' '.join(args)
        log.info(cmd_str)

        child_out, rc = None, 0
        if self.verbose:
            process = subprocess.Popen(args, stdout=sys.stderr, stderr=subprocess.STDOUT)
            rc = process.wait()
        else:
            process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
            child_out, _ = process.communicate()
            rc = process.returncode
        if rc != 0:
            if child_out is not None:
                log.error(child_out.strip())
            raise ValueError("copying %s to %s failed; %s returned %d" % (src, dst, cmd_str, rc))

        log.info('upload %s to %s complete' % (src, dst))

    def _s3_cp_args(self, src, dst):
        return self._base_cp_args() + [src, dst]

    def _s3_url_for(self, filepath):
        return "s3://%s" % os.path.join(self.bucket, self.prefix, os.path.basename(filepath))

    def _verify(self, s3_url, expected_md5):
        md5sum = self.remote_md5sum if self.ssh_to else md5_cmd

        dst = '$t/%s' % os.path.basename(s3_url)
        s3_cp_args = self._s3_cp_args(s3_url, dst)
        md5_args = md5sum_args(md5sum, dst)

        if self.ssh_to:
            args = ['t=$(mktemp -d)', '&&']
            for var in auth_env_vars:
                val = os.environ.get(var)
                if val is not None:
                    args.append('%s=%s' % (var, val))
            args = args + s3_cp_args + ['&&'] + md5_args
            cmd_str = ' '.join(args)
            cmd_str = "ssh %s '%s'" % (self.ssh_to, cmd_str)
        else:
            args = ['t=$(mktemp -d)', '&&'] + s3_cp_args + ['&&'] + md5_args
            cmd_str = ' '.join(args)

        cmd_sanitized = cmd_str
        for var in auth_env_vars:
            val = os.environ.get(var)
            if val is not None:
                val_sanitized = u'\u2022' * len(val)
                old = '%s=%s' % (var, val)
                new = '%s=%s' % (var, val_sanitized)
                cmd_sanitized = cmd_sanitized.replace(old, new)
        log.info(cmd_sanitized)

        child_out, child_err, rc = None, None, 0
        if self.verbose:
            process = subprocess.Popen(cmd_str, stdout=subprocess.PIPE, stderr=sys.stderr, shell=True)
            child_out, _ = process.communicate()
            rc = process.returncode
        else:
            process = subprocess.Popen(cmd_str, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
            child_out, child_err = process.communicate()
            rc = process.returncode

        output = child_out.strip()
        log.info(output)

        if rc != 0:
            if child_err is not None:
                log.error(child_err.strip())
            raise ValueError("verifying %s failed; %s returned %d" % (s3_url, cmd_str, rc))

        result = md5sum_re.match(output)
        if result is None:
            raise ValueError("can't parse '%s' as md5sum output" % output)
        md5, filepath = result.groups()
        if md5 != expected_md5:
            raise ValueError("digest mismatch for downloaded file %s: expected %s, was %s" % (filepath, expected_md5, md5))

    def upload(self):
        s3_url = self._s3_url_for(self.filepath)
        self._s3_cp(self.filepath, s3_url)
        if self.verify:
            if self.dryrun:
                log.info("(dry run; skipping verification of %s" % s3_url)
            else:
                self._verify(s3_url, self.expected_md5)


def uploaders_from_args():
    parser = argparse.ArgumentParser()

    parser.add_argument("-n", "--dryrun", help="dry run (don't upload anything)", action="store_true")
    parser.add_argument("-e", "--endpoint", help="AWS/Minio endpoint")
    parser.add_argument("-p", "--prefix", help="prefix for uploaded files (defaults to timestamp)")
    parser.add_argument("-v", "--verbose", help="verbose output", action="store_true")
    parser.add_argument("--verify", help="verify after upload", action="store_true")
    parser.add_argument("--ssh-to", help="ssh target for remote verification (default is local)")
    parser.add_argument("--remote-md5sum", help="alternate md5sum binary for remote verification (default is md5sum)")
    parser.add_argument("file", help="file or directory to upload")

    required = parser.add_argument_group('required arguments')
    required.add_argument("-b", "--bucket", required=True, help="S3/Minio bucket")

    args = parser.parse_args()
    bucket = args.bucket
    prefix = (args.prefix or datetime.utcnow().isoformat())
    endpoint = args.endpoint
    dryrun = args.dryrun
    verbose = args.verbose
    verify = args.verify
    ssh_to = args.ssh_to
    remote_md5sum = (args.remote_md5sum or 'md5sum')

    file_to_upload = args.file
    uploaders = []
    if os.path.isfile(file_to_upload):
        uploader = Uploader(bucket, prefix, file_to_upload, endpoint, dryrun, verbose, verify, ssh_to, remote_md5sum)
        uploaders.append(uploader)
    elif os.path.isdir(file_to_upload):
        for dirpath, dirnames, filenames in os.walk(file_to_upload):
            dir_prefix = os.path.join(prefix, dirpath)
            for basename in sorted(filenames):
                filepath = os.path.join(dirpath, basename)
                uploader = Uploader(bucket, dir_prefix, filepath, endpoint, dryrun, verbose, verify, ssh_to, remote_md5sum)
                uploaders.append(uploader)
    return uploaders


# ############################################################
# Helper methods

def check_auth_env_vars():
    for v in auth_env_vars:
        if os.environ.get(v) is None:
            log.warn("$%s not set; falling back to AWS default credentials", v)


# ############################################################
# Main program

def main():
    check_auth_env_vars()
    try:
        uploaders = uploaders_from_args()
        log.info("uploading %d files", len(uploaders))
        success, failure = 0, 0
        for uploader in uploaders:
            try:
                uploader.upload()
                success = success + 1
            except ValueError as e:
                failure = failure + 1
                log.error(e)
        log.info('uploaded %d files successfully' % success)
        if failure > 0:
            log.error('failed to verify %d files' % failure)
    except StandardError as e:
        log.exception(e)
        exit(1)


main()
