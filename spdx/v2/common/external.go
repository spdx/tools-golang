// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package common

// Constants for various string types
const (

	// F.2 Security types
	CategorySecurity      string = "SECURITY"
	TypeSecurityCPE23Type string = "cpe23Type"
	TypeSecurityCPE22Type string = "cpe22Type"
	TypeSecurityAdvisory  string = "advisory"
	TypeSecurityFix       string = "fix"
	TypeSecurityUrl       string = "url"
	TypeSecuritySwid      string = "swid"

	// F.3 Package-Manager types
	CategoryPackageManager         string = "PACKAGE-MANAGER"
	TypePackageManagerMavenCentral string = "maven-central"
	TypePackageManagerNpm          string = "npm"
	TypePackageManagerNuGet        string = "nuget"
	TypePackageManagerBower        string = "bower"
	TypePackageManagerPURL         string = "purl"

	// F.4 Persistent-Id types
	CategoryPersistentId   string = "PERSISTENT-ID"
	TypePersistentIdSwh    string = "swh"
	TypePersistentIdGitoid string = "gitoid"

	// F.5 Other
	CategoryOther string = "OTHER"

	// 11.1 Relationship field types
	TypeRelationshipDescribe                  string = "DESCRIBES"
	TypeRelationshipDescribeBy                string = "DESCRIBED_BY"
	TypeRelationshipContains                  string = "CONTAINS"
	TypeRelationshipContainedBy               string = "CONTAINED_BY"
	TypeRelationshipDependsOn                 string = "DEPENDS_ON"
	TypeRelationshipDependencyOf              string = "DEPENDENCY_OF"
	TypeRelationshipBuildDependencyOf         string = "BUILD_DEPENDENCY_OF"
	TypeRelationshipDevDependencyOf           string = "DEV_DEPENDENCY_OF"
	TypeRelationshipOptionalDependencyOf      string = "OPTIONAL_DEPENDENCY_OF"
	TypeRelationshipProvidedDependencyOf      string = "PROVIDED_DEPENDENCY_OF"
	TypeRelationshipTestDependencyOf          string = "TEST_DEPENDENCY_OF"
	TypeRelationshipRuntimeDependencyOf       string = "RUNTIME_DEPENDENCY_OF"
	TypeRelationshipExampleOf                 string = "EXAMPLE_OF"
	TypeRelationshipGenerates                 string = "GENERATES"
	TypeRelationshipGeneratedFrom             string = "GENERATED_FROM"
	TypeRelationshipAncestorOf                string = "ANCESTOR_OF"
	TypeRelationshipDescendantOf              string = "DESCENDANT_OF"
	TypeRelationshipVariantOf                 string = "VARIANT_OF"
	TypeRelationshipDistributionArtifact      string = "DISTRIBUTION_ARTIFACT"
	TypeRelationshipPatchFor                  string = "PATCH_FOR"
	TypeRelationshipPatchApplied              string = "PATCH_APPLIED"
	TypeRelationshipCopyOf                    string = "COPY_OF"
	TypeRelationshipFileAdded                 string = "FILE_ADDED"
	TypeRelationshipFileDeleted               string = "FILE_DELETED"
	TypeRelationshipFileModified              string = "FILE_MODIFIED"
	TypeRelationshipExpandedFromArchive       string = "EXPANDED_FROM_ARCHIVE"
	TypeRelationshipDynamicLink               string = "DYNAMIC_LINK"
	TypeRelationshipStaticLink                string = "STATIC_LINK"
	TypeRelationshipDataFileOf                string = "DATA_FILE_OF"
	TypeRelationshipTestCaseOf                string = "TEST_CASE_OF"
	TypeRelationshipBuildToolOf               string = "BUILD_TOOL_OF"
	TypeRelationshipDevToolOf                 string = "DEV_TOOL_OF"
	TypeRelationshipTestOf                    string = "TEST_OF"
	TypeRelationshipTestToolOf                string = "TEST_TOOL_OF"
	TypeRelationshipDocumentationOf           string = "DOCUMENTATION_OF"
	TypeRelationshipOptionalComponentOf       string = "OPTIONAL_COMPONENT_OF"
	TypeRelationshipMetafileOf                string = "METAFILE_OF"
	TypeRelationshipPackageOf                 string = "PACKAGE_OF"
	TypeRelationshipAmends                    string = "AMENDS"
	TypeRelationshipPrerequisiteFor           string = "PREREQUISITE_FOR"
	TypeRelationshipHasPrerequisite           string = "HAS_PREREQUISITE"
	TypeRelationshipRequirementDescriptionFor string = "REQUIREMENT_DESCRIPTION_FOR"
	TypeRelationshipSpecificationFor          string = "SPECIFICATION_FOR"
	TypeRelationshipOther                     string = "OTHER"
)
