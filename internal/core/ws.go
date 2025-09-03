package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	openai "github.com/sashabaranov/go-openai"
)

type wsMsg struct {
	TenantID  string `json:"tenant_id"`
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type wsReply struct {
	Reply string `json:"reply"`
	Lead  Lead   `json:"lead"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func ServeWS(w http.ResponseWriter, r *http.Request, cm *ConversationManager) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer c.Close()

	svc := NewService(cm)
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			return
		}
		var m wsMsg
		if err := json.Unmarshal(data, &m); err != nil {
			return
		}

		msgs := []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: SystemPrompt},
			{Role: openai.ChatMessageRoleSystem, Content: DevPrompt, Name: "developer"},
			{Role: openai.ChatMessageRoleUser, Content: m.Message},
		}
		assistant, err := svc.ai.RunWithTools(context.Background(), msgs, OpenAIToolsSpec(), svc.handleTool(m.TenantID, m.SessionID))
		if err != nil {
			return
		}
		reply := wsReply{Reply: assistant.Content, Lead: cm.Get(m.SessionID).Lead}
		b, _ := json.Marshal(reply)
		_ = c.WriteMessage(websocket.TextMessage, b)
	}
}
