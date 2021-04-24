package vern

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ess/hype"
)

type Client struct {
	driver  *hype.Driver
	headers []*hype.Header
}

func New(url, accessKey, secretKey string) (*Client, error) {
	driver, err := hype.New(url)
	if err != nil {
		return nil, err
	}

	headers := []*hype.Header{
		hype.NewHeader("Accept", "application/json"),
		hype.NewHeader("Content-Type", "application/json"),
		hype.NewHeader("Authorization", accessKey+":"+secretKey),
	}

	return &Client{driver: driver, headers: headers}, nil
}

func (client *Client) Analyze(text string) (*Analysis, error) {
	input := &payload{Text: text}

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	response := client.driver.
		Post("analyze", nil, data).
		WithHeaderSet(client.headers...).
		Response()

	if response.Okay() {
		analysis := &Analysis{}

		jErr := json.Unmarshal(response.Data(), analysis)
		if jErr != nil {
			return nil, jErr
		}

		return analysis, nil
	}

	return nil, response.Error()
}

type Score struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (score *Score) String() string {
	return fmt.Sprintf("%s: %f", score.Name, score.Value)
}

type Analysis struct {
	Text   string   `json:"text"`
	Scores []*Score `json:"scores"`
}

func (analysis *Analysis) String() string {
	scoreblock := "No scores"

	if len(analysis.Scores) > 0 {
		scores := make([]string, 0)

		for _, score := range analysis.Scores {
			scores = append(scores, score.String())
		}

		scoreblock = strings.Join(scores, ", ")
	}

	return fmt.Sprintf("%s (%s)", analysis.Text, scoreblock)
}

type payload struct {
	Text string `json:"text"`
}
