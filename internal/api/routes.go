package api

import (
	"ergo-tools-go/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	TenantID  string     `json:"tenant_id"`
	SessionID string     `json:"session_id"`
	Message   string     `json:"message"`
	Lead      *core.Lead `json:"lead,omitempty"`
}

func RegisterRoutes(r *gin.Engine) {
	cm := core.NewConversationManager()
	svc := core.NewService(cm)

	r.POST("/chat", func(c *gin.Context) {
		var req ChatRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Lead != nil {
			cm.SetLead(req.SessionID, *req.Lead)
		}

		reply, lead, err := svc.HandleTurn(c, req.TenantID, req.SessionID, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reply": reply, "lead": lead})
	})

	r.GET("/ws", func(c *gin.Context) {
		core.ServeWS(c.Writer, c.Request, cm)
	})
}
