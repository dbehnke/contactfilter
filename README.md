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

The program takes three command-line arguments:

1.  Path to the input CSV file.
2.  Path to the country filter file.
3.  Path for the new, filtered output CSV file.

### Example

```bash
./target/release/contactfilter \
  /path/to/your/Baofeng_DM-32UV_Everything-180days-20250803.csv \
  countries.txt \
  filtered_contacts.csv
```

This will read the source CSV, filter it using the countries listed in `countries.txt`, and create a new `filtered_contacts.csv` file in the current directory.