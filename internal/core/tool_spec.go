package core

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

func OpenAIToolsSpec() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "capture_contact_info",
				Description: "Persist structured fields gathered from the caller.",
				Parameters: rawParams(map[string]any{
					"type": "object",
					"properties": map[string]any{
						"first_name": map[string]any{"type": "string"},
						"last_name":  map[string]any{"type": "string"},
						"email":      map[string]any{"type": "string"},
						"phone":      map[string]any{"type": "string"},
						"language":   map[string]any{"type": "string"},
						"case_type":  map[string]any{"type": "string"},
						"notes":      map[string]any{"type": "string"},
					},
					"additionalProperties": false,
				}),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_missing_fields",
				Description: "Return which fields are still required for this tenant/case.",
				Parameters: rawParams(map[string]any{
					"type": "object",
					"properties": map[string]any{
						"case_type": map[string]any{"type": "string"},
					},
					"additionalProperties": false,
				}),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "export_to_crm",
				Description: "Export current lead to tenant CRM and return an export id.",
				Parameters: rawParams(map[string]any{
					"type": "object",
					"properties": map[string]any{
						"provider": map[string]any{"type": "string", "default": "lawmatics"},
					},
					"additionalProperties": false,
				}),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "calendly",
				Description: "Schedule a consultation using Calendly API and return a link/slot.",
				Parameters: rawParams(map[string]any{
					"type": "object",
					"properties": map[string]any{
						"desired_date": map[string]any{"type": "string"},
						"desired_time": map[string]any{"type": "string"},
					},
					"additionalProperties": false,
				}),
			},
		},
	}
}

// helper: encode schema to json.RawMessage
func rawParams(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}
