#!/bin/bash

set -e

rm -f *.deb

drone exec

for f in *.deb; do
  echo "*** FILE: ${f}"
  echo "INFO:"
  dpkg --info ${f}
  echo "CONTENTS:"
  dpkg --contents ${f}
done
