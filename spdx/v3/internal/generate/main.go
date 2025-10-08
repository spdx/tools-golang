package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kzantow/go-ld/shaclgen"
)

func main() {
	versions := []string{
		//"3.0.0",
		"3.0.1",
	}

	for _, version := range versions {
		packageName := "v" + strings.ReplaceAll(version, ".", "_")
		fileName, err := filepath.Abs(fmt.Sprintf("spdx/v3/%s/model.go", packageName))
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(fmt.Sprintf("spdx/v3/%s", packageName), 0o755)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Generating SPDX %s to %s\n", version, fileName)

		shaclgen.Generate(
			shaclgen.EnableLog(),
			shaclgen.PackageName(packageName),
			shaclgen.LicenseID("MIT"),
			shaclgen.OutputFile(fileName),
			shaclgen.RenameFunc(renameFunc),
			shaclgen.JsonLDContext(fmt.Sprintf("https://spdx.org/rdf/%s/spdx-context.jsonld", version)),
			shaclgen.SHACLTypes(fmt.Sprintf("https://spdx.org/rdf/%s/spdx-model.ttl", version)),
		)

		contents, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		// hack to work around go-ld not compacting certain types
		regexp.MustCompile(`(?m)": map[string]any{.*?},\n`).ReplaceAll(contents, []byte(`": map[string]any{
		},`))
	}
}

func renameFunc(typ shaclgen.NameType, name string, c *shaclgen.Class) string {
	if typ == shaclgen.NameTypeField {
		// get rid of property stutters
		if c != nil && c.GoName != name && !strings.EqualFold(name, "PackageURL") {
			shortName := strings.TrimPrefix(name, c.GoName)
			if !strings.EqualFold(shortName, "id") {
				name = shortName
			}
		}

		return replaceSuffixes(name, map[string]string{
			"Bies":           "By",
			"Tos":            "To",
			"CreatedUsings":  "CreatedUsing",
			"VerifiedUsings": "VerifiedUsing",
			"Id":             "ID",
			"Url":            "URL",
		})
	}
	switch name {
	case "AnyLicenseInfo":
		return "LicenseInfo"
	case "Bom":
		return "BOM"
	case "BOMS":
		return "BOMs"
	case "Sbom":
		return "SBOM"
	case "SBOMS":
		return "SBOMs"
	}
	return ""
}

func replaceSuffixes(value string, suffixToReplacement map[string]string) string {
	for suffix, replacement := range suffixToReplacement {
		if strings.HasSuffix(value, suffix) {
			value = strings.TrimSuffix(value, suffix)
			return value + replacement
		}
	}
	return value
}
