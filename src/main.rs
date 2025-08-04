use clap::Parser;
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::error::Error;
use std::fs;

/// Represents a single contact record in the CSV file.
/// The `serde` attributes allow us to automatically map the CSV columns
/// to the fields of this struct, even when the column names are not
/// valid Rust identifiers (like "No." or "Alert Call").
#[derive(Debug, Deserialize, Serialize)]
struct Contact {
    #[serde(rename = "No.")]
    no: u32,
    #[serde(rename = "ID")]
    id: u32,
    #[serde(rename = "Repeater")]
    repeater: String,
    #[serde(rename = "Name")]
    name: String,
    #[serde(rename = "City")]
    city: String,
    #[serde(rename = "Province")]
    province: String,
    #[serde(rename = "Country")]
    country: String,
    #[serde(rename = "Remark")]
    remark: String,
    #[serde(rename = "Type")]
    call_type: String,
    #[serde(rename = "Alert Call")]
    alert_call: String,
}

/// A simple command-line utility to filter a CSV file of contacts
/// based on a list of countries, with prioritization and size limiting.
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Cli {
    /// Path to the input CSV file.
    input_csv: String,
    /// Path to the country filter file (one country per line).
    filter_file: String,
    /// Path for the new, filtered output CSV file.
    output_csv: String,

    /// The country to prioritize, ensuring its contacts appear first.
    #[arg(long, default_value = "United States")]
    priority_country: String,

    /// The maximum number of contacts to include in the final output.
    #[arg(long, default_value_t = 50_000)]
    limit: usize,
}

/// The main function returns a Result, which is a common and idiomatic
/// way to handle errors in Rust applications.
fn main() -> Result<(), Box<dyn Error>> {
    // --- 1. Argument Parsing ---
    let cli = Cli::parse();

    // --- 2. Load Filtered Countries ---
    // We read the filter file and load the countries into a HashSet.
    // A HashSet provides very fast lookups (O(1) on average), which is
    // much more efficient than searching through a list for every row.
    let filter_content = fs::read_to_string(&cli.filter_file)?;
    let countries_to_keep: HashSet<String> = filter_content
        .lines()
        .map(String::from)
        .collect();

    println!("Filtering for {} countries.", countries_to_keep.len());

    // --- 3. Phase 1: Read and Filter Contacts ---
    // We'll read all matching contacts into memory first. This allows us
    // to prioritize and truncate the list before writing to the output file.
    println!("\n--- Phase 1: Reading and filtering contacts from {}... ---", cli.input_csv);
    let mut rdr = csv::Reader::from_path(&cli.input_csv)?;
    let headers = rdr.headers()?.clone(); // Store headers for later use.

    let mut us_contacts: Vec<Contact> = Vec::new();
    let mut other_contacts: Vec<Contact> = Vec::new();
    let mut records_read = 0;

    for result in rdr.deserialize() {
        let record: Contact = result?;
        records_read += 1;

        // If the contact's country is in our list of countries to keep...
        if countries_to_keep.contains(&record.country) {
            // ...separate the US contacts from the others to enforce priority.
            if record.country == cli.priority_country {
                us_contacts.push(record);
            } else {
                other_contacts.push(record);
            }
        }
    }
    println!("Read {} records.", records_read);
    println!("Found {} contacts from the priority country ({}).", us_contacts.len(), &cli.priority_country);
    println!("Found {} contacts from other filtered countries.", other_contacts.len());

    // --- 4. Phase 2: Prioritize and Truncate ---
    println!("\n--- Phase 2: Prioritizing list and truncating to {} records... ---", cli.limit);

    // Combine the lists, with US contacts taking precedence at the start of the list.
    let mut final_contacts = us_contacts;
    final_contacts.append(&mut other_contacts);

    let total_filtered = final_contacts.len();
    println!("Total filtered contacts before truncation: {}", total_filtered);

    // If the combined list exceeds the limit, truncate it.
    if final_contacts.len() > cli.limit {
        final_contacts.truncate(cli.limit);
        println!("List truncated to the first {} records.", cli.limit);
    }

    // --- 5. Phase 3: Write to Output File ---
    println!("\n--- Phase 3: Writing final list to output file... ---");
    // We use a WriterBuilder to explicitly disable automatic header writing.
    // This is because we are manually writing the headers we captured from
    // the input file, ensuring they are preserved exactly.
    let mut wtr = csv::WriterBuilder::new()
        .has_headers(false)
        .from_path(&cli.output_csv)?;
    wtr.write_record(&headers)?; // Write the headers we saved earlier.

    let records_to_write_count = final_contacts.len();
    // Iterate through the final list, renumbering the 'No.' column starting from 1.
    // We use enumerate() to get a 0-based index, so we add 1 for the row number.
    for (i, mut record) in final_contacts.into_iter().enumerate() {
        record.no = (i + 1) as u32;
        wtr.serialize(record)?;
    }

    // It's good practice to flush the writer to ensure all data is written to the file.
    wtr.flush()?;

    println!("\nProcessing complete!");
    println!("Wrote {} records to {}.", records_to_write_count, &cli.output_csv);

    Ok(())
}
