package core

import (
	"sync"
)

type ConversationManager struct {
	mu sync.Mutex
	s  map[string]*SessionState
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{s: map[string]*SessionState{}}
}

func (m *ConversationManager) Get(sessionID string) *SessionState {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.s[sessionID]
	if !ok {
		st = &SessionState{Lead: Lead{Language: "en"}}
		m.s[sessionID] = st
	}
	return st
}

func (m *ConversationManager) SetLead(sessionID string, lead Lead) {
	st := m.Get(sessionID)
	m.mu.Lock()
	st.Lead = lead
	m.mu.Unlock()
}

func (m *ConversationManager) MissingFields(t TenantConfig, lead Lead) []string {
	caseKey := lead.CaseType
	if caseKey == "" {
		caseKey = "default"
	}
	req := t.RequiredByCase[caseKey]
	if req == nil {
		req = t.RequiredByCase["default"]
	}
	missing := []string{}
	for _, f := range req {
		switch f {
		case "first_name":
			if lead.FirstName == "" {
				missing = append(missing, f)
			}
		case "phone":
			if lead.Phone == "" {
				missing = append(missing, f)
			}
		case "case_type":
			if lead.CaseType == "" {
				missing = append(missing, f)
			}
		case "incident_date": /* example */
			missing = append(missing, f)
		}
	}
	return missing
}
