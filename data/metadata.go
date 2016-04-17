// Contains the data model of metadata needed for any page.
package data

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	dateString := t.Format("2006-01-02 15:04:05+00:00")
	return json.Marshal(dateString)
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	var dateString string
	if err := json.Unmarshal(data, &dateString); err != nil {
		return err
	}
	timeValue, err := time.Parse("2006-01-02 15:04:05+00:00", dateString)
	*t = JSONTime{timeValue}
	return err
}

type Metadata struct {
	Title string
	Slug  string
	Date  JSONTime
	Url   string
	// The question is whether to just the content here or not and have it all in this pipeline...
	// I'd say no as it's easy to read it anyway and there is only one version in the repo that way.
	ContentPath string
}

func (metadata Metadata) Content() (string, error) {
	data, err := ioutil.ReadFile(metadata.ContentPath)
	return string(data), err
}
