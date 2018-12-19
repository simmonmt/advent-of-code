#!/bin/bash

cat testdata/expected |while read file rest ; do
    out=$(./solution <testdata/${file} 2>&1)
    if [[ "$out" = "$rest" ]] ; then
	echo $file PASSED
    else
	echo $file FAILED
    fi
done
