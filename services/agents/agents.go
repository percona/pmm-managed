package agents

import (
	"github.com/AlekSi/pointer"

	"github.com/percona/pmm-managed/models"
)

type debugValue int

const (
	enableDebug debugValue = iota
	disableDebug
)

func redactKeywords(s *models.Agent, debug debugValue) []string {
	var hideKeywords []string
	if s.Password != nil && debug == disableDebug {
		hideKeywords = append(hideKeywords, pointer.GetString(s.Password))
	}
	return hideKeywords
}
