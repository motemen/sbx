package sbapi

import "encoding/json"

type (
	_page    Page
	_project Project
)

type rawJSON interface {
	RawJSON() json.RawMessage
	setRawJSON(json.RawMessage)
}

func (p Page) RawJSON() json.RawMessage      { return p.rawJSON }
func (p *Page) setRawJSON(j json.RawMessage) { p.rawJSON = j }

func (p Project) RawJSON() json.RawMessage      { return p.rawJSON }
func (p *Project) setRawJSON(j json.RawMessage) { p.rawJSON = j }

func unmarshalJSON(data []byte, r rawJSON, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	r.setRawJSON(append([]byte{}, data...))

	return nil
}

func marshalJSON(r rawJSON, obj interface{}) ([]byte, error) {
	if j := r.RawJSON(); j != nil {
		return json.Marshal(j)
	}

	return json.Marshal(obj)
}

func (p *Page) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, p, (*_page)(p))
}

func (p Page) MarshalJSON() ([]byte, error) {
	return marshalJSON(&p, (_page)(p))
}

func (p *Project) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, p, (*_project)(p))
}

func (p Project) MarshalJSON() ([]byte, error) {
	return marshalJSON(&p, (_project)(p))
}
