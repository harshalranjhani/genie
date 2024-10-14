package assets

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
	switch templateName {
	case "minimal":
		return minimalTemplate
	case "detailed":
		return detailedTemplate
	default:
		return defaultTemplate
	}
}
