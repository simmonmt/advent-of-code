#!/bin/bash

cat testdata/data |while read file ap num; do
  NUM=$(./solution <testdata/$file --elf_attack_power $ap --num_turns=-1 2>/dev/null |tail -1 |awk '{print $6}')
  if [[ "X$NUM" = "X${num}" ]] ; then
    echo PASS $file
  else
    echo FAIL $file
  fi
done
