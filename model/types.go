package model

type Intent struct {
	Name  string       `json:"name"`
	Steps []IntentStep `json:"steps"`
}

type IntentStep struct {
	Match    []string `json:"match"`
	Respones []string `json:"responses"`
	Fallback []string `json:"fallback"`
}

type IntentClass struct {
	Intent string   `json:"intent"`
	Words  []string `json:"words"`
	CurrentStep int `json:"currentStep"`
}
