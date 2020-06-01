// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package management

import (
	"context"
	"testing"

	"github.com/percona/pmm/api/managementpb"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/percona/pmm-managed/services"
)

func TestStartSecurityChecks(t *testing.T) {
	t.Run("internal error", func(t *testing.T) {
		var checksService mockChecksService
		checksService.On("StartChecks", mock.Anything).Return(errors.New("random error"))

		s := NewChecksAPIService(&checksService)

		resp, err := s.StartSecurityChecks(context.Background(), &managementpb.StartSecurityChecksRequest{})
		assert.EqualError(t, err, "rpc error: code = Internal desc = Failed to start security checks.")
		assert.Nil(t, resp)
	})

	t.Run("STT disabled error", func(t *testing.T) {
		var checksService mockChecksService
		checksService.On("StartChecks", mock.Anything).Return(services.ErrSTTDisabled)

		s := NewChecksAPIService(&checksService)

		resp, err := s.StartSecurityChecks(context.Background(), &managementpb.StartSecurityChecksRequest{})
		assert.EqualError(t, err, "rpc error: code = FailedPrecondition desc = STT is disabled.")
		assert.Nil(t, resp)
	})
}
