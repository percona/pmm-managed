package utils

import (
	"github.com/percona-platform/saas/pkg/check"

	"github.com/percona/pmm-managed/models"
)

// Target contains required info about STT check target.
type Target struct {
	AgentID       string
	ServiceID     string
	ServiceName   string
	Labels        map[string]string
	Dsn           string
	Files         map[string]string
	Tdp           *models.DelimiterPair
	TLSSkipVerify bool
}

// STTCheckResult contains the output from the check file and other information.
type STTCheckResult struct {
	CheckName string
	Interval  check.Interval
	Target    Target
	Result    check.Result
}
