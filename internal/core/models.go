package core

import "time"

type Lead struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CaseType  string `json:"case_type,omitempty"`
	Language  string `json:"language,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

type TenantConfig struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Languages        []string            `json:"languages"`
	AllowedCaseTypes []string            `json:"allowed_case_types"`
	RequiredByCase   map[string][]string `json:"required_fields_by_case"`
}

func demoTenant() TenantConfig {
	return TenantConfig{
		ID: "dev", Name: "Demo Firm",
		Languages:        []string{"en"},
		AllowedCaseTypes: []string{"personal_injury", "family", "criminal", "immigration", "bankruptcy", "real_estate", "wills_trusts"},
		RequiredByCase: map[string][]string{
			"default":         {"first_name", "phone", "case_type"},
			"personal_injury": {"first_name", "phone", "incident_date"},
		},
	}
}

type SessionState struct {
	Lead     Lead
	History  []Message
	TenantID string
	ExportID string
	Updated  time.Time
}

type Message struct {
	Role    string
	Content string
}
