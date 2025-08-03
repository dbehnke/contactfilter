use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::env;
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

/// The main function returns a Result, which is a common and idiomatic
/// way to handle errors in Rust applications.
fn main() -> Result<(), Box<dyn Error>> {
    // --- 1. Argument Parsing ---
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 {
        eprintln!("Usage: {} <input_csv> <filter_file> <output_csv>", args[0]);
        std::process::exit(1);
    }

    let input_path = &args[1];
    let filter_path = &args[2];
    let output_path = &args[3];

    println!("Input file: {}", input_path);
    println!("Filter file: {}", filter_path);
    println!("Output file: {}", output_path);

    // --- 2. Load Filtered Countries ---
    // We read the filter file and load the countries into a HashSet.
    // A HashSet provides very fast lookups (O(1) on average), which is
    // much more efficient than searching through a list for every row.
    let filter_content = fs::read_to_string(filter_path)?;
    let countries_to_keep: HashSet<String> = filter_content
        .lines()
        .map(String::from)
        .collect();

    println!("Filtering for {} countries.", countries_to_keep.len());

    // --- 3. Process the CSV ---
    // Create a CSV reader and writer.
    let mut rdr = csv::Reader::from_path(input_path)?;
    let mut wtr = csv::Writer::from_path(output_path)?;

    let mut records_read = 0;
    let mut records_written = 0;

    // We write the headers to the output file first.
    wtr.write_record(rdr.headers()?)?;

    // Iterate over each record in the input CSV.
    // `deserialize` automatically converts each row into our `Contact` struct.
    for result in rdr.deserialize() {
        let record: Contact = result?;
        records_read += 1;

        // Check if the record's country is in our filter set.
        if countries_to_keep.contains(&record.country) {
            // If it is, serialize it back into a CSV row and write it.
            wtr.serialize(&record)?;
            records_written += 1;
        }
    }

    // It's good practice to flush the writer to ensure all data is written to the file.
    wtr.flush()?;

    println!("\nProcessing complete!");
    println!("Read {} records.", records_read);
    println!("Wrote {} records.", records_written);

    Ok(())
}