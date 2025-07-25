package domain

import (
	"encoding/json"
	"time"
)

// User JSON
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	return json.Unmarshal(data, &aux)
}

// Transaction JSON
func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(t),
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	t.CreatedAt = parsed
	return nil
}

// Balance JSON
func (b *Balance) MarshalJSON() ([]byte, error) {
	type Alias Balance
	return json.Marshal(&struct {
		LastUpdatedAt string `json:"last_updated_at"`
		*Alias
	}{
		LastUpdatedAt: b.LastUpdatedAt.Format(time.RFC3339),
		Alias:         (*Alias)(b),
	})
}

func (b *Balance) UnmarshalJSON(data []byte) error {
	type Alias Balance
	aux := &struct {
		LastUpdatedAt string `json:"last_updated_at"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, aux.LastUpdatedAt)
	if err != nil {
		return err
	}
	b.LastUpdatedAt = parsed
	return nil
}
