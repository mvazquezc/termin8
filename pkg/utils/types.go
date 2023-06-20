package utils

type StuckResource struct {
	ResourceName      string
	ResourceType      string
	ResourceNamespace string
	ResourceGroup     string
	ResourceVersion   string
}

type RunResult struct {
	Namespace           string   `json:"namespace"`
	TerminatedResources []string `json:"terminated_resources"`
}
