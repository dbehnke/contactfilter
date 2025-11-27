#!/bin/bash
set -e

echo "Building contactfilter..."
go build -o contactfilter

echo "Creating dummy input..."
echo "No.,ID,Repeater,Name,City,Province,Country,Remark,Type,Alert Call" > test_input.csv
echo "1,1234567,Repeater1,John Doe,City1,State1,USA,Remark1,Type1,Alert1" >> test_input.csv

echo "USA" > countries.txt

for fmt in baofeng-dm32uv anytone opengd77 db25d; do
    echo "Testing format: $fmt"
    ./contactfilter --input-csv test_input.csv --output-csv "output_$fmt.csv" --radio "$fmt" --filter-file countries.txt
    if [ -f "output_$fmt.csv" ]; then
        echo "Generated output_$fmt.csv"
        head -n 2 "output_$fmt.csv"
    else
        echo "Failed to generate output_$fmt.csv"
        exit 1
    fi
done

echo "Cleaning up..."
rm contactfilter test_input.csv output_*.csv
echo "Verification complete!"
