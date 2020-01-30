package storage

import "encoding/json"

// Chat is Telegram chat
type Chat struct {
	ID               string
	SubscribedAppIDs []string
}

// FromJSON is decoding json to Chat
func (c *Chat) FromJSON(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, c)
}

// ToJSON is encoding Chat to json
func (c Chat) ToJSON() ([]byte, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
