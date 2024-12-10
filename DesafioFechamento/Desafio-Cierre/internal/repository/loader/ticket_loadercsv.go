package loader

import (
	"app/internal"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func NewLoaderTicketCSV(filePath string) LoaderTicketCSV {
	return LoaderTicketCSV{
		filePath: filePath,
	}
}

type LoaderTicketCSV struct {
	filePath string
}

func (l *LoaderTicketCSV) Load() (map[int]*internal.TicketAttributes, error) {
	f, err := os.Open(l.filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	tickets := make(map[int]*internal.TicketAttributes)
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading record: %v", err)
		}

		if len(record) < 6 {
			return nil, fmt.Errorf("insufficient fields in record: %v", record)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %v", record[0])
		}

		price, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price format: %v", record[5])
		}

		ticket := &internal.TicketAttributes{
			Name:    record[1],
			Email:   record[2],
			Country: record[3],
			Hour:    record[4],
			Price:   price,
		}

		tickets[id] = ticket
	}

	return tickets, nil
}
