package service

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"merchant-bank-api/models"
)

type HistoryService interface {
	LogHistory(customerID string, action string) error
}

type historyService struct{}

func (s *historyService) LogHistory(customerID string, action string) error {
	history := s.createHistoryEntry(customerID, action)

	histories, err := s.readHistoriesFromFile("database/history.json")
	if err != nil {
		log.Printf("Error reading histories: %v", err)
		return nil
	}

	histories = append(histories, history)

	if err := s.writeHistoriesToFile("database/history.json", histories); err != nil {
		log.Printf("Error writing histories: %v", err)
	}

	return nil
}

// createHistoryEntry creates a new history entry.
func (s *historyService) createHistoryEntry(customerID, action string) models.History {
	return models.History{
		CustomerID: customerID,
		Action:     action,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// readHistoriesFromFile reads history entries from a file.
func (s *historyService) readHistoriesFromFile(filePath string) ([]models.History, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var histories []models.History
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&histories); err != nil && err != io.EOF {
		return nil, err
	}

	return histories, nil
}

// writeHistoriesToFile writes history entries to a file.
func (s *historyService) writeHistoriesToFile(filePath string, histories []models.History) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(histories); err != nil {
		return err
	}

	return nil
}

func NewHistoryService() HistoryService {
	return &historyService{}
}
