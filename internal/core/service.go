package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

type Service struct {
	cm *ConversationManager
	ai *OpenAIClient
	rt *ToolRuntime
}

func NewService(cm *ConversationManager) *Service {
	return &Service{cm: cm, ai: NewOpenAI(), rt: NewToolRuntime(cm)}
}

func (s *Service) handleTool(tenantID, sessionID string) ToolHandler {
	return func(name string, args map[string]any) (any, error) {
		s.cm.Get(sessionID).TenantID = tenantID
		switch name {
		case "capture_contact_info":
			return s.rt.CaptureContactInfo(sessionID, args)
		case "get_missing_fields":
			return s.rt.GetMissingFields(sessionID, tenantID, args)
		case "export_to_crm":
			provider, _ := args["provider"].(string)
			if provider == "" {
				provider = "lawmatics"
			}
			return s.rt.ExportToCRM(sessionID, provider)
		case "calendly":
			return s.rt.Calendly(sessionID, args)
		default:
			return nil, fmt.Errorf("unknown tool: %s", name)
		}
	}
}

func (s *Service) HandleTurn(ctx *gin.Context, tenantID, sessionID, userMsg string) (string, Lead, error) {
	msgs := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: SystemPrompt},
		{Role: openai.ChatMessageRoleSystem, Content: DevPrompt, Name: "developer"},
		{Role: openai.ChatMessageRoleUser, Content: userMsg},
	}
	assistant, err := s.ai.RunWithTools(ctx, msgs, OpenAIToolsSpec(), s.handleTool(tenantID, sessionID))
	if err != nil {
		return "", Lead{}, err
	}
	lead := s.cm.Get(sessionID).Lead
	return assistant.Content, lead, nil
}
