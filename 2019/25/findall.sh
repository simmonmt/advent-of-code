#!/bin/bash
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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

