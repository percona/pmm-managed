package agents

import "github.com/percona/pmm-managed/models"

func init() {
	models.HashPassword = func(password string) string {
		return password
	}
}
