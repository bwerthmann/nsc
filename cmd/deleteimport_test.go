/*
 * Copyright 2018 The NATS Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"testing"

	"github.com/nats-io/jwt"
	"github.com/stretchr/testify/require"
)

func Test_DeleteImport(t *testing.T) {
	ts := NewTestStore(t, "delete import")
	defer ts.Done(t)

	ts.AddAccount(t, "A")
	ts.AddExport(t, "A", jwt.Stream, "foo", false)
	ts.AddExport(t, "A", jwt.Stream, "bar", false)

	ts.AddAccount(t, "B")
	ts.AddImport(t, "A", "foo", "B")
	ts.AddImport(t, "A", "bar", "B")

	tests := CmdTests{
		{createDeleteImportCmd(), []string{"delete", "import", "--account", "A"}, nil, []string{"account \"A\" doesn't have imports"}, true},
		{createDeleteImportCmd(), []string{"delete", "import", "--account", "B"}, nil, []string{"subject is required"}, true},
		{createDeleteImportCmd(), []string{"delete", "import", "--account", "B", "--subject", "baz"}, nil, []string{"no stream import matching \"baz\" found"}, true},
		{createDeleteImportCmd(), []string{"delete", "import", "--account", "B", "--subject", "foo"}, nil, []string{"deleted stream import \"foo\""}, false},
	}

	tests.Run(t, "root", "delete")
}

func Test_DeleteImportAccountRequired(t *testing.T) {
	ts := NewTestStore(t, "delete import")
	defer ts.Done(t)

	ts.AddAccount(t, "A")
	ts.AddExport(t, "A", jwt.Stream, "foo", false)
	ts.AddAccount(t, "B")
	ts.AddImport(t, "A", "foo", "B")

	GetConfig().SetAccount("")
	_, _, err := ExecuteCmd(createDeleteImportCmd(), "--subject", "A")
	require.Error(t, err)
	require.Contains(t, err.Error(), "account is required")
}

func Test_DeleteImportInteractive(t *testing.T) {
	ts := NewTestStore(t, "test")
	defer ts.Done(t)

	ts.AddAccount(t, "A")
	ts.AddExport(t, "A", jwt.Stream, "foo", false)
	ts.AddExport(t, "A", jwt.Stream, "bar", false)

	ts.AddAccount(t, "B")
	ts.AddImport(t, "A", "foo", "B")
	ts.AddImport(t, "A", "bar", "B")

	input := []interface{}{1, false, 0, 0}
	cmd := createDeleteImportCmd()
	HoistRootFlags(cmd)
	_, _, err := ExecuteInteractiveCmd(cmd, input)
	require.NoError(t, err)

	ac, err := ts.Store.ReadAccountClaim("B")
	require.NoError(t, err)
	require.Len(t, ac.Imports, 1)
}

func Test_DeleteImportInteractiveNoService(t *testing.T) {
	ts := NewTestStore(t, "test")
	defer ts.Done(t)

	ts.AddAccount(t, "A")
	ts.AddExport(t, "A", jwt.Stream, "foo", false)
	ts.AddExport(t, "A", jwt.Stream, "bar", false)

	ts.AddAccount(t, "B")
	ts.AddImport(t, "A", "foo", "B")
	ts.AddImport(t, "A", "bar", "B")

	input := []interface{}{1, true, 0, ts.OperatorKeyPath}
	cmd := createDeleteImportCmd()
	HoistRootFlags(cmd)
	_, _, err := ExecuteInteractiveCmd(cmd, input)
	require.Contains(t, err.Error(), "no service imports defined")
}

func Test_DeleteImportInteractiveNoStreams(t *testing.T) {
	ts := NewTestStore(t, "test")
	defer ts.Done(t)

	ts.AddAccount(t, "A")
	ts.AddExport(t, "A", jwt.Service, "foo", false)
	ts.AddExport(t, "A", jwt.Service, "bar", false)

	ts.AddAccount(t, "B")
	ts.AddImport(t, "A", "foo", "B")
	ts.AddImport(t, "A", "bar", "B")

	input := []interface{}{1, false, 0, ts.OperatorKeyPath}
	cmd := createDeleteImportCmd()
	HoistRootFlags(cmd)
	_, _, err := ExecuteInteractiveCmd(cmd, input)
	require.Contains(t, err.Error(), "no stream imports defined")
}
