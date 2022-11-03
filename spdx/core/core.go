// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package core

import (
	"fmt"
	"reflect"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_1"
	"github.com/spdx/tools-golang/spdx/v2_2"
)

// Metadata includes additional information to map back and forth
// to other SPDX types
type Metadata struct {
	SpdxVersion string
}

var (
	cmV2_1 Metadata = Metadata{common.SpdxV2_1}
	cmV2_2 Metadata = Metadata{common.SpdxV2_2}
)

func Convert(s interface{}) interface{} {
	// have variable to track if return type is pointer or not
	switch s.(type) {

	case *v2_1.File:
		f := File{}
		remainder, err := copyOver(&f, s)
		if err != nil {
			panic(err)
		}
		for _, r := range remainder {
			// do stuff
			fmt.Println(r)
		}
		return &f
	case *v2_2.File:
		return File{}
	default:
		return nil
	}
}
func convertFile2_2(s v2_2.File) File {
	return File{}
}

type fieldValue struct {
	field string
	value interface{}
}

// Asssume dst is a core struct or pointer thereof
func copyOver(dst, src interface{}) ([]fieldValue, error) {
	var notHandled []fieldValue

	// Assume the inputs to these objects are pointers to the
	// actual objects
	srcV := reflect.ValueOf(src).Elem()
	dstV := reflect.ValueOf(dst).Elem()
	dstT := dstV.Type()

	for i := 0; i < dstV.NumField(); i++ {
		f := dstV.Field(i)
		fmt.Printf("dst %d: %s %s = %v (Settable? %v) \n", i,
			dstT.Field(i).Name, f.Type(), f.Interface(), f.CanSet())

		fs := srcV.FieldByName(dstT.Field(i).Name)
		if !fs.IsValid() {
			fmt.Printf("skipping field because it is not valid in src %v\n", dstT.Field(i).Name)
			continue
		}
		fmt.Printf("src %d: %s %s = %v (Settable? %v) \n", i,
			dstT.Field(i).Name, fs.Type(), fs.Interface(), fs.CanSet())

		if f.Type() != fs.Type() {
			fmt.Printf("skipping field because types do not match %v (%s vs %s)\n",
				dstT.Field(i).Name, f.Type(), fs.Type())
			// TODO(lumjjb): Add check for simple conversion and recurse of Convert
			// TODO (lumjjb): Handle slices of objects here
			// TODO (lumjjb): Handle maps with either simple types / common type here (check if you change the value type name if it becomes the same type)

			notHandled = append(notHandled, fieldValue{dstT.Field(i).Name, fs.Interface()})
			continue
		}
		f.Set(fs)
	}

	// fieldsToCopy := reflect.VisibleFields(dstT)
	// for _, field := range fieldsToCopy {
	// 	fName := field.Name
	// 	if ok := srcV.FieldByName(fName).IsValid(); ok {
	// 		v := reflect.ValueOf(srcV.FieldByName(fName))
	// 		dstV.FieldByName(fName).Set(v)
	// 	}
	// }
	return notHandled, nil
}

/*

// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type A struct {
	F1 string
	F2 string
	FA string
}

type B struct {
	F1 string
	F2 string
	FB string
}

func main() {
	a := A{
		F1: "f1",
		F2: "f2",
		FA: "fa",
	}
	//fmt.Printf("a: %+v", a)

	b := A{}
	copyOver(a, b)

}

func copyOver(src, dst interface{}) error {
	dstT := reflect.TypeOf(dst)

	srcV := reflect.ValueOf(src)
	dstV := reflect.ValueOf(dst)

	fieldsToCopy := reflect.VisibleFields(dstT)
	for _, field := range fieldsToCopy {
		fName := field.Name
		fmt.Println(fName)
		if ok := srcV.FieldByName(fName).IsValid(); ok {
			srcF := srcV.FieldByName(fName)
			//fmt.Println(srcF)

			dstF := dstV.FieldByName(fName)
			dstF.Set(reflect.ValueOf(srcF.Interface()))
			fmt.Println(dstF)
			//			dstV.FieldByName(fName).Set(v)
		}
	}
	return nil
}

*/
