#!/bin/bash

declare -a OperatingSystems=("windows" "darwin" "linux")

declare -a Archs=("386")

for o in "${OperatingSystems[@]}"
do
	for a in "${Archs[@]}"
	do
		echo Building $o.$a
		if [ $o = "windows" ]
		then
			ext=".exe"
		else
			ext=""
		fi
		GOOS=$o GOARCH=$a go build -o "./releases/$o.$a/reach$ext"
	done
done

