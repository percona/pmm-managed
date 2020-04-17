package alertmanager

import (
	"fmt"

	"github.com/percona/pmm/api/alertmanager/ammodels"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"

	"github.com/percona/pmm-managed/models"
)

// Severity defines alert severity.
type Severity string

// severities
const (
	Error   = Severity("error")
	Warning = Severity("warning")
	Info    = Severity("info")
)

// AlertParams defines alert parameters.
type AlertParams struct {
	Name        string
	Summary     string
	Description string
	Severity    Severity

	Node    *models.Node
	Service *models.Service
	Agent   *models.Agent
}

// validate checks parameters and fills defaults.
func (ap *AlertParams) validate() error {
	if ap.Name == "" {
		return errors.New("empty Name")
	}
	if ap.Summary == "" {
		return errors.New("empty Summary")
	}
	if ap.Description == "" {
		return errors.New("empty Description")
	}

	if ap.Severity == "" {
		ap.Severity = Info
	}

	return nil
}

// makeAlert makes alert from given parameters.
func makeAlert(params *AlertParams) (*ammodels.PostableAlert, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	labels, err := models.MergeLabels(params.Node, params.Service, params.Agent)
	if err != nil {
		return nil, err
	}

	labels[model.AlertNameLabel] = params.Name
	labels["severity"] = string(params.Severity)

	return &ammodels.PostableAlert{
		Alert: ammodels.Alert{
			// GeneratorURL: "TODO",
			Labels: labels,
		},

		// StartsAt and EndAt can't be added there without changes in registry

		Annotations: map[string]string{
			"summary":     params.Summary,
			"description": params.Description,
		},
	}, nil
}

// makeAlertPMMAgentNotConnected makes pmm_agent_not_connected alert.
func makeAlertPMMAgentNotConnected(agent *models.Agent, node *models.Node) (string, *ammodels.PostableAlert, error) {
	name := "pmm_agent_not_connected"
	alert, err := makeAlert(&AlertParams{
		Name:        name,
		Summary:     "pmm-agent is not connected to PMM Server",
		Description: fmt.Sprintf("Node name: %s", node.NodeName),
		Severity:    Warning,

		Node:  node,
		Agent: agent,
	})
	if err != nil {
		return "", nil, err
	}
	return name, alert, nil
}

// makeAlertPMMAgentIsOutdated makes pmm_agent_outdated alert.
func makeAlertPMMAgentIsOutdated(agent *models.Agent, node *models.Node, serverVersion string) (string, *ammodels.PostableAlert, error) {
	name := "pmm_agent_outdated"
	alert, err := makeAlert(&AlertParams{
		Name:    name,
		Summary: "pmm-agent is outdated",
		Description: fmt.Sprintf(
			"Node name: %s\npmm-agent version: %s\nPMM Server version: %s",
			node.NodeName, *agent.Version, serverVersion,
		),
		Severity: Info,

		Node:  node,
		Agent: agent,
	})
	if err != nil {
		return "", nil, err
	}
	return name, alert, nil
}
