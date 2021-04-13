package eveapi

import (
	"encoding/json"
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// ChatMessage takes in an eve API Model (as an interface)
// and returns a formatted string for a Chat Message response/result
func ChatMessage(model interface{}) string {
	if model == nil {
		return ""
	}

	switch v := model.(type) {
	case eve.DeployJob:
		return deployJobMsg(v)
	case eve.DeployJobs:
		return deployJobsMsg(v)
	case *eve.DeployService:
		return deployServiceMsg(v)
	case eve.DeployService:
		return deployServiceMsg(&v)
	case eve.DeployServices:
		return deployServicesMsg(v)
	case *eve.NSDeploymentPlan:
		return nsDeployPlanMsg(v)
	case eve.Release:
		return releaseMsg(v)
	case eve.Metadata:
		return metadataMsg(v)
	case []eve.Service:
		return servicesMsg(v)
	case []eve.Namespace:
		return namespacesMsg(v)
	case []eve.Environment:
		return environmentsMsg(v)
	case []eve.Job:
		return jobsMsg(v)
	default:
		return ""
	}
}

func deployJobMsg(j eve.DeployJob) string {
	if j.ArtifactName == j.JobName {
		return fmt.Sprintf("\n%s:%s", j.JobName, j.AvailableVersion)
	}
	return fmt.Sprintf("\n%s (%s):%s", j.JobName, j.ArtifactName, j.AvailableVersion)
}

func deployJobsMsg(v eve.DeployJobs) string {
	msg := ""
	if msg = initListString(v, "jobs"); len(msg) == 0 {
		for _, val := range v {
			msg += deployJobMsg(*val)
		}
	}
	return msg
}

func jobMsg(v eve.Job) string {
	if v.ArtifactName == v.Name {
		return fmt.Sprintf("\n%s:%s", v.Name, v.DeployedVersion)
	}
	return fmt.Sprintf("\n%s (%s):%s", v.Name, v.ArtifactName, v.DeployedVersion)
}

func jobsMsg(v []eve.Job) string {
	msg := ""
	if msg = initListString(v, "jobs"); len(msg) == 0 {
		for _, val := range v {
			msg += jobMsg(val)
		}
	}
	return msg
}

func deployServiceMsg(v *eve.DeployService) string {
	if v.ArtifactName == v.ServiceName {
		return fmt.Sprintf("\n%s:%s", v.ServiceName, v.AvailableVersion)
	}
	return fmt.Sprintf("\n%s (%s):%s", v.ServiceName, v.ArtifactName, v.AvailableVersion)
}

func deployServicesMsg(v eve.DeployServices) string {
	msg := ""
	if msg = initListString(v, "services"); len(msg) == 0 {
		for _, svc := range v {
			if len(msg) == 0 {
				msg = ChatMessage(svc)
			} else {
				msg += ChatMessage(svc)
			}
		}
	}
	return msg
}

func metadataMsg(v eve.Metadata) string {
	if v.ID == 0 || len(v.Value) == 0 {
		return "no metadata"
	}
	jsonB, err := json.MarshalIndent(v.Value, "", "	")
	if err != nil {
		return "invalid json metadata"
	}
	return "```" + string(jsonB) + "```"
}

func servicesMsg(v []eve.Service) string {
	msg := ""
	if msg = initListString(v, "services"); len(msg) == 0 {
		for _, val := range v {
			msg += "`" + val.Name + "` - _" + val.DeployedVersion + "_ ( *" + val.ArtifactName + "* )" + "\n"
		}
	}
	return msg
}

func namespacesMsg(v []eve.Namespace) string {
	msg := ""
	if msg = initListString(v, "namespaces"); len(msg) == 0 {
		for _, val := range v {
			msg += "`" + val.Alias + "` ( _" + val.RequestedVersion + "_ )" + "\n"
		}
	}
	return msg
}

func environmentsMsg(v []eve.Environment) string {
	msg := ""
	if msg = initListString(v, "environments"); len(msg) == 0 {
		for _, val := range v {
			msg += "Name: " + "`" + val.Name + "`" + "\n" + "Description: " + "_" + val.Description + "_" + "\n\n"
		}
	}
	return msg
}

func nsDeployPlanMsg(v *eve.NSDeploymentPlan) string {
	if v == nil || v.Status == eve.DeploymentPlanStatusMessage || v.Namespace == nil {
		return ""
	}
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s\nCluster: %s```", v.Namespace.Alias, v.EnvironmentName, v.Namespace.ClusterName)
}

func releaseMsg(v eve.Release) string {
	return fmt.Sprintf("Artifact: `%s`\nVersion: `%s`\nFrom: `%s`\nTo: `%s`", v.Artifact, v.Version, v.FromFeed, v.ToFeed)
}

func initListString(v interface{}, s string) string {
	if v != nil {
		if items, ok := v.([]interface{}); ok {
			if len(items) == 0 {
				return fmt.Sprintf("no %s", s)
			}
		}
	}
	return ""
}
