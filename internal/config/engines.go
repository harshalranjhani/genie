package config

import "strings"

const (
	GeminiEngine   = "Gemini"
	GPTEngine      = "GPT"
	DeepSeekEngine = "DeepSeek"
	OllamaEngine   = "Ollama"
)

// Engine represents an AI engine configuration
type Engine struct {
	Name         string
	Models       []string
	DefaultModel string
	Features     EngineFeatures
}

// EngineFeatures represents supported features for an engine
type EngineFeatures struct {
	SupportsImageGen      bool
	SupportsChat          bool
	SupportsSafeMode      bool
	SupportsReasoning     bool
	SupportsDocumentation bool
}

// EngineMap stores all available engines
var EngineMap = map[string]Engine{
	"GPT": {
		Name: "GPT",
		Models: []string{
			"gpt-4",
			"gpt-4-turbo-preview",
			"gpt-3.5-turbo",
			"gpt-4o",
			"gpt-4o-2024-11-20",
			"gpt-4o-2024-08-06",
			"gpt-4o-mini",
			"gpt-4o-mini-2024-07-18",
		},
		DefaultModel: "gpt-4",
		Features: EngineFeatures{
			SupportsImageGen:      true,
			SupportsChat:          true,
			SupportsSafeMode:      true,
			SupportsReasoning:     true,
			SupportsDocumentation: true,
		},
	},
	"Gemini": {
		Name: "Gemini",
		Models: []string{
			"gemini-1.5-pro",
			"gemini-1.5-flash",
			"gemini-1.5-flash-8b",
		},
		DefaultModel: "gemini-1.5-pro",
		Features: EngineFeatures{
			SupportsImageGen:      false,
			SupportsChat:          true,
			SupportsSafeMode:      true,
			SupportsReasoning:     false,
			SupportsDocumentation: false,
		},
	},
	"DeepSeek": {
		Name: "DeepSeek",
		Models: []string{
			"deepseek-chat",
			"deepseek-reasoner",
		},
		DefaultModel: "deepseek-chat",
		Features: EngineFeatures{
			SupportsImageGen:      false,
			SupportsChat:          true,
			SupportsSafeMode:      false,
			SupportsReasoning:     true,
			SupportsDocumentation: true,
		},
	},
	"Ollama": {
		Name: "Ollama",
		Models: []string{
			"llama3.2",
		},
		DefaultModel: "llama3.2",
		Features: EngineFeatures{
			SupportsImageGen:      false,
			SupportsChat:          true,
			SupportsSafeMode:      false,
			SupportsReasoning:     false,
			SupportsDocumentation: true,
		},
	},
}

// Helper functions
func GetEngine(name string) (Engine, bool) {
	engine, exists := EngineMap[name]
	return engine, exists
}

func CheckAndGetEngine(name string) (Engine, bool) {
	lookupName := strings.ToLower(name)

	for engineName, engine := range EngineMap {
		if strings.ToLower(engineName) == lookupName {
			return engine, true
		}
	}
	return Engine{}, false
}

func GetNextEngine(currentEngine string) string {
	switch currentEngine {
	case "Gemini":
		return "GPT"
	case "GPT":
		return "DeepSeek"
	case "DeepSeek":
		return "Ollama"
	case "Ollama":
		return "Gemini"
	default:
		return "Gemini"
	}
}

func IsValidModel(engineName string, modelName string) bool {
	engine, exists := EngineMap[engineName]
	if !exists {
		return false
	}

	for _, model := range engine.Models {
		if model == modelName {
			return true
		}
	}
	return false
}
