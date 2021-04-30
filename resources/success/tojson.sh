#!/bin/bash

for f in *.yaml
do
  yaml2json $f | jq . > "${f%.*}.json"
done