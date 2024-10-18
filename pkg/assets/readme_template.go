package assets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/harshalranjhani/genie/internal/middleware"
)

var defaultTemplate = `# {{.projectName}}

## Description
{{.description}}

## Installation
{{.installation}}

## Usage
{{.usage}}

## Contributing
{{.contributing}}

## License
{{.license}}
`

var minimalTemplate = `# {{.projectName}}

{{.description}}

## Quick Start
{{.installation}}
{{.usage}}

## License
{{.license}}
`

var detailedTemplate = `# {{.projectName}}

## Description
{{.description}}

## Features
{{.features}}

## Prerequisites
{{.prerequisites}}

## Installation
{{.installation}}

## Configuration
{{.configuration}}

## Usage
{{.usage}}

## Testing
{{.testing}}

## Deployment
{{.deployment}}

## Contributing
{{.contributing}}

## License
{{.license}}

## Acknowledgements
{{.acknowledgements}}
`

func GetReadmeTemplate(templateName string) string {
	status, err := middleware.LoadStatus()
	if err != nil {
		return defaultTemplate
	}
	token := status.Token
	switch templateName {
	case "minimal":
		return minimalTemplate
	case "detailed":
		return detailedTemplate
	case "animated", "interactive":
		return getProReadmeTemplate(templateName, token)
	default:
		return defaultTemplate
	}
}

func getProReadmeTemplate(templateName string, token string) string {
	url := "https://api.harshalranjhani.in/genie/pro-readme-template"
	payload := strings.NewReader(fmt.Sprintf(`{"token":"%s","templateName":"%s"}`, token, templateName))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching pro template:", err)
		return defaultTemplate
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Error fetching pro template. Status:", res.Status)
		return defaultTemplate
	}

	var response struct {
		Template string `json:"template"`
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return defaultTemplate
	}

	return response.Template
}
