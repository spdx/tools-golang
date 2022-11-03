package core

import (
	"fmt"
	"testing"

	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

func TestCopyOver(t *testing.T) {
	f1 := v2_2.File{}
	f2 := v2_1.File{
		FileName: "abc",
	}

	notHandled, err := copyOver(&f1, &f2)
	if err == nil {
		for _, e := range notHandled {
			fmt.Printf("field not handled: %s: %+v\n", e.field, e.value)
		}
	}

	// Convert test
	fcore := Convert(&f2)
	fmt.Printf("Converted to core: %+v\n", fcore)
}
