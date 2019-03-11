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

package models_test

import (
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/percona/pmm-managed/utils/tests"
)

func assertCantBeNull(t *testing.T, err error, column string) {
	t.Helper()

	require.IsType(t, &pq.Error{}, err)
	pgErr := err.(*pq.Error)
	// see: https://www.postgresql.org/docs/10/errcodes-appendix.html
	assert.EqualValues(t, pq.ErrorCode("23502"), pgErr.Code)
	assert.Equal(t, fmt.Sprintf("null value in column \"%s\" violates not-null constraint", column), pgErr.Message)
}

func assertDuplicate(t *testing.T, err error, constraint string) {
	t.Helper()

	require.IsType(t, &pq.Error{}, err)
	pgErr := err.(*pq.Error)
	// see: https://www.postgresql.org/docs/10/errcodes-appendix.html
	assert.EqualValues(t, pq.ErrorCode("23505"), pgErr.Code)
	assert.Equal(t, fmt.Sprintf("duplicate key value violates unique constraint \"%s\"", constraint), pgErr.Message)
}

func TestDatabaseUniqueIndexes(t *testing.T) {
	db := tests.OpenTestPostgresDB(t)
	defer func() {
		require.NoError(t, db.Close())
	}()

	var err error

	t.Run("Nodes", func(t *testing.T) {
		// node_id
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name) " +
				"VALUES ('1', '', 'name')",
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name) " +
				"VALUES ('1', '', 'other name')",
		)
		assertDuplicate(t, err, "nodes_pkey")

		// node_name
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name) " +
				"VALUES ('2', '', 'name')",
		)
		assertDuplicate(t, err, "nodes_node_name_key")

		// machine_id
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, machine_id) " +
				"VALUES ('31', '', 'name31', 'machine-id')",
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, machine_id) " +
				"VALUES ('32', '', 'name32', 'machine-id')",
		)
		assertDuplicate(t, err, "nodes_machine_id_key")

		// docker_container_id
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, docker_container_id) " +
				"VALUES ('41', '', 'name41', 'docker-container-id')",
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, docker_container_id) " +
				"VALUES ('42', '', 'name42', 'docker-container-id')",
		)
		assertDuplicate(t, err, "nodes_docker_container_id_key")

		// (address, region)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, address, region) " +
				"VALUES ('51', '', 'name51', 'instance1', 'region1')",
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, address, region) " +
				"VALUES ('52', '', 'name52', 'instance1', 'region1')",
		)
		assertDuplicate(t, err, "nodes_address_region_key")
		// same address, NULL region is fine
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, address) " +
				"VALUES ('53', '', 'name53', 'instance1')",
		)
		require.NoError(t, err)
		_, err = db.Exec(
			"INSERT INTO nodes (node_id, node_type, node_name, address) " +
				"VALUES ('54', '', 'name54', 'instance1')",
		)
		require.NoError(t, err)
	})

	t.Run("Services", func(t *testing.T) {
		t.Skip("TODO")
	})

	t.Run("Agents", func(t *testing.T) {
		t.Skip("TODO")
	})
}
