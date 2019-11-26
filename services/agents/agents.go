package agents

import (
	"github.com/AlekSi/pointer"

	"github.com/percona/pmm-managed/models"
)

type redactMode int

const (
	redactSecrets redactMode = iota
	exposeSecrets
)

// redactWords returns words that should be redacted from given Agent logs/output.
func redactWords(agent *models.Agent) []string {
	var words []string
	if s := pointer.GetString(agent.Password); s != "" {
		words = append(words, s)
	}
	if s := pointer.GetString(agent.AWSSecretKey); s != "" {
		words = append(words, s)
	}
	return words
}
