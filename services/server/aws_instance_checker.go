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

package server

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/reform.v1"

	"github.com/percona/pmm-managed/models"
)

type Checker struct {
	db               *reform.DB
	telemetryService telemetryService
	l                *logrus.Entry

	rw      sync.RWMutex
	checked bool
}

func NewChecker(db *reform.DB, telemetryService telemetryService) *Checker {
	return &Checker{
		db:               db,
		telemetryService: telemetryService,
		l:                logrus.WithField("component", "server/checker"),
	}
}

func (c *Checker) NeedsCheck() bool {
	// fast-path
	c.rw.RLock()
	checked := c.checked
	c.rw.RUnlock()
	if checked {
		return false
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	// TODO
	// if c.telemetryService.DistributionMethod() != serverpb.DistributionMethod_AMI {
	// 	c.checked = true
	// 	return false
	// }

	settings, err := models.GetSettings(c.db.Querier)
	if err != nil {
		c.l.Error(err)
		return true
	}
	if settings.AWSInstanceIDChecked {
		c.checked = true
		return false
	}

	return true
}

func (c *Checker) CheckInstanceID(instanceID string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	doc, err := ec2metadata.New(sess).GetInstanceIdentityDocument()
	if err != nil {
		return errors.Wrap(err, "cannot get Instance Identity Document to validate the instance ID")
	}

	if instanceID == doc.InstanceID {
		return nil
	}

	return status.Error(codes.InvalidArgument, "invalid instance ID")
}
