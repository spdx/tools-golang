// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package rdf2v1

import (
	"fmt"
	"strings"

	"github.com/deltamobile/goraptor"
)

type updater func(goraptor.Term) error

// takes in string updates string
func update(val *ValueStr) updater {
	key := false
	return func(term goraptor.Term) error {
		if key {
			return fmt.Errorf("Property Already Defined")
		}
		val.Val = termStr(term)
		key = true
		return nil
	}
}
func updateTrimPrefix(prefix string, ptr *ValueStr) updater {
	key := false
	return func(term goraptor.Term) error {
		if key {
			return fmt.Errorf("Property Already Defined")
		}

		ptr.Val = strings.TrimPrefix(termStr(term), prefix)
		key = true
		return nil
	}
}

func updateList(sl *[]ValueStr) updater {
	return func(term goraptor.Term) error {
		*sl = append(*sl, Str(termStr(term))) // convert Str to ValStr
		return nil
	}
}

// Update a ValueCreator pointer
func updateCreator(ptr *ValueCreator) updater {
	key := false
	return func(term goraptor.Term) error {
		if key {
			return fmt.Errorf("Property Already Defined")
		}
		ptr.SetValue(termStr(term))
		key = true
		return nil
	}
}

func updateListCreator(sl *[]ValueCreator) updater {
	return func(term goraptor.Term) error {
		*sl = append(*sl, ValueCreatorNew(termStr(term)))
		return nil
	}
}

// Update a ValueDate pointer
func updateDate(ptr *ValueDate) updater {
	key := false
	return func(term goraptor.Term) error {
		if key {
			return fmt.Errorf("Property Already Defined")
		}
		ptr.SetValue(termStr(term))
		key = true
		return nil
	}
}
