#!/bin/bash

cat testdata/data |while read file num; do
  NUM=$(./solution <testdata/$file --num_turns=-1 2>/dev/null |tail -1 |awk '{print $6}')
  if [[ "X$NUM" = "X${num}" ]] ; then
    echo PASS $file
  else
    echo FAIL $file
  fi
done
