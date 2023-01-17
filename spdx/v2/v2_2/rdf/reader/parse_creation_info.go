// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package reader

import (
	"fmt"

	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_2"
)

// Cardinality: Mandatory, one.
func (parser *rdfParser2_2) parseCreationInfoFromNode(ci *v2_2.CreationInfo, node *gordfParser.Node) error {
	for _, triple := range parser.nodeToTriples(node) {
		switch triple.Predicate.ID {
		case SPDX_LICENSE_LIST_VERSION: // 2.7
			// cardinality: max 1
			ci.LicenseListVersion = triple.Object.ID
		case SPDX_CREATOR: // 2.8
			// cardinality: min 1
			err := setCreator(triple.Object.ID, ci)
			if err != nil {
				return err
			}
		case SPDX_CREATED: // 2.9
			// cardinality: exactly 1
			ci.Created = triple.Object.ID
		case RDFS_COMMENT: // 2.10
			ci.CreatorComment = triple.Object.ID
		case RDF_TYPE:
			continue
		default:
			return fmt.Errorf("unknown predicate %v while parsing a creation info", triple.Predicate)
		}
	}
	return nil
}

func setCreator(creatorStr string, ci *v2_2.CreationInfo) error {
	entityType, entity, err := ExtractSubs(creatorStr, ":")
	if err != nil {
		return fmt.Errorf("error setting creator of a creation info: %s", err)
	}

	creator := common.Creator{Creator: entity}

	switch entityType {
	case "Person", "Organization", "Tool":
		creator.CreatorType = entityType
	default:
		return fmt.Errorf("unknown creatorType %v in a creation info", entityType)
	}

	ci.Creators = append(ci.Creators, creator)

	return nil
}
