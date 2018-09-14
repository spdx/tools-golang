// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package tvloader

import "github.com/swinslow/spdx-go/v0/spdx"

func (parser *tvParser2_1) parsePairFromPackage2_1(tag string, value string) error {
	switch tag {
	case "PackageName":
		// if package already has a name, create and go on to a new package
		if parser.pkg.PackageName != "" {
			parser.pkg = &spdx.Package2_1{IsUnpackaged: false}
			parser.doc.Packages = append(parser.doc.Packages, parser.pkg)
		}
		parser.pkg.PackageName = value
	}

	return nil
}
