# contactfilter

A command-line utility written in Go to filter large CSV files of contacts. It keeps records based on a specified list of countries, allows prioritizing one country's contacts to appear first, and can limit the final output to a maximum number of records.

It also supports **merging** "Last Heard" data from Brandmeister with the full contact details from RadioID.net, allowing you to create a contact list of only recently active users.

## Features

* **Merge Mode**: Combine Brandmeister "Last Heard" activity with RadioID.net details.
* **Country-based Filtering**: Filters a CSV file based on a newline-separated list of countries.
* **Priority Country**: Ensures all contacts from a specific country (e.g., "United States") are placed at the top of the output file.
* **Size Limiting**: Truncates the final list to a specified maximum number of contacts.
* **Cross-Platform**: Builds and runs on Linux, macOS, and Windows.

## Installation

### From Source

You need Go installed on your machine.

1. **Clone the repository**:

    ```bash
    git clone https://github.com/dbehnke/contactfilter.git
    cd contactfilter
    ```

2. **Build the project**:

    ```bash
    go build -o contactfilter
    ```

## Usage

### Arguments

```bash
./contactfilter [flags] [input-csv] [filter-file] [output-csv]
```

* `--merge`: Activate merge mode (Brandmeister + RadioID).
* `--brandmeister-csv`: Path to Brandmeister "Last Heard" CSV (required for merge mode).
* `--radioid-csv`: Path to RadioID.net `user.csv`. If omitted in merge mode, it is downloaded automatically.
* `--priority-country`: The country to prioritize (default: "United States").
* `--limit`: Maximum number of contacts (default: 50000).

### Examples

#### 1. Standard Filtering (Baofeng CSV Input)

Filter an existing Baofeng-format CSV:

```bash
./contactfilter \
  --input-csv input.csv \
  --filter-file countries.txt \
  --output-csv filtered.csv
```

#### 2. Merge Mode (Brandmeister + RadioID)

Create a list of active users from your Brandmeister export:

```bash
./contactfilter \
  --merge \
  --brandmeister-csv bm_contacts.csv \
  --filter-file countries.txt \
  --output-csv active_contacts.csv
```

## How to Export Brandmeister Contacts

To get the "Last Heard" list for merge mode:

1. **Log in** to [Brandmeister Network](https://brandmeister.network/).
2. Go to the **Contacts Export** page: [https://brandmeister.network/?page=contactsexport](https://brandmeister.network/?page=contactsexport).
3. **Configure the Export**:
    * **Talkgroups**: You **must** enter a list of talkgroups.
    * *Example (Michigan Starter)*: `3126,31261,31262,31266,313136,3200449`
    * **Format**: Select **CSV**.
4. Click **Run** (or Download) to save the file.
5. Use this file as the `--brandmeister-csv` input (e.g., `bm_contacts.csv`).
