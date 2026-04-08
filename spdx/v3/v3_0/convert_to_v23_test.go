package v3_0

import (
	"testing"
	"github.com/spdx/tools-golang/spdx/v2/common"
)

func TestConvertPackageNameToV23(t *testing.T) {

	t.Run("basic name mapping", func(t *testing.T) {
		pkg := &AIPackage{
			ID:   "pkg-1",
			Name: "example1",
		}

		result := ConvertPackageNameToV23(pkg)

		if result == nil {
			t.Fatalf("expected non-nil result")
		}

		if result.PackageName != "example1" {
			t.Fatalf("expected 'example1', got '%s'", result.PackageName)
		}
	})

	t.Run("nil input", func(t *testing.T) {
		result := ConvertPackageNameToV23(nil)

		if result != nil {
			t.Fatalf("expected nil result for nil input")
		}
	})
}

func TestConvertV3ElementIDtoV2ElementID(t *testing.T) {

	t.Run("basic ID mapping", func(t *testing.T) {
		idMap := make(map[string]common.ElementID)

		pkg := &AIPackage{
			ID:   "pkg-1",
			Name: "example",
		}

		result := ConvertPackageToV23(pkg, idMap, 0)

		if result.PackageSPDXIdentifier != "SPDXRef-Package-1" {
			t.Fatalf("expected SPDXRef-Package-1, got %s", result.PackageSPDXIdentifier)
		}
	})

	t.Run("reuse existing ID mapping", func(t *testing.T) {
		idMap := make(map[string]common.ElementID)

		pkg := &AIPackage{
			ID:   "pkg-1",
			Name: "example",
		}

		first := ConvertPackageToV23(pkg, idMap, 0)
		second := ConvertPackageToV23(pkg, idMap, 1)

		if first.PackageSPDXIdentifier != second.PackageSPDXIdentifier {
			t.Fatalf("expected same SPDXID, got %s and %s", first.PackageSPDXIdentifier, second.PackageSPDXIdentifier)
		}
	})

	t.Run("multiple packages get unique IDs", func(t *testing.T) {
		idMap := make(map[string]common.ElementID)

		pkg1 := &AIPackage{ID: "pkg-1", Name: "one"}
		pkg2 := &AIPackage{ID: "pkg-2", Name: "two"}

		res1 := ConvertPackageToV23(pkg1, idMap, 0)
		res2 := ConvertPackageToV23(pkg2, idMap, 1)

		if res1.PackageSPDXIdentifier == res2.PackageSPDXIdentifier {
			t.Fatalf("expected different SPDXIDs, got %s", res1.PackageSPDXIdentifier)
		}
	})
}