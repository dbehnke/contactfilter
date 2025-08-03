# contactfilter

A simple command-line utility written in Rust to filter a CSV file of contacts based on a list of countries.

This program reads a source CSV file, checks the "Country" column for each record, and writes the record to a new output CSV file only if the country is present in a specified filter file.

## Dependencies

*   Rust and Cargo
*   `csv` crate
*   `serde` crate

## Setup

1.  **Create a filter file**: Create a plain text file (e.g., `countries.txt`) and list the countries you want to keep, one per line.

    ```text
    United States
    Canada
    United Kingdom
    Australia
    New Zealand
    ```

2.  **Build the project**: Compile the program in release mode for optimal performance.

    ```bash
    cargo build --release
    ```

    The executable will be located at `target/release/contactfilter`.

## Usage

The program takes three required positional arguments and has several optional flags.

### Arguments

`contactfilter <INPUT_CSV> <FILTER_FILE> <OUTPUT_CSV> [OPTIONS]`

*   `<INPUT_CSV>`: Path to the input CSV file.
*   `<FILTER_FILE>`: Path to the country filter file.
*   `<OUTPUT_CSV>`: Path for the new, filtered output CSV file.

### Options

*   `--priority-country <COUNTRY>`: The country to prioritize, ensuring its contacts appear first in the output. [default: `United States`]
*   `--limit <LIMIT>`: The maximum number of contacts to include in the final output. [default: `50000`]
*   `-h, --help`: Print help information.
*   `-V, --version`: Print version information.

### Examples

```bash
# Basic usage
./target/release/contactfilter \
  /path/to/your/Baofeng_DM-32UV_Everything-180days-20250803.csv \
  countries.txt \
  filtered_contacts.csv

# Advanced usage: Prioritize Canada and limit the output to 10,000 contacts
./target/release/contactfilter \
  input.csv \
  countries.txt \
  filtered_canadian_contacts.csv \
  --priority-country "Canada" \
  --limit 10000
```

This will read the source CSV, filter it using the countries listed in `countries.txt`, and create a new `filtered_contacts.csv` file in the current directory.