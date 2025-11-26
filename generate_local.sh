#!/bin/bash
set -e

# 1. Build the binary (mimicking the build step)
echo "Building contactfilter..."
go build -o contactfilter

# 2. Define variables (mimicking CI env vars)
TIMESTAMP=$(date +%s)
BM_FILE="bm_contacts.csv"
OUTPUT_FILE="Baofeng_DM-32UV_Contacts_${TIMESTAMP}.csv"

# 3. Run the generation logic (mimicking the CI run step)
if [ -f "$BM_FILE" ]; then
    echo "Found Brandmeister contacts file: $BM_FILE"
    echo "Generating '$OUTPUT_FILE'..."
    
    ./contactfilter \
        --merge \
        --brandmeister-csv "$BM_FILE" \
        --filter-file countries.txt \
        --output-csv "$OUTPUT_FILE"
        
    echo "Success! Generated: $OUTPUT_FILE"
else
    echo "Error: $BM_FILE not found."
    exit 1
fi
