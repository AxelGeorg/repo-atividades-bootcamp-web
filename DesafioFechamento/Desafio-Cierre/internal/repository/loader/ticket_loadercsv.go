package loader

import (
	"app/internal"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// NewLoaderTicketCSV creates a new ticket loader from a CSV file
func NewLoaderTicketCSV(filePath string) LoaderTicketCSV {
	return LoaderTicketCSV{
		filePath: filePath,
	}
}

// LoaderTicketCSV represents a ticket loader from a CSV file
type LoaderTicketCSV struct {
	filePath string
}

// Load loads the tickets from the CSV file
func (l *LoaderTicketCSV) Load() (map[int]*internal.TicketAttributes, error) {
	// Open the file
	f, err := os.Open(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Read the file
	r := csv.NewReader(f)

	// Read the records
	tickets := make(map[int]*internal.TicketAttributes)
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		// Verifique se há campos suficientes
		if len(record) < 6 {
			return nil, fmt.Errorf("insufficient fields in record: %v", record)
		}

		// Serialize the record
		id, err := strconv.Atoi(record[0]) // Converte id de string para int
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %v", record[0])
		}

		price, err := strconv.ParseFloat(record[5], 64) // Converte price de string para float64
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %v", record[5])
		}

		ticket := &internal.TicketAttributes{ // Use & para criar um ponteiro
			Name:    record[1],
			Email:   record[2],
			Country: record[3],
			Hour:    record[4],
			Price:   price, // Agora price é o float64 convertido
		}

		// Add the ticket to the map
		tickets[id] = ticket
	}

	return tickets, nil
}
