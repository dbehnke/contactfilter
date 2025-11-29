package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Contact represents a single contact record in the Baofeng CSV format.
type Contact struct {
	No        int
	ID        int
	Repeater  string
	Name      string
	City      string
	Province  string
	Country   string
	Remark    string
	CallType  string
	AlertCall string
}

// RadioIDRecord represents a record from the RadioID.net user.csv
type RadioIDRecord struct {
	ID       int
	Callsign string
	Fname    string
	Surname  string
	City     string
	State    string
	Country  string
}

func main() {
	// Flags
	inputCSV := flag.String("input-csv", "", "Path to the input CSV file. In --merge mode, this is the Brandmeister Last Heard CSV.")
	filterFile := flag.String("filter-file", "", "Path to the country filter file (one country per line).")
	outputCSV := flag.String("output-csv", "", "Path for the new, filtered output CSV file.")
	priorityCountry := flag.String("priority-country", "United States", "The country to prioritize.")
	limit := flag.Int("limit", 50000, "The maximum number of contacts to include.")
	merge := flag.Bool("merge", false, "Activate merge mode.")
	brandmeisterCSV := flag.String("brandmeister-csv", "", "Path to Brandmeister 'Last Heard' CSV export (optional alias for input-csv in merge mode).")
	radioidCSV := flag.String("radioid-csv", "", "Path to RadioID.net user.csv. If not provided in merge mode, it will be downloaded.")

	radioType := flag.String("radio", "baofeng-dm32uv", "The radio output format: baofeng-dm32uv, anytone, opengd77")

	flag.Parse()

	// Handle positional arguments for backward compatibility or ease of use
	args := flag.Args()
	if *inputCSV == "" && len(args) > 0 {
		*inputCSV = args[0]
	}
	if *filterFile == "" && len(args) > 1 {
		*filterFile = args[1]
	}
	if *outputCSV == "" && len(args) > 2 {
		*outputCSV = args[2]
	}

	// If brandmeister-csv is set, use it as input-csv
	if *brandmeisterCSV != "" {
		*inputCSV = *brandmeisterCSV
		*merge = true
	}

	if *inputCSV == "" && !*merge {
		fmt.Println("Error: Input CSV is required.")
		flag.Usage()
		os.Exit(1)
	}
	if *filterFile == "" {
		fmt.Println("Error: Filter file is required.")
		flag.Usage()
		os.Exit(1)
	}
	if *outputCSV == "" {
		fmt.Println("Error: Output CSV is required.")
		flag.Usage()
		os.Exit(1)
	}

	// Load Filtered Countries
	countriesToKeep, err := loadCountries(*filterFile)
	if err != nil {
		fmt.Printf("Error loading filter file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Filtering for %d countries.\n", len(countriesToKeep))

	var contacts []Contact

	if *merge {
		fmt.Println("\n--- Merge Mode Activated ---")
		bmPath := *inputCSV
		if bmPath == "" {
			fmt.Println("Error: Brandmeister CSV path required in merge mode.")
			os.Exit(1)
		}

		// 1. Read Brandmeister CSV to get active IDs
		fmt.Printf("Reading Brandmeister Last Heard data from: %s\n", bmPath)
		activeIDs, err := readBrandmeisterIDs(bmPath)
		if err != nil {
			fmt.Printf("Error reading Brandmeister CSV: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Found %d unique active IDs in Brandmeister data.\n", len(activeIDs))

		// 2. Get RadioID data
		var ridReader io.Reader
		if *radioidCSV != "" {
			fmt.Printf("Reading RadioID data from: %s\n", *radioidCSV)
			f, err := os.Open(*radioidCSV)
			if err != nil {
				fmt.Printf("Error opening RadioID CSV: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			ridReader = f
		} else {
			fmt.Println("Downloading latest RadioID database...")
			resp, err := http.Get("https://www.radioid.net/static/user.csv")
			if err != nil {
				fmt.Printf("Error downloading RadioID database: %v\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			ridReader = resp.Body
		}

		// 3. Process RadioID data
		fmt.Println("Processing RadioID database...")
		contacts, err = processRadioID(ridReader, activeIDs)
		if err != nil {
			fmt.Printf("Error processing RadioID data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Matched %d contacts from RadioID database.\n", len(contacts))

	} else {
		// Standard Mode
		fmt.Printf("\n--- Standard Mode: Reading from %s ---\n", *inputCSV)
		var err error
		contacts, err = readBaofengCSV(*inputCSV)
		if err != nil {
			fmt.Printf("Error reading input CSV: %v\n", err)
			os.Exit(1)
		}
	}

	// Filter and Prioritize
	var usContacts []Contact
	var otherContacts []Contact

	// Check if we need to filter
	shouldFilter := len(contacts) > *limit
	if !shouldFilter {
		fmt.Printf("Total contacts (%d) is within limit (%d). Skipping country filter.\n", len(contacts), *limit)
		// Keep all, just prioritize
		for _, c := range contacts {
			if c.Country == *priorityCountry {
				usContacts = append(usContacts, c)
			} else {
				otherContacts = append(otherContacts, c)
			}
		}
	} else {
		fmt.Printf("Total contacts (%d) exceeds limit (%d). Applying country filter.\n", len(contacts), *limit)
		for _, c := range contacts {
			if _, ok := countriesToKeep[c.Country]; ok {
				if c.Country == *priorityCountry {
					usContacts = append(usContacts, c)
				} else {
					otherContacts = append(otherContacts, c)
				}
			}
		}
	}

	fmt.Printf("Found %d contacts from priority country (%s).\n", len(usContacts), *priorityCountry)
	fmt.Printf("Found %d contacts from other countries.\n", len(otherContacts))

	finalContacts := append(usContacts, otherContacts...)

	if len(finalContacts) > *limit {
		fmt.Printf("Truncating list from %d to %d records.\n", len(finalContacts), *limit)
		finalContacts = finalContacts[:*limit]
	}

	// Write Output
	fmt.Printf("\nWriting to %s (Format: %s)...\n", *outputCSV, *radioType)
	if err := writeCSV(*outputCSV, finalContacts, *radioType); err != nil {
		fmt.Printf("Error writing output CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func loadCountries(path string) (map[string]bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	countries := make(map[string]bool)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			countries[trimmed] = true
		}
	}
	return countries, nil
}

func readBrandmeisterIDs(path string) (map[int]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Allow variable number of fields
	r.FieldsPerRecord = -1

	// Read header
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}

	// Find ID column index
	idIdx := -1
	for i, h := range headers {
		h = strings.TrimSpace(h)
		if strings.EqualFold(h, "Sending ID") || strings.EqualFold(h, "Radio ID") || strings.EqualFold(h, "ID") {
			idIdx = i
			break
		}
	}

	// If header not found, try index 0 as fallback or fail?
	// Let's assume index 0 if not found, but warn? Or just fail.
	// For dummy data "Sending ID" is at 0.
	if idIdx == -1 {
		// Fallback to 0
		idIdx = 0
	}

	ids := make(map[int]bool)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip malformed lines
		}
		if len(record) <= idIdx {
			continue
		}

		idStr := strings.TrimSpace(record[idIdx])
		if id, err := strconv.Atoi(idStr); err == nil {
			ids[id] = true
		}
	}
	return ids, nil
}

func processRadioID(r io.Reader, activeIDs map[int]bool) ([]Contact, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// Map headers to indices
	idxMap := make(map[string]int)
	for i, h := range headers {
		idxMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	var contacts []Contact
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// Helper to get value safely
		getVal := func(key string) string {
			if idx, ok := idxMap[key]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		idStr := getVal("radio_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}

		if activeIDs[id] {
			fname := getVal("first_name")
			surname := getVal("last_name")
			name := strings.TrimSpace(fname + " " + surname)

			c := Contact{
				ID:        id,
				Repeater:  getVal("callsign"), // Use callsign as repeater/name
				Name:      name,
				City:      getVal("city"),
				Province:  getVal("state"),
				Country:   getVal("country"),
				CallType:  "Private Call",
				AlertCall: "None",
			}
			contacts = append(contacts, c)
		}
	}
	return contacts, nil
}

func readBaofengCSV(path string) ([]Contact, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	// Skip header
	_, err = r.Read()
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	// Assuming standard Baofeng column order:
	// No.,ID,Repeater,Name,City,Province,Country,Remark,Type,Alert Call
	// 0   1  2        3    4    5        6       7      8    9

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) < 10 {
			continue
		}

		id, _ := strconv.Atoi(record[1])

		c := Contact{
			ID:        id,
			Repeater:  record[2],
			Name:      record[3],
			City:      record[4],
			Province:  record[5],
			Country:   record[6],
			Remark:    record[7],
			CallType:  record[8],
			AlertCall: record[9],
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func writeCSV(path string, contacts []Contact, radioType string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if radioType == "anytone" {
		// Manual writing for Anytone to ensure all fields are quoted and use CRLF
		w := bufio.NewWriter(f)
		defer w.Flush()

		// Helper to quote and join fields
		quoteAndJoin := func(fields []string) string {
			quoted := make([]string, len(fields))
			for i, f := range fields {
				quoted[i] = `"` + strings.ReplaceAll(f, `"`, `""`) + `"`
			}
			return strings.Join(quoted, ",")
		}

		// Header
		header := []string{"No.", "Radio ID", "Callsign", "Name", "City", "State", "Country", "Remarks", "Call Type", "Call Alert"}
		if _, err := w.WriteString(quoteAndJoin(header) + "\r\n"); err != nil {
			return err
		}

		for i, c := range contacts {
			row := []string{
				strconv.Itoa(i + 1),
				strconv.Itoa(c.ID),
				c.Repeater, // Using Repeater field for Callsign
				c.Name,
				c.City,
				c.Province,
				c.Country,
				c.Remark,
				c.CallType,
				c.AlertCall,
			}
			if _, err := w.WriteString(quoteAndJoin(row) + "\r\n"); err != nil {
				return err
			}
		}
		return nil
	}

	w := csv.NewWriter(f)
	if radioType == "db25d" {
		w.UseCRLF = true
	}
	defer w.Flush()

	var header []string
	switch strings.ToLower(radioType) {
	case "anytone":
		header = []string{"No.", "Radio ID", "Callsign", "Name", "City", "State", "Country", "Remarks", "Call Type", "Call Alert"}
	case "opengd77":
		header = []string{"Radio ID", "Callsign", "Name", "Nickname", "City", "State", "Country", "Remarks"}
	case "db25d":
		// Sample format: Radio-ID,Callsign,First-Name,City,State/Prov,Country;
		// Note the trailing semicolon in the header and rows
		header = []string{"Radio-ID", "Callsign", "First-Name", "City", "State/Prov", "Country;"}
	default: // baofeng-dm32uv
		header = []string{"No.", "ID", "Repeater", "Name", "City", "Province", "Country", "Remark", "Type", "Alert Call"}
	}

	if err := w.Write(header); err != nil {
		return err
	}

	for i, c := range contacts {
		var row []string
		switch strings.ToLower(radioType) {
		case "anytone":
			row = []string{
				strconv.Itoa(i + 1),
				strconv.Itoa(c.ID),
				c.Repeater, // Using Repeater field for Callsign as per read logic
				c.Name,
				c.City,
				c.Province,
				c.Country,
				c.Remark,
				c.CallType,
				c.AlertCall,
			}
		case "opengd77":
			// Nickname logic: First word of Name
			nickname := strings.Split(c.Name, " ")[0]
			row = []string{
				strconv.Itoa(c.ID),
				c.Repeater, // Using Repeater field for Callsign
				c.Name,
				nickname,
				c.City,
				c.Province,
				c.Country,
				c.Remark,
			}
		case "db25d":
			// First-Name logic: First word of Name
			firstName := strings.Split(c.Name, " ")[0]
			row = []string{
				strconv.Itoa(c.ID),
				c.Repeater, // Using Repeater field for Callsign
				firstName,
				c.City,
				c.Province,
				c.Country + ";", // Append trailing semicolon
			}
		default: // baofeng-dm32uv
			row = []string{
				strconv.Itoa(i + 1),
				strconv.Itoa(c.ID),
				c.Repeater,
				c.Name,
				c.City,
				c.Province,
				c.Country,
				c.Remark,
				c.CallType,
				c.AlertCall,
			}
		}

		if err := w.Write(row); err != nil {
			return err
		}
	}
	return nil
}
