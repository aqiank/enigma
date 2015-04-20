package enigma

import (
	"encoding/json"
	"fmt"
	"os"
)

type buildInfo struct {
	Components []componentInfo `json:"components"`
}

type componentInfo struct {
	Type string `json:"type"`
	In string `json:"in"`
	Out string `json:"out"`
}

func FromJSON(data []byte) (*Component, error) {
	var info buildInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("FromJSON: %v", err)
	}
	return build(info)
}

func FromJSONFile(filename string) (*Component, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("FromJSONFile: %v", err)
	}
	defer file.Close()

	var info buildInfo
	err = json.NewDecoder(file).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("FromJSONFile: %v", err)
	}

	return build(info)
}

func build(info buildInfo) (*Component, error) {
	var first, prev *Component
	for _, v := range info.Components {
		var next *Component
		switch v.Type {
		case "rotor":
			next = NewComponent(Rotor)
		case "reflector":
			next = NewComponent(Reflector)
		case "plugboard":
			next = NewComponent(Plugboard)
		default:
			return nil, fmt.Errorf("build: unknown component type")
		}
		next.Set(v.In, v.Out)
		if first == nil {
			first = next
			prev = next
		} else {
			prev = prev.Connect(next)
		}
	}
	return first, nil
}
