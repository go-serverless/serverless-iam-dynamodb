#!/usr/bin/env bash

dep ensure

echo "************************************************"
echo "* Compiling functions to bin/handlers/ ...      "
echo "************************************************"

rm -rf bin/

cd src/handlers/
for f in *.go; do
  filename="${f%.go}"
  if GOOS=linux go build -o "../../bin/handlers/$filename" ${f}; then
    echo "* Compiled $filename"
  else
    echo "* Failed to compile $filename!"
    exit 1
  fi
done

echo "************************************************"
echo "* Formatting Code ...                           "
echo "************************************************"
go fmt

echo "************************************************"
echo "* Build Completed                               "
echo "************************************************"