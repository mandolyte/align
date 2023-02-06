#!/bin/sh

echo "Must be in a directory of chapter JSON files from Translation Core"
rm -f *.csv
for i in *.json 
do 
    echo "Working on $i"
    go run ../../../convert.go -i $i -o $(basename $i .json).csv -b "Titus" -c $(basename $i .json)
done 

echo "Combining into one file for this book"
catcsv -o all.csv 1.csv 2.csv 3.csv