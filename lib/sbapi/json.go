package sbapi

import "encoding/json"

type pageWithoutRaw Page

func (p *Page) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*pageWithoutRaw)(p))
	if err != nil {
		return err
	}

	p.RawMessage = make([]byte, len(data))
	copy(p.RawMessage, data)

	return nil
}

func (p Page) MarshalJSON() ([]byte, error) {
	if p.RawMessage != nil {
		return json.Marshal(p.RawMessage)
	}

	return json.Marshal(pageWithoutRaw(p))
}

type projectWithoutRaw Project

func (p *Project) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*projectWithoutRaw)(p))
	if err != nil {
		return err
	}

	p.RawMessage = make([]byte, len(data))
	copy(p.RawMessage, data)

	return nil
}

func (p Project) MarshalJSON() ([]byte, error) {
	if p.RawMessage != nil {
		return json.Marshal(p.RawMessage)
	}

	return json.Marshal(projectWithoutRaw(p))
}
