package agents

import (
	"github.com/AlekSi/pointer"

	"github.com/percona/pmm-managed/models"
)

type debugValue bool

var enableDebug = debugValue(true)
var disableDebug = debugValue(false)

func redactKeywords(s *models.Agent, debug debugValue) []string {
	var hideKeywords []string
	if s.Password != nil && debug == disableDebug {
		hideKeywords = append(hideKeywords, pointer.GetString(s.Password))
	}
	return hideKeywords
}
