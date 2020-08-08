package eveapimodels

// DeploymentPlanType data structure
type DeploymentPlanType string

// String returns the string deployment type
func (dpt DeploymentPlanType) String() string {
	return string(dpt)
}

// DeploymentPlanOptions data structure
type DeploymentPlanOptions struct {
	Artifacts        ArtifactDefinitions `json:"artifacts"`
	ForceDeploy      bool                `json:"force_deploy"`
	User             string              `json:"user"`
	DryRun           bool                `json:"dry_run"`
	CallbackURL      string              `json:"callback_url"`
	Environment      string              `json:"environment"`
	NamespaceAliases StringList          `json:"namespaces,omitempty"`
	Messages         []string            `json:"messages,omitempty"`
	Type             DeploymentPlanType  `json:"type"`
}
