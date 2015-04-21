package enigma

import (
	"encoding/json"
	"fmt"
	"os"
)

type ComponentInfo struct {
	Type string `json:"type"`
	In   string `json:"in"`
	Out  string `json:"out"`
	Offset int `json:"offset"`
}

func FromJSON(data []byte) (*Component, error) {
	var info []ComponentInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("FromJSON: %v", err)
	}
	return Build(info)
}

func FromJSONFile(filename string) (*Component, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("FromJSONFile: %v", err)
	}
	defer file.Close()

	var info []ComponentInfo
	err = json.NewDecoder(file).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("FromJSONFile: %v", err)
	}

	return Build(info)
}

func Build(info []ComponentInfo) (*Component, error) {
	if len(info) <= 0 {
		return nil, fmt.Errorf("Build: empty component info")
	}

	comps := make([]*Component, len(info))
	for i, v := range info {
		var c *Component
		switch v.Type {
		case "rotor":
			c = NewComponent(Rotor)
		case "reflector":
			c = NewComponent(Reflector)
		case "plugboard":
			c = NewComponent(Plugboard)
		default:
			return nil, fmt.Errorf("Build: unknown component type")
		}
		c.SetCharacterMap(v.In, v.Out)
		c.Step(v.Offset)
		comps[i] = c
	}

	Connect(comps...)

	return comps[0], nil
}
