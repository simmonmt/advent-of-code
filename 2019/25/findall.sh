#!/bin/bash

# Run through all the combinations of dropped items

BINDIR=${HOME}/src/aoc/2019/bazel-bin/25/a/src

rm -fr /tmp/droid-out
mkdir /tmp/droid-out

for i in $(seq 0 255) ; do
  echo $i

  cat ~/src/aoc/2019/25/commands_prefix.txt >/tmp/out
  ${BINDIR}/util/darwin_amd64_stripped/drop --seq $i >>/tmp/out
  echo north >>/tmp/out

  ${BINDIR}/darwin_amd64_stripped/solution --ram ~/src/aoc/2019/25/input.txt --input /tmp/out >/tmp/droid-out/$i
done

