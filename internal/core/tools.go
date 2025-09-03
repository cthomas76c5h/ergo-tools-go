package core

import (
	"fmt"
)

type ToolRuntime struct {
	cm      *ConversationManager
	Tenants map[string]TenantConfig
}

func NewToolRuntime(cm *ConversationManager) *ToolRuntime {
	return &ToolRuntime{cm: cm, Tenants: map[string]TenantConfig{"dev": demoTenant()}}
}

func (rt *ToolRuntime) CaptureContactInfo(sessionID string, args map[string]any) (any, error) {
	st := rt.cm.Get(sessionID)
	lead := st.Lead
	if v, ok := args["first_name"].(string); ok && v != "" {
		lead.FirstName = v
	}
	if v, ok := args["last_name"].(string); ok && v != "" {
		lead.LastName = v
	}
	if v, ok := args["email"].(string); ok && v != "" {
		lead.Email = normalizeEmail(v)
	}
	if v, ok := args["phone"].(string); ok && v != "" {
		lead.Phone = normalizePhone(v)
	}
	if v, ok := args["language"].(string); ok && v != "" {
		lead.Language = v
	}
	if v, ok := args["case_type"].(string); ok && v != "" {
		lead.CaseType = v
	}
	if v, ok := args["notes"].(string); ok && v != "" {
		lead.Notes = v
	}
	rt.cm.SetLead(sessionID, lead)
	return map[string]any{"ok": true, "lead": lead}, nil
}

func (rt *ToolRuntime) GetMissingFields(sessionID, tenantID string, args map[string]any) (any, error) {
	st := rt.cm.Get(sessionID)
	if v, ok := args["case_type"].(string); ok && v != "" && st.Lead.CaseType == "" {
		st.Lead.CaseType = v
		rt.cm.SetLead(sessionID, st.Lead)
	}
	t, ok := rt.Tenants[tenantID]
	if !ok {
		t = demoTenant()
	}
	missing := rt.cm.MissingFields(t, st.Lead)
	return map[string]any{"missing": missing}, nil
}

func (rt *ToolRuntime) ExportToCRM(sessionID, provider string) (any, error) {
	st := rt.cm.Get(sessionID)
	exportID := fmt.Sprintf("%s_%s_001", provider, sessionID)
	st.ExportID = exportID
	return map[string]any{"export_id": exportID}, nil
}

func (rt *ToolRuntime) Calendly(sessionID string, args map[string]any) (any, error) {
	return map[string]any{"booking_url": "https://calendly.com/firm/consult", "status": "placeholder"}, nil
}

func normalizeEmail(s string) string { /* TODO: john dot doe at gmail -> john.doe@gmail.com */ return s }
func normalizePhone(s string) string { /* TODO: strip non-digits, E.164 */ return s }
