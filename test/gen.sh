#!/bin/sh

rm -rf *_gen
go build -o schema-generate_gen ../cmd/schema-generate/
for f in *.json; do
    name=`echo $f |sed s/.json//`
    dir=${name}_gen
    echo $f $name $dir
    mkdir -p $dir
    args=$(jq '.__test_args__' -r $f)
    if [ "$args" = "null" ]; then
        args=""
    fi
    echo $args
    ./schema-generate_gen -p $name -o $dir/generated.go $args $f
done
