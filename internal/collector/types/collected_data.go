package types

import (
	"encoding/json"
	"fmt"
)

type CollectedData struct {
	YandexData   *YandexData
	MagneticData *MagneticData
	EventsData   *EventsData
}

func (c *CollectedData) Clone() (*CollectedData, error) {
	copy, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal collected data: %w", err)
	}

	var collectedData CollectedData
	err = json.Unmarshal(copy, &collectedData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal collected data: %w", err)
	}

	return &collectedData, nil
}
