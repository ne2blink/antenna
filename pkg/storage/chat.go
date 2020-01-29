package storage

import "encoding/json"

// Chat is
type Chat struct {
	ID             string
	SubscribedApps []App
}

// FromJSON is
func (c *Chat) FromJSON(jsonBytes []byte) {
	json.Unmarshal(jsonBytes, c)
}

// ToJSON is
func (c Chat) ToJSON() ([]byte, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
