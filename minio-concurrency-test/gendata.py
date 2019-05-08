#!/usr/bin/env python

import argparse
import logging
import math
import os

# ############################################################
# Constants

KB = 1024
cwd = os.getcwd()

# ############################################################
# Logging

logging.basicConfig(level=logging.INFO, format='%(asctime)-25s %(message)s')
log = logging.getLogger(__file__)


# ############################################################
# Generator

class Generator:
    def __init__(self, data_dir, num, size_mb, force):
        self.data_dir = data_dir
        self.num = num
        self.size_mb = size_mb
        self.force = force
        self.numwidth = int(math.ceil(math.log10(num)))

    def _filepath(self, index):
        filename = 'file-%s.bin' % (str(index).zfill(self.numwidth))
        filepath = os.path.join(self.data_dir, filename)
        return os.path.relpath(filepath, cwd)

    def _write_file(self, filepath, expected_bytes):
        actual_bytes = 0
        try:
            with open(filepath, 'wb') as outfile:
                for _ in range(0, self.size_mb):
                    outfile.write(os.urandom(KB))
        finally:
            if os.path.exists(filepath):
                actual_bytes = os.path.getsize(filepath)
            log.info("wrote %d of %d bytes to %s", actual_bytes, expected_bytes, filepath)
        if actual_bytes != expected_bytes:
            raise ValueError("expected %d bytes, but only wrote %d" % (expected_bytes, actual_bytes))

    def generate(self):
        expected_bytes = self.size_mb * KB
        log.info("writing %d files of size %d bytes to %s", self.num, expected_bytes, self.data_dir)
        if self.num == 0:
            return
        if not os.path.isdir(self.data_dir):
            os.makedirs(self.data_dir)
        for index in range(0, self.num):
            filepath = self._filepath(index)
            if os.path.isfile(filepath) and not self.force:
                log.warn("skipping existing file %s", filepath)
                continue
            self._write_file(filepath, expected_bytes)


def generator_from_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--dir", help="directory for generated files (default ./data)")
    parser.add_argument("-n", "--num", help="number of files to generate (default 1)", type=int)
    parser.add_argument("-s", "--size", help="size in KiB (default 1)", type=int)
    parser.add_argument("-f", "--force", help="overwrite existing files (default false)", action="store_true")

    args = parser.parse_args()

    if args.num is not None and args.num <= 0:
        raise ValueError("can't write %d files" % args.num)
    if args.size is not None and args.size < 0:
        raise ValueError("can't write file of %d bytes" % args.size)

    return Generator(
        args.dir or os.path.join(cwd, 'data'),
        args.num or 1,
        args.size or 1,
        args.force
    )


# ############################################################
# Main program

def main():
    try:
        generator = generator_from_args()
        generator.generate()
    except ValueError as e:
        log.error(e)
        exit(1)
    except StandardError as e:
        log.exception(e)
        exit(1)


main()
