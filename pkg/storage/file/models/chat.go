package models

import "encoding/json"

// Chat is Telegram chat
type Chat struct {
	ID               int64    `json:"id,omitempty"`
	SubscribedAppIDs []string `json:"subscribed_app_ids,omitempty"`
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
