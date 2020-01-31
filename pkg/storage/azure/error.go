package azure

import "encoding/json"

type errGroup []error

func (g *errGroup) Append(err error) {
	if err != nil {
		*g = append(*g, err)
	}
}

func (g errGroup) Error() string {
	switch len(g) {
	case 0:
		return ""
	case 1:
		return g[0].Error()
	}

	messages := make([]string, 0, len(g))
	for _, err := range g {
		messages = append(messages, err.Error())
	}
	message, _ := json.Marshal(messages)
	return string(message)
}

func (g errGroup) HasError() bool {
	return len(g) > 0
}

func (g errGroup) Simplify() error {
	if len(g) == 0 {
		return nil
	}
	return g
}
