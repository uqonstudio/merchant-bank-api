package service

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"merchant-bank-api/models"
)

// HistoryService defines the interface for logging customer history actions.
type HistoryService interface {
	// LogHistory logs a customer's action by creating a history entry and saving it to a file.
	LogHistory(customerID string, action string) error
}

// historyService is a concrete implementation of the HistoryService interface.
type historyService struct{}

// LogHistory logs a customer's action by creating a history entry and appending it to the history file.
func (s *historyService) LogHistory(customerID string, action string) error {
	// Create a new history entry for the given customer ID and action.
	history := s.createHistoryEntry(customerID, action)

	// Read existing history entries from the file.
	histories, err := s.readHistoriesFromFile("database/history.json")
	if err != nil {
		log.Printf("Error reading histories: %v", err)
		return nil
	}

	// Append the new history entry to the list of histories.
	histories = append(histories, history)

	// Write the updated list of histories back to the file.
	if err := s.writeHistoriesToFile("database/history.json", histories); err != nil {
		log.Printf("Error writing histories: %v", err)
	}

	return nil
}

// createHistoryEntry creates a new history entry with the current timestamp.
func (s *historyService) createHistoryEntry(customerID, action string) models.History {
	return models.History{
		CustomerID: customerID,
		Action:     action,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// readHistoriesFromFile reads history entries from a specified JSON file.
func (s *historyService) readHistoriesFromFile(filePath string) ([]models.History, error) {
	// Open the file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data into a slice of History objects.
	var histories []models.History
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&histories); err != nil && err != io.EOF {
		return nil, err
	}

	return histories, nil
}

// writeHistoriesToFile writes a slice of history entries to a specified JSON file.
func (s *historyService) writeHistoriesToFile(filePath string, histories []models.History) error {
	// Open the file for writing, creating it if it doesn't exist, and truncating it if it does.
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the slice of History objects into JSON and write it to the file.
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(histories); err != nil {
		return err
	}

	return nil
}

// NewHistoryService creates a new instance of historyService and returns it as a HistoryService.
func NewHistoryService() HistoryService {
	return &historyService{}
}
