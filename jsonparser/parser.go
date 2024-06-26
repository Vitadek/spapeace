package jsonparser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Rule struct {
	GroupID   string `json:"groupid"`
	Status    string `json:"status"`
	FixText   string `json:"fix_text"`
	CheckText string `json:"check_text"`
}

type JSONData struct {
	Stig []struct {
		Rule []Rule `json:"rule"`
	} `json:"stig"`
}

func ParseJSONFile(filePath string) (*JSONData, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read file: " + err.Error())
	}

	var data JSONData
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, errors.New("failed to parse JSON: " + err.Error())
	}

	return &data, nil
}

func SaveJSONFile(filePath string, data *JSONData) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.New("failed to marshal JSON: " + err.Error())
	}

	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return errors.New("failed to write file: " + err.Error())
	}

	return nil
}
