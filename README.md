# contactfilter

A command-line utility written in Rust to filter large CSV files of contacts. It keeps records based on a specified list of countries, allows prioritizing one country's contacts to appear first, and can limit the final output to a maximum number of records.

## Features

*   **Country-based Filtering**: Filters a CSV file based on a newline-separated list of countries.
*   **Priority Country**: Ensures all contacts from a specific country (e.g., "United States") are placed at the top of the output file.
*   **Size Limiting**: Truncates the final list to a specified maximum number of contacts.
*   **Cross-Platform**: Builds and runs on Linux, macOS, and Windows.
*   **Performant**: Written in Rust for fast processing of large files.

## Installation

There are two ways to install and use `contactfilter`.

### From GitHub Releases (Recommended)

Pre-compiled binaries for Linux, macOS, and Windows are available on the project's **GitHub Releases page**. This is the easiest way to get started.

1.  Go to the latest release.
2.  Download the appropriate asset for your operating system (e.g., `contactfilter-linux-x86_64`, `contactfilter-macos-x86_64`, or `contactfilter-windows-x86_64.exe`).
3.  (Optional for Linux/macOS) Make the binary executable: `chmod +x ./contactfilter-linux-x86_64`

### From Source (For Developers)

If you have the Rust toolchain installed, you can build the project from source.

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/dbehnke/contactfilter.git
    cd contactfilter
    ```

2.  **Build the project**:
    ```bash
    cargo build --release
    ```
    The executable will be located at `target/release/contactfilter`.

## Usage

The program requires a source CSV, a filter file, and a path for the output file.

### Arguments

`contactfilter <INPUT_CSV> <FILTER_FILE> <OUTPUT_CSV> [OPTIONS]`

*   `<INPUT_CSV>`: Path to the input CSV file.
*   `<FILTER_FILE>`: Path to a plain text file with one country per line. An example `countries.txt` is included in this repository and in each release.
*   `<OUTPUT_CSV>`: Path for the new, filtered output CSV file.

### Options

*   `--priority-country <COUNTRY>`: The country to prioritize, ensuring its contacts appear first in the output. [default: `United States`]
*   `--limit <LIMIT>`: The maximum number of contacts to include in the final output. [default: `50000`]
*   `-h, --help`: Print help information.
*   `-V, --version`: Print version information.

### Example

First, create your filter file (or use the one provided in the release assets).

**`countries.txt`:**
```text
United States
Canada
United Kingdom
Australia
New Zealand
```

Then, run the program. The following example uses a downloaded binary on Linux.

```bash
# Basic usage
./contactfilter-linux-x86_64 \
  /path/to/your/Baofeng_DM-32UV_Everything-180days-20250803.csv \
  countries.txt \
  filtered_contacts.csv

# Advanced usage: Prioritize Canada and limit the output to 10,000 contacts
./contactfilter-linux-x86_64 \
  input.csv \
  countries.txt \
  filtered_canadian_contacts.csv \
  --priority-country "Canada" \
  --limit 10000
```

This will read the source CSV, filter it using the countries listed in `countries.txt`, and create a new `filtered_contacts.csv` file in the current directory.