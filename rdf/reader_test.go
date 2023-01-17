// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spdx/tools-golang/spdx"
)

func Test_Read(t *testing.T) {
	fileName := "../examples/sample-docs/rdf/SPDXRdfExample-v2.2.spdx.rdf"

	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("error opening File: %s", err))
	}

	got, err := Read(file)
	if err != nil {
		t.Errorf("rdf.Read() error = %v", err)
		return
	}

	assert.IsType(t, &spdx.Document{}, got)
}
