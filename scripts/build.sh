#!/usr/bin/env bash

echo "************************************************"
echo "* Formatting Code ...                           "
echo "************************************************"
go fmt src/handlers/*

echo "************************************************"
echo "* Compiling functions to bin/handlers/ ...      "
echo "************************************************"

rm -rf bin/

for f in src/handlers/*.go; do
  filename="${f%.go}"
  if GOOS=linux go build -o "../../bin/handlers/$filename" ${f}; then
    echo "* Compiled $filename"
  else
    echo "* Failed to compile $filename!"
    exit 1
  fi
done

echo "************************************************"
echo "* Build Completed                               "
echo "************************************************"