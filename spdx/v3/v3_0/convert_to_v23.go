package v3_0

import (
	"github.com/spdx/tools-golang/spdx/v2/common" 
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"fmt"
)

//we want conversion from v3.0 -> v2.3 for name specifically v3.GetName -> v2PackageName

// Extract Name from v3 package and place it in v2 package
func ConvertPackageNameToV23(p AnyPackage) *v2_3.Package {
	if p == nil {
		return nil
	}
	name := p.GetName()
	return &v2_3.Package{
		PackageName: name,
	}
}

//Extract V3ID -> Check mapping -> Generate NewID -> StoreMapping -> Assign To v2

func ConvertPackageToV23(p AnyPackage, idMap map[string]common.ElementID, idx int) *v2_3.Package {
	if p == nil {
		return nil
	}

	v3ID := p.GetID()

	// Check if already mapped
	if existing, ok := idMap[v3ID]; ok {
		return &v2_3.Package{
			PackageName:           p.GetName(),
			PackageSPDXIdentifier: existing,
		}
	}

	// Generate new SPDX ID
	spdxID := fmt.Sprintf("SPDXRef-Package-%d", idx+1)

	// Store mapping
	idMap[v3ID] = common.ElementID(spdxID)

	return &v2_3.Package{
		PackageName:           p.GetName(),
		PackageSPDXIdentifier: common.ElementID(spdxID),
	}
}