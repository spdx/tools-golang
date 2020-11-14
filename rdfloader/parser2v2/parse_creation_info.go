// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"
	gordfParser "github.com/spdx/gordf/rdfloader/parser"
	"github.com/spdx/tools-golang/spdx"
)

// Cardinality: Mandatory, one.
func (parser *rdfParser2_2) parseCreationInfoFromNode(ci *spdx.CreationInfo2_2, node *gordfParser.Node) error {
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

func setCreator(creator string, ci *spdx.CreationInfo2_2) error {
	entityType, entity, err := ExtractSubs(creator, ":")
	if err != nil {
		return fmt.Errorf("error setting creator of a creation info: %s", err)
	}
	switch entityType {
	case "Person":
		ci.CreatorPersons = append(ci.CreatorPersons, entity)
	case "Organization":
		ci.CreatorOrganizations = append(ci.CreatorOrganizations, entity)
	case "Tool":
		ci.CreatorTools = append(ci.CreatorTools, entity)
	default:
		return fmt.Errorf("unknown creatorType %v in a creation info", entityType)
	}
	return nil
}
