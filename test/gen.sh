#!/bin/sh

rm -rf *_gen
go build -o schema-generate_gen ../cmd/schema-generate/
for f in *.json; do
    name=`echo $f |sed s/.json//`
    dir=${name}_gen
    echo $f $name $dir
    mkdir -p $dir
    ./schema-generate_gen -p $name -o $dir/gen.go $f
done
