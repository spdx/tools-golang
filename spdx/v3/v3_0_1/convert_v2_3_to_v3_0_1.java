/**
 * SPDX-FileCopyrightText: Copyright (c) 2024 Source Auditor Inc.
 * SPDX-FileType: SOURCE
 * SPDX-License-Identifier: Apache-2.0
 * <p>
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 * <p>
 *       http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */
package org.spdx.library.conversion;

import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import javax.annotation.Nullable;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.spdx.core.IModelCopyManager;
import org.spdx.core.InvalidSPDXAnalysisException;
import org.spdx.library.ListedLicenses;
import org.spdx.library.model.v2.SpdxConstantsCompatV2;
import org.spdx.library.model.v2.SpdxCreatorInformation;
import org.spdx.library.model.v2.pointer.ByteOffsetPointer;
import org.spdx.library.model.v2.pointer.LineCharPointer;
import org.spdx.library.model.v2.pointer.StartEndPointer;
import org.spdx.library.model.v3_0_1.ModelObjectV3;
import org.spdx.library.model.v3_0_1.SpdxConstantsV3;
import org.spdx.library.model.v3_0_1.SpdxModelClassFactoryV3;
import org.spdx.library.model.v3_0_1.core.Agent;
import org.spdx.library.model.v3_0_1.core.Annotation;
import org.spdx.library.model.v3_0_1.core.AnnotationType;
import org.spdx.library.model.v3_0_1.core.CreationInfo;
import org.spdx.library.model.v3_0_1.core.Element;
import org.spdx.library.model.v3_0_1.core.ExternalElement;
import org.spdx.library.model.v3_0_1.core.ExternalIdentifier;
import org.spdx.library.model.v3_0_1.core.ExternalIdentifierType;
import org.spdx.library.model.v3_0_1.core.ExternalMap;
import org.spdx.library.model.v3_0_1.core.ExternalRefType;
import org.spdx.library.model.v3_0_1.core.Hash;
import org.spdx.library.model.v3_0_1.core.HashAlgorithm;
import org.spdx.library.model.v3_0_1.core.IntegrityMethod;
import org.spdx.library.model.v3_0_1.core.LifecycleScopeType;
import org.spdx.library.model.v3_0_1.core.LifecycleScopedRelationship;
import org.spdx.library.model.v3_0_1.core.NamespaceMap;
import org.spdx.library.model.v3_0_1.core.NoAssertionElement;
import org.spdx.library.model.v3_0_1.core.NoneElement;
import org.spdx.library.model.v3_0_1.core.Organization;
import org.spdx.library.model.v3_0_1.core.PackageVerificationCode;
import org.spdx.library.model.v3_0_1.core.Person;
import org.spdx.library.model.v3_0_1.core.Relationship;
import org.spdx.library.model.v3_0_1.core.RelationshipCompleteness;
import org.spdx.library.model.v3_0_1.core.RelationshipType;
import org.spdx.library.model.v3_0_1.core.SpdxDocument;
import org.spdx.library.model.v3_0_1.core.Tool;
import org.spdx.library.model.v3_0_1.expandedlicensing.ConjunctiveLicenseSet;
import org.spdx.library.model.v3_0_1.expandedlicensing.CustomLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.CustomLicenseAddition;
import org.spdx.library.model.v3_0_1.expandedlicensing.DisjunctiveLicenseSet;
import org.spdx.library.model.v3_0_1.expandedlicensing.ExtendableLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.ExternalCustomLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.License;
import org.spdx.library.model.v3_0_1.expandedlicensing.LicenseAddition;
import org.spdx.library.model.v3_0_1.expandedlicensing.ListedLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.ListedLicenseException;
import org.spdx.library.model.v3_0_1.expandedlicensing.NoAssertionLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.NoneLicense;
import org.spdx.library.model.v3_0_1.expandedlicensing.OrLaterOperator;
import org.spdx.library.model.v3_0_1.expandedlicensing.WithAdditionOperator;
import org.spdx.library.model.v3_0_1.simplelicensing.AnyLicenseInfo;
import org.spdx.library.model.v3_0_1.simplelicensing.LicenseExpression;
import org.spdx.library.model.v3_0_1.software.ContentIdentifierType;
import org.spdx.library.model.v3_0_1.software.Snippet;
import org.spdx.library.model.v3_0_1.software.SoftwareArtifact;
import org.spdx.library.model.v3_0_1.software.SoftwarePurpose;
import org.spdx.library.model.v3_0_1.software.SpdxFile;
import org.spdx.library.model.v3_0_1.software.SpdxPackage;
import org.spdx.storage.IModelStore;
import org.spdx.storage.IModelStore.IdType;
import org.spdx.storage.listedlicense.SpdxListedLicenseModelStore;

/**
 * Converts SPDX spec version 2.X objects to SPDX spec version 3.X and stores the result in the
 * toModelStore
 *
 * @author Gary O'Neall
 */
@SuppressWarnings({"OptionalGetWithoutIsPresent", "LoggingSimilarMessage"})
public class Spdx2to3Converter implements ISpdxConverter {

    static final Logger logger = LoggerFactory.getLogger(Spdx2to3Converter.class);

    static final Pattern SPDX_2_CREATOR_PATTERN = Pattern.compile("(Person|Organization):\\s*([^(]+)\\s*(\\(([^)]*)\\))?");

    public static final Map<org.spdx.library.model.v2.enumerations.RelationshipType, RelationshipType> RELATIONSHIP_TYPE_MAP;

    public static final Map<org.spdx.library.model.v2.enumerations.RelationshipType, LifecycleScopeType> LIFECYCLE_SCOPE_MAP;

    public static final Set<org.spdx.library.model.v2.enumerations.RelationshipType> SWAP_TO_FROM_REL_TYPES;

    public static final Map<org.spdx.library.model.v2.enumerations.AnnotationType, AnnotationType> ANNOTATION_TYPE_MAP;

    public static final Map<org.spdx.library.model.v2.enumerations.ChecksumAlgorithm, HashAlgorithm> HASH_ALGORITH_MAP;

    public static final Map<String, ContentIdentifierType> CONTENT_IDENTIFIER_TYPE_MAP;

    public static final Map<String, ExternalIdentifierType> EXTERNAL_IDENTIFIER_TYPE_MAP;

    public static final Map<String, ExternalRefType> EXTERNAL_REF_TYPE_MAP;

    public static final Map<org.spdx.library.model.v2.enumerations.Purpose, SoftwarePurpose> PURPOSE_MAP;

    static {
    }

    String toSpecVersion;
    IModelStore toModelStore;
    Map<String, String> alreadyConverted = Collections.synchronizedMap(new HashMap<>());
    CreationInfo defaultCreationInfo;
    String defaultUriPrefix;
    /**
     * Map of the documentUri to information captured from the ExternalDocumentRef from the document
     */
    Map<String, Map<Collection<ExternalMap>, ExternalMapInfo>> docUriToExternalMap = Collections.synchronizedMap(new HashMap<>());

    private final IModelCopyManager copyManager;

    private int documentIndex = 0;

    private final boolean complexLicenses;

    /**
     * @param fromObjectUri object URI of the SPDX object copied from
     * @return true if the SPDX object has already been copied
     */
    public boolean alreadyCopied(String fromObjectUri) {
        return this.alreadyConverted.containsKey(fromObjectUri);
    }

    /**
     * @param spdx2CreatorInfo SPDX 2 creation information
     * @param spdx3CreationInfo SPDX 3 creation information
     * @return true of the values of the SPDX 2 creation information are equivalent to the SPDX 3 creation information
     * @throws InvalidSPDXAnalysisException on error fetching model data
     */
    private boolean equivalentCreationInfo(SpdxCreatorInformation spdx2CreatorInfo,
                                           CreationInfo spdx3CreationInfo) throws InvalidSPDXAnalysisException {
    }

    /**
     * Coy all element properties from the SPDX spec version 2 element to the SPDX version 3 element
     * @param fromElement SPDX spec version 2 SpdxElement
     * @param toElement SPDX spec version 3 element
     * @throws InvalidSPDXAnalysisException on any errors converting element properties
     */
    private void convertElementProperties(org.spdx.library.model.v2.SpdxElement fromElement, Element toElement) throws InvalidSPDXAnalysisException {
        toElement.setCreationInfo(defaultCreationInfo);
        for (org.spdx.library.model.v2.Annotation fromAnnotation:fromElement.getAnnotations()) {
            convertAndStore(fromAnnotation, toElement);
        }
        toElement.setComment(fromElement.getComment().orElse(null));
        toElement.setName(fromElement.getName().orElse(null));
        for (org.spdx.library.model.v2.Relationship fromRelationship:fromElement.getRelationships()) {
            convertAndStore(fromRelationship, toElement);
        }
    }

    /**
     * @param fromObjectUri Object URI of the SPDX spec version 2 object being converted from
     * @param toType SPDX spec version 3 type
     * @return optional of the existing object - if it exists
     * @throws InvalidSPDXAnalysisException if there is an error creating the existing model object
     */
    protected Optional<ModelObjectV3> getExistingObject(String fromObjectUri, String toType) throws InvalidSPDXAnalysisException {
        String toObjectUri = alreadyConverted.get(fromObjectUri);
        if (Objects.isNull(toObjectUri)) {
            return Optional.empty();
        } else {
            return Optional.of(SpdxModelClassFactoryV3.getModelObject(toModelStore,
                    toObjectUri, toType, copyManager, false, defaultUriPrefix));
        }
    }

    /**
     * Converts an SPDX spec version 2 relationship to an SPDX spec version 3 relationship
     * @param fromRelationship relationship to convert from
     * @param containingElement Element which contains the property referring to the fromRelationship
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public Relationship convertAndStore(org.spdx.library.model.v2.Relationship fromRelationship,
                                        Element containingElement) throws InvalidSPDXAnalysisException {
        org.spdx.library.model.v2.enumerations.RelationshipType fromRelationshipType = fromRelationship.getRelationshipType();
        LifecycleScopeType scope = LIFECYCLE_SCOPE_MAP.get(fromRelationshipType);
        String fromUri = fromRelationship.getObjectUri();
        String relationshipType = Objects.isNull(scope) ? SpdxConstantsV3.CORE_RELATIONSHIP : SpdxConstantsV3.CORE_LIFECYCLE_SCOPED_RELATIONSHIP;
        Optional<ModelObjectV3> existing = getExistingObject(fromUri,
                relationshipType);
        if (existing.isPresent()) {
            return (Relationship)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String exitingUri = alreadyConverted.putIfAbsent(fromUri, toObjectUri);
        if (Objects.nonNull(exitingUri)) {
            return (Relationship)getExistingObject(fromUri,
                    relationshipType).get();
        }
        Relationship toRelationship;
        if (Objects.isNull(scope)) {
            toRelationship = (Relationship)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                    toObjectUri, SpdxConstantsV3.CORE_RELATIONSHIP, copyManager, true, defaultUriPrefix);
        } else {
            toRelationship = (LifecycleScopedRelationship)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                    toObjectUri, SpdxConstantsV3.CORE_LIFECYCLE_SCOPED_RELATIONSHIP, copyManager, true, defaultUriPrefix);
        }
        toRelationship.setCreationInfo(defaultCreationInfo);
        toRelationship.setRelationshipType(RELATIONSHIP_TYPE_MAP.get(fromRelationshipType));
        if (SWAP_TO_FROM_REL_TYPES.contains(fromRelationshipType)) {
            toRelationship.getTos().add(containingElement);
        } else {
            toRelationship.setFrom(containingElement);
        }
        toRelationship.setFrom(containingElement);
        toRelationship.setComment(fromRelationship.getComment().orElse(null));
        Optional<org.spdx.library.model.v2.SpdxElement> relatedSpdxElement = fromRelationship.getRelatedSpdxElement();
        RelationshipCompleteness completeness = RelationshipCompleteness.NO_ASSERTION;
        if (relatedSpdxElement.isPresent() && relatedSpdxElement.get() instanceof org.spdx.library.model.v2.SpdxNoneElement) {
            completeness = RelationshipCompleteness.COMPLETE;
        }
        toRelationship.setCompleteness(completeness);
        if (relatedSpdxElement.isPresent() &&
                !((relatedSpdxElement.get() instanceof org.spdx.library.model.v2.SpdxNoneElement) ||
                        (relatedSpdxElement.get() instanceof org.spdx.library.model.v2.SpdxNoAssertionElement))) {
            if (SWAP_TO_FROM_REL_TYPES.contains(fromRelationshipType)) {
                toRelationship.setFrom(convertAndStore(relatedSpdxElement.get()));
            } else {
                toRelationship.getTos().add(convertAndStore(relatedSpdxElement.get()));
            }
        }
        if (Objects.nonNull(scope)) {
            ((LifecycleScopedRelationship)toRelationship).setScope(scope);
        }
        toRelationship.setCreationInfo(defaultCreationInfo);
        return toRelationship;
    }

    /**
     * Converts an SPDX spec version 2 annotation to an SPDX spec version 3 annotation
     * @param fromAnnotation annotation to convert from
     * @param toElement Element which contains the property referring to the fromAnnotation
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public Annotation convertAndStore(org.spdx.library.model.v2.Annotation fromAnnotation, Element toElement) throws InvalidSPDXAnalysisException {
        String fromUri = fromAnnotation.getObjectUri();
        Optional<ModelObjectV3> existing = getExistingObject(fromUri, SpdxConstantsV3.CORE_ANNOTATION);
        if (existing.isPresent()) {
            return (Annotation)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String exitingUri = alreadyConverted.putIfAbsent(fromUri, toObjectUri);
        if (Objects.nonNull(exitingUri)) {
            return (Annotation)getExistingObject(fromUri, SpdxConstantsV3.CORE_ANNOTATION).get();
        }
        Annotation toAnnotation = (Annotation)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.CORE_ANNOTATION, copyManager, true, defaultUriPrefix);

        toAnnotation.setAnnotationType(ANNOTATION_TYPE_MAP.get(fromAnnotation.getAnnotationType()));
        toAnnotation.setStatement(fromAnnotation.getComment());
        toAnnotation.setSubject(toElement);
        CreationInfo creationInfo = new CreationInfo.CreationInfoBuilder(toModelStore, toModelStore.getNextId(IdType.Anonymous), null)
                .setCreated(fromAnnotation.getAnnotationDate())
                .setSpecVersion(SpdxConstantsV3.MODEL_SPEC_VERSION)
                .addAllCreatedUsing(defaultCreationInfo.getCreatedUsings())
                .build();
        creationInfo.setIdPrefix(defaultUriPrefix);
        creationInfo.getCreatedBys().add(stringToAgent(fromAnnotation.getAnnotator(), creationInfo));
        toAnnotation.setCreationInfo(creationInfo);
        return toAnnotation;
    }

    /**
     * Converts an SPDX spec version 2 SPDX document to an SPDX spec version 3 SPDX document and store the result
     * in the toStore
     * @param fromDoc SPDX spec version 2 document to convert from
     * @return SPDX spec version 3 document converted from the version 2 document
     * @throws InvalidSPDXAnalysisException on any errors converting the SPDX document
     */
    public SpdxDocument convertAndStore(org.spdx.library.model.v2.SpdxDocument fromDoc) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromDoc.getObjectUri(), SpdxConstantsV3.CORE_SPDX_DOCUMENT);
        if (existing.isPresent()) {
            return (SpdxDocument)existing.get();
        }
        String toObjectUri = defaultUriPrefix + "document" + documentIndex++;
        String existingUri = this.alreadyConverted.putIfAbsent(fromDoc.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (SpdxDocument)getExistingObject(fromDoc.getObjectUri(), SpdxConstantsV3.CORE_SPDX_DOCUMENT).get();
        }
        SpdxDocument toDoc = (SpdxDocument)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.CORE_SPDX_DOCUMENT, copyManager, true, defaultUriPrefix);
        // NOTE: We have to add the external doc refs first so that the ExternalMap will be properly populated
        for (org.spdx.library.model.v2.ExternalDocumentRef externalDocRef:fromDoc.getExternalDocumentRefs()) {
            toDoc.getNamespaceMaps().add(convertAndStore(externalDocRef, toDoc.getSpdxImports()));
        }
        convertElementProperties(fromDoc, toDoc);
        if (!equivalentCreationInfo(Objects.requireNonNull(fromDoc.getCreationInfo()), defaultCreationInfo)) {
            toDoc.setCreationInfo(convertCreationInfo(fromDoc.getCreationInfo(), this.toModelStore, this.defaultUriPrefix));
        }
        toDoc.setDataLicense(convertAndStore(fromDoc.getDataLicense()));
        toDoc.getRootElements().addAll(fromDoc.getDocumentDescribes().stream().map(spdxElement -> {
                    try {
                        return convertAndStore(spdxElement);
                    } catch (InvalidSPDXAnalysisException e) {
                        logger.error("Error converting SPDX elements from spec version 2 to spec version 3", e);
                        throw new RuntimeException(e);
                    }
                }
        ).collect(Collectors.toList()));
        for (org.spdx.library.model.v2.license.ExtractedLicenseInfo extractedLicense:fromDoc.getExtractedLicenseInfos()) {
            convertAndStore(extractedLicense);
        }
        return toDoc;
    }

    /**
     * Converts the externalDocRef to a NamespaceMap and store the NamespaceMap.
     * The document information is also retained in the externalDocRefMap such that any subsequent
     * references to the external documents can be captured in the ExternalMap for the document
     * @param externalDocRef SPDX Model V2 external document reference
     * @param docImports SPDX document imports to track any added external references
     * @return the namespace map correlating the documentRef ID to the document URI
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public NamespaceMap convertAndStore(org.spdx.library.model.v2.ExternalDocumentRef externalDocRef,
                                        Collection<ExternalMap> docImports) throws InvalidSPDXAnalysisException {
        // Add information to the docUriToExternalMap - this will be used if we run across any references
        // to the external document namespaces
        Optional<org.spdx.library.model.v2.Checksum> docChecksum = externalDocRef.getChecksum();
        Optional<Hash> externalDocumentHash = docChecksum.isPresent() ? Optional.of(convertAndStore(docChecksum.get())) : Optional.empty();
        docUriToExternalMap.putIfAbsent(externalDocRef.getSpdxDocumentNamespace(), Collections.synchronizedMap(new HashMap<>()));
        Map<Collection<ExternalMap>, ExternalMapInfo> externalMapInfoMap = docUriToExternalMap.get(externalDocRef.getSpdxDocumentNamespace());
        externalMapInfoMap.put(docImports, new ExternalMapInfo(externalDocRef.getId(), externalDocRef.getSpdxDocumentNamespace(),
                externalDocumentHash, docImports));
        Optional<ModelObjectV3> existing = getExistingObject(externalDocRef.getObjectUri(), SpdxConstantsV3.CORE_NAMESPACE_MAP);
        if (existing.isPresent()) {
            return (NamespaceMap)existing.get();
        }
        String toObjectUri = toModelStore.getNextId(IdType.Anonymous);
        String existingUri = this.alreadyConverted.putIfAbsent(externalDocRef.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (NamespaceMap)getExistingObject(externalDocRef.getObjectUri(), SpdxConstantsV3.CORE_NAMESPACE_MAP).get();
        }
        NamespaceMap toNamespaceMap = (NamespaceMap)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.CORE_NAMESPACE_MAP, copyManager, true, defaultUriPrefix);
        toNamespaceMap.setPrefix(externalDocRef.getId());
        toNamespaceMap.setNamespace(externalDocRef.getSpdxDocumentNamespace() + "#");
        return toNamespaceMap;
    }

    /**
     * Converts an SPDX spec version 2 SPDX ConjunctiveLicenseSet to an SPDX spec version 3 SPDX ConjunctiveLicenseSet and store the result
     * in the toStore
     * @param fromConjunctiveLicenseSet an SPDX spec version 2 ConjunctiveLicenseSet to convert from
     * @return an SPDX spec version 3 ConjunctiveLicenseSet
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public ConjunctiveLicenseSet convertAndStore(org.spdx.library.model.v2.license.ConjunctiveLicenseSet fromConjunctiveLicenseSet) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromConjunctiveLicenseSet.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CONJUNCTIVE_LICENSE_SET);
        if (existing.isPresent()) {
            return (ConjunctiveLicenseSet)existing.get();
        }
        String toObjectUri = defaultUriPrefix +  toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromConjunctiveLicenseSet.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (ConjunctiveLicenseSet)getExistingObject(fromConjunctiveLicenseSet.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CONJUNCTIVE_LICENSE_SET).get();
        }
        ConjunctiveLicenseSet toConjunctiveLicenseSet = (ConjunctiveLicenseSet)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_CONJUNCTIVE_LICENSE_SET, copyManager, true, defaultUriPrefix);
        for (org.spdx.library.model.v2.license.AnyLicenseInfo fromMember:fromConjunctiveLicenseSet.getMembers()) {
            toConjunctiveLicenseSet.getMembers().add(convertAndStore(fromMember));
        }
        toConjunctiveLicenseSet.setCreationInfo(defaultCreationInfo);
        return toConjunctiveLicenseSet;
    }

    /**
     * Converts an SPDX spec version 2 SPDX DisjunctiveLicenseSet to an SPDX spec version 3 SPDX DisjunctiveLicenseSet and store the result
     * in the toStore
     * @param fromDisjunctiveLicenseSet an SPDX spec version 2 DisjunctiveLicenseSet to convert from
     * @return an SPDX spec version 3 DisjunctiveLicenseSet
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public DisjunctiveLicenseSet convertAndStore(org.spdx.library.model.v2.license.DisjunctiveLicenseSet fromDisjunctiveLicenseSet) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromDisjunctiveLicenseSet.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_DISJUNCTIVE_LICENSE_SET);
        if (existing.isPresent()) {
            return (DisjunctiveLicenseSet)existing.get();
        }
        String toObjectUri = defaultUriPrefix +  toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromDisjunctiveLicenseSet.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (DisjunctiveLicenseSet)getExistingObject(fromDisjunctiveLicenseSet.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_DISJUNCTIVE_LICENSE_SET).get();
        }
        DisjunctiveLicenseSet toDisjunctiveLicenseSet = (DisjunctiveLicenseSet)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_DISJUNCTIVE_LICENSE_SET, copyManager, true, defaultUriPrefix);
        for (org.spdx.library.model.v2.license.AnyLicenseInfo fromMember:fromDisjunctiveLicenseSet.getMembers()) {
            toDisjunctiveLicenseSet.getMembers().add(convertAndStore(fromMember));
        }
        toDisjunctiveLicenseSet.setCreationInfo(defaultCreationInfo);
        return toDisjunctiveLicenseSet;
    }

    /**
     * Converts an SPDX spec version 2 SPDX ExtractedLicenseInfo to an SPDX spec version 3 SPDX CustomLicense and store the result
     * in the toStore
     * @param fromExtractedLicenseInfo an SPDX spec version 2 ExtractedLicenseInfo to convert from
     * @return an SPDX spec version 3 CustomLicense
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public CustomLicense convertAndStore(org.spdx.library.model.v2.license.ExtractedLicenseInfo fromExtractedLicenseInfo) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromExtractedLicenseInfo.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE);
        if (existing.isPresent()) {
            return (CustomLicense)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromExtractedLicenseInfo.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (CustomLicense)getExistingObject(fromExtractedLicenseInfo.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE).get();
        }
        CustomLicense toCustomLicense = (CustomLicense)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE, copyManager, true, defaultUriPrefix);
        toCustomLicense.setCreationInfo(defaultCreationInfo);
        toCustomLicense.setLicenseText(fromExtractedLicenseInfo.getExtractedText());
        toCustomLicense.setName(fromExtractedLicenseInfo.getName());
        toCustomLicense.getSeeAlsos().addAll(fromExtractedLicenseInfo.getSeeAlso());
        toCustomLicense.setComment(fromExtractedLicenseInfo.getComment());
        return toCustomLicense;
    }

    /**
     * Converts an SPDX spec version 2 SPDX OrLaterOperator to an SPDX spec version 3 SPDX OrLaterOperator and store the result
     * in the toStore
     * @param fromOrLaterOperator an SPDX spec version 2 OrLaterOperator to convert from
     * @return an SPDX spec version 3 OrLaterOperator
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public OrLaterOperator convertAndStore(org.spdx.library.model.v2.license.OrLaterOperator fromOrLaterOperator) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromOrLaterOperator.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_OR_LATER_OPERATOR);
        if (existing.isPresent()) {
            return (OrLaterOperator)existing.get();
        }
        String toObjectUri = defaultUriPrefix +  toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromOrLaterOperator.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (OrLaterOperator)getExistingObject(fromOrLaterOperator.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_OR_LATER_OPERATOR).get();
        }
        OrLaterOperator toOrLaterOperator = (OrLaterOperator)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_OR_LATER_OPERATOR, copyManager, true, defaultUriPrefix);
        toOrLaterOperator.setCreationInfo(defaultCreationInfo);
        toOrLaterOperator.setSubjectLicense((License)convertAndStore(fromOrLaterOperator.getLicense()));
        return toOrLaterOperator;
    }

    /**
     * Converts an SPDX spec version 2 SPDX SpdxListedLicense to an SPDX spec version 3 SPDX ListedLicense and store the result
     * in the toStore
     * @param fromSpdxListedLicense an SPDX spec version 2 SpdxListedLicense to convert from
     * @return an SPDX spec version 3 ListedLicense
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public ListedLicense convertAndStore(org.spdx.library.model.v2.license.SpdxListedLicense fromSpdxListedLicense) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromSpdxListedLicense.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE);
        if (existing.isPresent()) {
            return (ListedLicense)existing.get();
        }
        String existingUri = this.alreadyConverted.putIfAbsent(fromSpdxListedLicense.getObjectUri(), fromSpdxListedLicense.getObjectUri());
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (ListedLicense)getExistingObject(fromSpdxListedLicense.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE).get();
        }
        String licenseId = SpdxListedLicenseModelStore.objectUriToLicenseOrExceptionId(fromSpdxListedLicense.getObjectUri());
        if (ListedLicenses.getListedLicenses().isSpdxListedLicenseId(licenseId)) {
            ListedLicense retval = ListedLicenses.getListedLicenses().getListedLicenseById(licenseId);
            if (complexLicenses) {
                copyManager.copy(toModelStore, fromSpdxListedLicense.getObjectUri(), retval.getModelStore(),
                        fromSpdxListedLicense.getObjectUri(), toSpecVersion, null);
            }
            return retval;
        }
        ListedLicense toListedLicense = (ListedLicense)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                fromSpdxListedLicense.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE, copyManager, true, defaultUriPrefix);
        toListedLicense.setCreationInfo(defaultCreationInfo);
        toListedLicense.setComment(fromSpdxListedLicense.getComment());
        // fromSpdxListedLicense.getCrossRef()) - no equivalent in SPDX version 3.X
        toListedLicense.setDeprecatedVersion(fromSpdxListedLicense.getDeprecatedVersion());
        toListedLicense.setIsFsfLibre(fromSpdxListedLicense.getFsfLibre());
        // fromSpdxListedLicense.getLicenseHeaderHtml();  - no equivalent in SPDX version 3.X
        toListedLicense.setLicenseText(fromSpdxListedLicense.getLicenseText());
        // fromSpdxListedLicense.getLicenseTextHtml();   - no equivalent in SPDX version 3.X
        toListedLicense.setName(fromSpdxListedLicense.getName());
        toListedLicense.getSeeAlsos().addAll(fromSpdxListedLicense.getSeeAlso());
        toListedLicense.setStandardLicenseHeader(fromSpdxListedLicense.getStandardLicenseHeader());
        // fromSpdxListedLicense.getStandardLicenseHeaderTemplate(); - no equivalent in SPDX version 3.X
        toListedLicense.setStandardLicenseTemplate(fromSpdxListedLicense.getStandardLicenseTemplate());
        toListedLicense.setIsDeprecatedLicenseId(fromSpdxListedLicense.isDeprecated());
        toListedLicense.setIsOsiApproved(fromSpdxListedLicense.isOsiApproved());
        return toListedLicense;
    }

    /**
     * Converts an SPDX spec version 2 SPDX WithExceptionOperator to an SPDX spec version 3 SPDX WithAdditionOperator and store the result
     * in the toStore
     * @param fromWithExceptionOperator an SPDX spec version 2 WithExceptionOperator to convert from
     * @return an SPDX spec version 3 WithAdditionOperator
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public WithAdditionOperator convertAndStore(org.spdx.library.model.v2.license.WithExceptionOperator fromWithExceptionOperator) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromWithExceptionOperator.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_WITH_ADDITION_OPERATOR);
        if (existing.isPresent()) {
            return (WithAdditionOperator)existing.get();
        }
        String toObjectUri = defaultUriPrefix +  toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromWithExceptionOperator.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (WithAdditionOperator)getExistingObject(fromWithExceptionOperator.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_WITH_ADDITION_OPERATOR).get();
        }
        WithAdditionOperator toWithAdditionOperator = (WithAdditionOperator)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_WITH_ADDITION_OPERATOR, copyManager, true, defaultUriPrefix);
        toWithAdditionOperator.setCreationInfo(defaultCreationInfo);
        toWithAdditionOperator.setSubjectAddition(convertAndStore(fromWithExceptionOperator.getException()));
        toWithAdditionOperator.setSubjectExtendableLicense((ExtendableLicense)convertAndStore(fromWithExceptionOperator.getLicense()));
        return toWithAdditionOperator;
    }


    /**
     * Converts an SPDX spec version 2 SPDX LicenseException to an SPDX spec version 3 LicenseAddition  and store the result
     * in the toStore
     * @param fromException an SPDX spec version 2 LicenseException to convert from
     * @return an SPDX spec version 3 LicenseAddition
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public LicenseAddition convertAndStore(org.spdx.library.model.v2.license.LicenseException fromException) throws InvalidSPDXAnalysisException {
        if (fromException instanceof org.spdx.library.model.v2.license.ListedLicenseException) {
            return convertAndStore((org.spdx.library.model.v2.license.ListedLicenseException)fromException);
        }
        Optional<ModelObjectV3> existing = getExistingObject(fromException.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE_ADDITION);
        if (existing.isPresent()) {
            return (CustomLicenseAddition)existing.get();
        }
        String toObjectUri = defaultUriPrefix +  toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromException.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (CustomLicenseAddition)getExistingObject(fromException.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE_ADDITION).get();
        }
        CustomLicenseAddition toCustomAddition = (CustomLicenseAddition)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.EXPANDED_LICENSING_CUSTOM_LICENSE_ADDITION, copyManager, true, defaultUriPrefix);
        convertLicenseAdditionProperties(fromException, toCustomAddition);
        return toCustomAddition;
    }

    /**
     * Convert and add properties from the fromException to the toAddition
     * @param fromException SPDX spec version 2 LicenseException to copy properties from
     * @param toAddition SPDX spec version 3 LicenseAddition to copy properties to
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    private void convertLicenseAdditionProperties(
            org.spdx.library.model.v2.license.LicenseException fromException,
            LicenseAddition toAddition) throws InvalidSPDXAnalysisException {
        toAddition.setCreationInfo(defaultCreationInfo);
        toAddition.setAdditionText(fromException.getLicenseExceptionText());
        toAddition.setComment(fromException.getComment());
        toAddition.setName(fromException.getName());
        toAddition.setStandardAdditionTemplate(fromException.getLicenseExceptionTemplate());
        toAddition.getSeeAlsos().addAll(fromException.getSeeAlso());
    }

    /**
     * Converts an SPDX spec version 2 SPDX ListedLicenseException to an SPDX spec version 3 ListedLicenseAddition  and store the result
     * in the toStore
     * @param fromException an SPDX spec version 2 ListedLicenseException to convert from
     * @return an SPDX spec version 3 ListedLicenseException
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public ListedLicenseException convertAndStore(org.spdx.library.model.v2.license.ListedLicenseException fromException) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromException.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE_EXCEPTION);
        if (existing.isPresent()) {
            return (ListedLicenseException)existing.get();
        }
        String existingUri = this.alreadyConverted.putIfAbsent(fromException.getObjectUri(), fromException.getObjectUri());
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (ListedLicenseException)getExistingObject(fromException.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE_EXCEPTION).get();
        }
        String exceptionId = SpdxListedLicenseModelStore.objectUriToLicenseOrExceptionId(fromException.getObjectUri());
        if (ListedLicenses.getListedLicenses().isSpdxListedExceptionId(exceptionId)) {
            ListedLicenseException retval = ListedLicenses.getListedLicenses().getListedExceptionById(exceptionId);
            if (complexLicenses) {
                copyManager.copy(toModelStore, fromException.getObjectUri(), retval.getModelStore(),
                        fromException.getObjectUri(), toSpecVersion, null);
            }
            return retval;
        }
        ListedLicenseException toListedException = (ListedLicenseException)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                fromException.getObjectUri(), SpdxConstantsV3.EXPANDED_LICENSING_LISTED_LICENSE_EXCEPTION, copyManager, true, defaultUriPrefix);
        convertLicenseAdditionProperties(fromException, toListedException);
        toListedException.setDeprecatedVersion(fromException.getDeprecatedVersion());
        toListedException.setIsDeprecatedAdditionId(fromException.isDeprecated());
        // fromException.getExample(); - no SPDX spec version 3 equivalent
        // fromException.getExceptionTextHtml(); - no SPDX spec version 3 equivalent
        return toListedException;
    }

    /**
     * Converts an SPDX spec version 2 SPDX AnyLicenseInfo to an SPDX spec version 3 LicenseExpression
     * @param fromLicense an SPDX spec version 2 AnyLicenseInfo
     * @return an SPDX spec version 3 LicenseExpression
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public LicenseExpression convertToLicenseExpression(org.spdx.library.model.v2.license.AnyLicenseInfo fromLicense) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromLicense.getObjectUri(), SpdxConstantsV3.SIMPLE_LICENSING_LICENSE_EXPRESSION);
        if (existing.isPresent()) {
            return (LicenseExpression)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromLicense.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (LicenseExpression)getExistingObject(fromLicense.getObjectUri(), SpdxConstantsV3.SIMPLE_LICENSING_LICENSE_EXPRESSION).get();
        }
        LicenseExpression licenseExpression = (LicenseExpression)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.SIMPLE_LICENSING_LICENSE_EXPRESSION, copyManager, true, defaultUriPrefix);
        licenseExpression.setCreationInfo(defaultCreationInfo);
        String expression = fromLicense.toString();
        licenseExpression.setLicenseExpression(expression);
        StringTokenizer tokenizer = new StringTokenizer(expression, "() ");
        while (tokenizer.hasMoreTokens()) {
            String token = tokenizer.nextToken().trim();
            if (token.startsWith(SpdxConstantsCompatV2.NON_STD_LICENSE_ID_PRENUM)) {
                licenseExpression.getCustomIdToUris().add(licenseExpression.createDictionaryEntry(toModelStore.getNextId(IdType.Anonymous))
                        .setKey(token)
                        .setValue(defaultUriPrefix + token)
                        .build());
            }
        }
        return licenseExpression;
    }

    /**
     * Converts an SPDX spec version 2 SPDX AnyLicenseInfo to an SPDX spec version 3 SPDX AnyLicenseInfo and store the result
     * @param fromLicense an SPDX spec version 2 AnyLicenseInfo
     * @return an SPDX spec version 3 AnyLicenseInfo
     * @throws InvalidSPDXAnalysisException on any errors converting
     */
    public AnyLicenseInfo convertAndStore(org.spdx.library.model.v2.license.AnyLicenseInfo fromLicense) throws InvalidSPDXAnalysisException {
        if (!complexLicenses) {
            return convertToLicenseExpression(fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.ConjunctiveLicenseSet) {
            return convertAndStore((org.spdx.library.model.v2.license.ConjunctiveLicenseSet)fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.DisjunctiveLicenseSet) {
            return convertAndStore((org.spdx.library.model.v2.license.DisjunctiveLicenseSet)fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.ExternalExtractedLicenseInfo) {
            String externalUri = ((org.spdx.library.model.v2.license.ExternalExtractedLicenseInfo)fromLicense).getIndividualURI();
            logger.warn("Referencing an external SPDX 2 element with URI {} while converting from SPDX 2 to 3", externalUri);
            addExternalMapInfo(externalUri);
            return new ExternalCustomLicense(externalUri);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.ExtractedLicenseInfo) {
            return convertAndStore((org.spdx.library.model.v2.license.ExtractedLicenseInfo)fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.OrLaterOperator) {
            return convertAndStore((org.spdx.library.model.v2.license.OrLaterOperator)fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.SpdxListedLicense) {
            return convertAndStore((org.spdx.library.model.v2.license.SpdxListedLicense)fromLicense);
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.SpdxNoneLicense) {
            return new NoneLicense();
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.SpdxNoAssertionLicense) {
            return new NoAssertionLicense();
        } else if (fromLicense instanceof org.spdx.library.model.v2.license.WithExceptionOperator) {
            return convertAndStore((org.spdx.library.model.v2.license.WithExceptionOperator)fromLicense);
        } else {
            throw new InvalidSPDXAnalysisException("Can not convert the from AnyLicenseInfo type "+fromLicense.getType());
        }
    }

    /**
     * Converts the Element and stores all properties in the toStore
     * @param fromElement element to convert from
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public Element convertAndStore(org.spdx.library.model.v2.SpdxElement fromElement) throws InvalidSPDXAnalysisException {
        if (fromElement instanceof org.spdx.library.model.v2.SpdxFile) {
            return convertAndStore((org.spdx.library.model.v2.SpdxFile)fromElement);
        } else if (fromElement instanceof org.spdx.library.model.v2.SpdxPackage) {
            return convertAndStore((org.spdx.library.model.v2.SpdxPackage)fromElement);
        } else if (fromElement instanceof org.spdx.library.model.v2.SpdxSnippet) {
            return convertAndStore((org.spdx.library.model.v2.SpdxSnippet)fromElement);
        } else if (fromElement instanceof org.spdx.library.model.v2.SpdxNoAssertionElement) {
            return new NoAssertionElement();
        } else if (fromElement instanceof org.spdx.library.model.v2.SpdxNoneElement) {
            return new NoneElement();
        } else if (fromElement instanceof org.spdx.library.model.v2.ExternalSpdxElement) {
            String externalUri = ((org.spdx.library.model.v2.ExternalSpdxElement)fromElement).getIndividualURI();
            logger.warn("Referencing an external SPDX 2 element with URI {} while converting from SPDX 2 to 3", externalUri);
            addExternalMapInfo(externalUri);
            return new ExternalElement(externalUri);
        } else if (fromElement instanceof org.spdx.library.model.v2.SpdxDocument) {
            return convertAndStore((org.spdx.library.model.v2.SpdxDocument)fromElement);
        } else {
            throw new InvalidSPDXAnalysisException("Conversion of SPDX 2 type" + fromElement.getType()+" is not currently supported");
        }
    }

    /**
     * Creates ExternalMaps for a reference to an external SPDX element or license
     * @param externalUri URI of the external element
     * @throws InvalidSPDXAnalysisException on error creating ExternalMap
     */
    private void addExternalMapInfo(String externalUri) throws InvalidSPDXAnalysisException {
        Objects.requireNonNull(externalUri, "External URI can not be null");
        String[] parts = externalUri.split("#");
        if (parts.length != 2) {
            logger.warn("{} is not a valid SPDX Spec version 2 external referenced - should have a document uri + '#' + ID", externalUri);
            return;
        }
        Map<Collection<ExternalMap>, ExternalMapInfo> externalMapMap = docUriToExternalMap.get(parts[0]);
        if (Objects.isNull(externalMapMap)) {
            logger.warn("No corresponding ExternalDocumentRefs for {}", externalUri);
            return;
        }
        synchronized(externalMapMap) {
            for (ExternalMapInfo mapInfo:externalMapMap.values()) {
                mapInfo.addExternalMap(externalUri, toModelStore);
            }
        }

    }

    /**
     * Converts the SPDX 2 SpdxFile to an SPDX 3 SpdxFile and returns the converted file
     * @param spdxFile SPDX file to convert from
     * @return SPDX 3 SpdxFile
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public SpdxFile convertAndStore(org.spdx.library.model.v2.SpdxFile spdxFile) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(spdxFile.getObjectUri(), SpdxConstantsV3.SOFTWARE_SPDX_FILE);
        if (existing.isPresent()) {
            return (SpdxFile)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(spdxFile.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (SpdxFile)getExistingObject(spdxFile.getObjectUri(), SpdxConstantsV3.SOFTWARE_SPDX_FILE).get();
        }
        SpdxFile toFile = (SpdxFile)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.SOFTWARE_SPDX_FILE, copyManager, true, defaultUriPrefix);
        convertItemProperties(spdxFile, toFile);

        for (org.spdx.library.model.v2.Checksum checksum:spdxFile.getChecksums()) {
            toFile.getVerifiedUsings().add(convertAndStore(checksum));
        }
        // spdxFile.getFileContributors(); - No equivalent SPDX 3 property
        // spdxFile.getFileDependency(); - deprecated - skipping

        for (org.spdx.library.model.v2.enumerations.FileType fileType : spdxFile.getFileTypes()) {
            convertAndAddFileType(fileType, toFile);
        }
        Optional<String> noticeText = spdxFile.getNoticeText();

        noticeText.ifPresent(s -> toFile.getAttributionTexts().add(s));
        // - this is already captured in the checksums - String sha1 = spdxFile.getSha1();
        return toFile;
    }

    /**
     * Converts the SPDX spec version 2 Checksum to an SPDX spec version 3 Hash and store the result
     * @param checksum SPDX spec version 2 Checksum
     * @return SPDX spec version 3 Hash
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public Hash convertAndStore(org.spdx.library.model.v2.Checksum checksum) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(checksum.getObjectUri(), SpdxConstantsV3.CORE_HASH);
        if (existing.isPresent()) {
            return (Hash)existing.get();
        }
        String toObjectUri = toModelStore.getNextId(IdType.Anonymous);
        String existingUri = this.alreadyConverted.putIfAbsent(checksum.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (Hash)getExistingObject(checksum.getObjectUri(), SpdxConstantsV3.CORE_HASH).get();
        }
        Hash toHash = (Hash)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.CORE_HASH, copyManager, true, defaultUriPrefix);
        toHash.setAlgorithm(HASH_ALGORITH_MAP.get(checksum.getAlgorithm()));
        toHash.setHashValue(checksum.getValue());
        return toHash;
    }

    /**
     * Converts an SPDX spec version 2 FileType to the corresponding SPDX model 3 software purpose and/or content type
     * and adds that information to the file
     * @param fileType SPDX spec version 2 FileType to convert and add
     * @param file SPDX spec version 3 SpdxFile to add the software purpose or content type to
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    private void convertAndAddFileType(org.spdx.library.model.v2.enumerations.FileType fileType, SpdxFile file) throws InvalidSPDXAnalysisException {
        switch (fileType) {
            case ARCHIVE: addSoftwarePurpose(SoftwarePurpose.ARCHIVE, file); break;
            case BINARY: file.setContentType("application/octet-stream"); break;
            case SOURCE: addSoftwarePurpose(SoftwarePurpose.SOURCE, file); break;
            case TEXT: file.setContentType("text/plain"); break;
            case APPLICATION: addSoftwarePurpose(SoftwarePurpose.APPLICATION, file); break;
            case AUDIO: file.setContentType("audio/*"); break;
            case IMAGE: file.setContentType("image/*"); break;
            case VIDEO: file.setContentType("video/*"); break;
            case DOCUMENTATION: addSoftwarePurpose(SoftwarePurpose.DOCUMENTATION, file); break;
            case SPDX: file.setContentType("text/spdx"); break;
            case OTHER: addSoftwarePurpose(SoftwarePurpose.OTHER, file); break;

            default: throw new InvalidSPDXAnalysisException("Unknown file type "+fileType);
        }
    }

    /**
     * Adds a Software Purpose to a SoftwareArtifact.  If the primaryPurpose is already used, add as an additionalPurpose
     * @param purpose purpose to add
     * @param artifact artifact to add the purpose to
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    private void addSoftwarePurpose(SoftwarePurpose purpose,
                                    SoftwareArtifact artifact) throws InvalidSPDXAnalysisException {
        if (artifact.getPrimaryPurpose().isPresent()) {
            artifact.getAdditionalPurposes().add(purpose);
        } else {
            artifact.setPrimaryPurpose(purpose);
        }
    }

    /**
     * Converts and copies properties from the fromItem to the toArtifact
     * @param fromItem SPDX spec version 2 Item to copy properties from
     * @param toArtifact SPDX spec version 3 SoftwareArtifact to copy properties to
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    private void convertItemProperties(org.spdx.library.model.v2.SpdxItem fromItem, SoftwareArtifact toArtifact) throws InvalidSPDXAnalysisException {
        convertElementProperties(fromItem, toArtifact);
        toArtifact.getAttributionTexts().addAll(fromItem.getAttributionText());
        toArtifact.setCopyrightText(fromItem.getCopyrightText());
        Optional<String> licenseComments = fromItem.getLicenseComments();
        if (licenseComments.isPresent()) {
            Optional<String> existingComment = toArtifact.getComment();
            toArtifact.setComment(existingComment.map(s -> s + ";" + licenseComments.get()).orElseGet(licenseComments::get));
        }
        org.spdx.library.model.v2.license.AnyLicenseInfo concludedLicense = fromItem.getLicenseConcluded();
        if (Objects.nonNull(concludedLicense)) {
            Relationship concludedRelationship = (Relationship)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                    defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId), SpdxConstantsV3.CORE_RELATIONSHIP, copyManager, true, defaultUriPrefix);
            concludedRelationship.setCreationInfo(defaultCreationInfo);
            concludedRelationship.setFrom(toArtifact);
            concludedRelationship.getTos().add(convertAndStore(concludedLicense));
            concludedRelationship.setRelationshipType(RelationshipType.HAS_CONCLUDED_LICENSE);
        }
        if (!(fromItem instanceof org.spdx.library.model.v2.SpdxPackage)) {
            // we use the license concluded for the SPDX package
            for (org.spdx.library.model.v2.license.AnyLicenseInfo declaredLicense:fromItem.getLicenseInfoFromFiles()) {
                Relationship declaredRelationship = (Relationship)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                        defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId), SpdxConstantsV3.CORE_RELATIONSHIP, copyManager, true, defaultUriPrefix);
                declaredRelationship.setCreationInfo(defaultCreationInfo);
                declaredRelationship.setFrom(toArtifact);
                declaredRelationship.getTos().add(convertAndStore(declaredLicense));
                declaredRelationship.setRelationshipType(RelationshipType.HAS_DECLARED_LICENSE);
            }
        }
    }

    /**
     * Converts the SPDX 2 SpdxPackage to an SPDX 3 SpdxPackage and returns the converted package
     * @param spdxPackage SPDX package to convert from
     * @return SPDX 3 SpdxPackage
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public SpdxPackage convertAndStore(org.spdx.library.model.v2.SpdxPackage spdxPackage) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(spdxPackage.getObjectUri(), SpdxConstantsV3.SOFTWARE_SPDX_PACKAGE);
        if (existing.isPresent()) {
            return (SpdxPackage)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(spdxPackage.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (SpdxPackage)getExistingObject(spdxPackage.getObjectUri(), SpdxConstantsV3.SOFTWARE_SPDX_PACKAGE).get();
        }
        SpdxPackage toPackage = (SpdxPackage)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.SOFTWARE_SPDX_PACKAGE, copyManager, true, defaultUriPrefix);
        convertItemProperties(spdxPackage, toPackage);
        toPackage.setBuiltTime(spdxPackage.getBuiltDate().orElse(null));
        toPackage.setDescription(spdxPackage.getDescription().orElse(null));
        toPackage.setDownloadLocation(spdxPackage.getDownloadLocation().orElse(null));
        for (org.spdx.library.model.v2.ExternalRef externalRef:spdxPackage.getExternalRefs()) {
            addExternalRefToArtifact(externalRef, toPackage);
        }
        // spdxPackage.getFiles() - these should be captured in relationships
        toPackage.setHomePage(spdxPackage.getHomepage().orElse(null));
        Optional<String> originator = spdxPackage.getOriginator();
        if (originator.isPresent()) {
            toPackage.getOriginatedBys().add(stringToAgent(originator.get(), toPackage.getCreationInfo()));
        }
        Optional<String> packageFileName = spdxPackage.getPackageFileName();
        if (packageFileName.isPresent()) {
            addPackageFileNameToPackage(packageFileName.get(), toPackage, spdxPackage.getChecksums());
        }
        Optional<org.spdx.library.model.v2.SpdxPackageVerificationCode> pkgVerificationCode = spdxPackage.getPackageVerificationCode();
        if (pkgVerificationCode.isPresent()) {
            toPackage.getVerifiedUsings().add(convertAndStore(pkgVerificationCode.get()));
        }
        Optional<org.spdx.library.model.v2.enumerations.Purpose> primaryPurpose = spdxPackage.getPrimaryPurpose();
        if (primaryPurpose.isPresent()) {
            if (toPackage.getPrimaryPurpose().isPresent()) {
                toPackage.getAdditionalPurposes().add(toPackage.getPrimaryPurpose().get());
            }
            toPackage.setPrimaryPurpose(PURPOSE_MAP.get(primaryPurpose.get()));
        }
        toPackage.setReleaseTime(spdxPackage.getReleaseDate().orElse(null));
        toPackage.setSourceInfo(spdxPackage.getSourceInfo().orElse(null));
        toPackage.setSummary(spdxPackage.getSummary().orElse(null));
        Optional<String> supplier = spdxPackage.getSupplier();
        if (supplier.isPresent()) {
            toPackage.setSuppliedBy(stringToAgent(supplier.get(), toPackage.getCreationInfo()));
        }
        toPackage.setValidUntilTime(spdxPackage.getValidUntilDate().orElse(null));
        toPackage.setPackageVersion(spdxPackage.getVersionInfo().orElse(null));

        org.spdx.library.model.v2.license.AnyLicenseInfo declaredLicense = spdxPackage.getLicenseDeclared();
        if (Objects.nonNull(declaredLicense)) {
            Relationship declaredRelationship = (Relationship)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                    defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId), SpdxConstantsV3.CORE_RELATIONSHIP, copyManager, true, defaultUriPrefix);
            declaredRelationship.setCreationInfo(defaultCreationInfo);
            declaredRelationship.setFrom(toPackage);
            declaredRelationship.getTos().add(convertAndStore(declaredLicense));
            declaredRelationship.setRelationshipType(RelationshipType.HAS_DECLARED_LICENSE);
        }
        return toPackage;
    }

    /**
     * Converts the spdxPackageVerificationCode to an IntegrityMethod and store the result in the toModelStore
     * @param spdxPackageVerificationCode SPDX Spec version 2 package verification code
     * @return the package verification code integrity method
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public IntegrityMethod convertAndStore(
            org.spdx.library.model.v2.SpdxPackageVerificationCode spdxPackageVerificationCode) throws InvalidSPDXAnalysisException {
        PackageVerificationCode pkgVerificationCode = (PackageVerificationCode)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toModelStore.getNextId(IdType.Anonymous), SpdxConstantsV3.CORE_PACKAGE_VERIFICATION_CODE,
                copyManager, true, defaultUriPrefix);

        pkgVerificationCode.setAlgorithm(HashAlgorithm.SHA1);
        pkgVerificationCode.setHashValue(spdxPackageVerificationCode.getValue());
        pkgVerificationCode.getPackageVerificationCodeExcludedFiles().addAll(spdxPackageVerificationCode.getExcludedFileNames());
        return pkgVerificationCode;
    }

    /**
     * Create a File artifact and add that to toPackage as a relationship
     * @param fileName Name of the File artifact
     * @param toPackage package to add the file to
     * @param fileChecksums checksums for the file
     * @throws InvalidSPDXAnalysisException on SPDX parsing errors
     */
    private void addPackageFileNameToPackage(String fileName,
                                             SpdxPackage toPackage, Collection<org.spdx.library.model.v2.Checksum> fileChecksums) throws InvalidSPDXAnalysisException {
        SpdxFile file = toPackage.createSpdxFile(defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId))
                .setName(fileName)
                .build();
        for (org.spdx.library.model.v2.Checksum checksum : fileChecksums) {
            file.getVerifiedUsings().add(convertAndStore(checksum));
        }
        toPackage.createRelationship(defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId))
                .setRelationshipType(RelationshipType.HAS_DISTRIBUTION_ARTIFACT)
                .setFrom(toPackage)
                .addTo(file)
                .setCompleteness(RelationshipCompleteness.COMPLETE)
                .build();
    }

    /**
     * @param externalRef SPDX Spec version 2 External Ref to add to the package
     * @param artifact SPDX Spec version 3 Artifact to add either an ExternalRef or ExternalId depending on the externalRef type
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    private void addExternalRefToArtifact(org.spdx.library.model.v2.ExternalRef externalRef,
                                          SoftwareArtifact artifact) throws InvalidSPDXAnalysisException {
        addExternalRefToArtifact(externalRef, artifact, toModelStore);
    }

    /**
     * @param externalRef SPDX Spec version 2 External Ref to add to the package
     * @param artifact SPDX Spec version 3 Artifact to add either an ExternalRef or ExternalId depending on the externalRef type
     * @param modelStore modelStore to use for creating any SPDX objects
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public static void addExternalRefToArtifact(org.spdx.library.model.v2.ExternalRef externalRef,
                                                SoftwareArtifact artifact, IModelStore modelStore) throws InvalidSPDXAnalysisException {
        addExternalRefToArtifact(externalRef.getReferenceCategory(), externalRef.getReferenceType(),
                externalRef.getReferenceLocator(), externalRef.getComment().orElse(null), artifact, modelStore);
    }

    /**
     * @param referenceCategory Reference category for external ref
     * @param referenceType Reference type for external ref
     * @param referenceLocator Locator for external ref
     * @param comment External reference comment
     * @param artifact Artifact which contains the external ref
     * @param modelStore modelStore to use for creating any SPDX objects
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public static void addExternalRefToArtifact(org.spdx.library.model.v2.enumerations.ReferenceCategory referenceCategory,
                                                org.spdx.library.model.v2.ReferenceType referenceType, String referenceLocator, @Nullable String comment,
                                                SoftwareArtifact artifact, IModelStore modelStore) throws InvalidSPDXAnalysisException {
        Objects.requireNonNull(referenceType);
        switch (referenceType.getIndividualURI()) {
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "cpe22Type":
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "cpe23Type":
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "swid":
                artifact.getExternalIdentifiers().add(artifact.createExternalIdentifier(modelStore.getNextId(IdType.Anonymous))
                        .setExternalIdentifierType(EXTERNAL_IDENTIFIER_TYPE_MAP.get(referenceType.getIndividualURI()))
                        .setIdentifier(referenceLocator)
                        .setComment(comment)
                        .build()); break;
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "purl": {
                if (artifact instanceof SpdxPackage) {
                    ((SpdxPackage)artifact).setPackageUrl(referenceLocator);
                } else {
                    artifact.getExternalIdentifiers().add(artifact.createExternalIdentifier(modelStore.getNextId(IdType.Anonymous))
                            .setExternalIdentifierType(EXTERNAL_IDENTIFIER_TYPE_MAP.get(referenceType.getIndividualURI()))
                            .setIdentifier(referenceLocator)
                            .setComment(comment)
                            .build()); break;
                }
            } break;
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "swh":
            case SpdxConstantsCompatV2.SPDX_LISTED_REFERENCE_TYPES_PREFIX + "gitoid":
                artifact.getContentIdentifiers().add(artifact.createContentIdentifier(modelStore.getNextId(IdType.Anonymous))
                        .setContentIdentifierType(CONTENT_IDENTIFIER_TYPE_MAP.get(referenceType.getIndividualURI()))
                        .setContentIdentifierValue(referenceLocator)
                        .setComment(comment)
                        .build()); break;
            default: {
                ExternalRefType externalRefType = EXTERNAL_REF_TYPE_MAP.get(referenceType.getIndividualURI());
                if (Objects.isNull(externalRefType)) {
                    switch (referenceCategory) {
                        case PACKAGE_MANAGER: externalRefType = ExternalRefType.BUILD_SYSTEM; break;
                        case SECURITY: externalRefType = ExternalRefType.SECURITY_OTHER; break;
                        default: externalRefType = ExternalRefType.OTHER;
                    }
                }
                artifact.getExternalRefs().add(artifact.createExternalRef(modelStore.getNextId(IdType.Anonymous))
                        .setExternalRefType(externalRefType)
                        .addLocator(referenceLocator)
                        .setComment(comment)
                        .build());
            }
        }
    }

    /**
     * Converts the SPDX 2 SpdxSnippet to an SPDX 3 Snippet and returns the converted snippet
     * @param fromSnippet SPDX 2 snippet to convert from
     * @return SPDX 3 Snippet
     * @throws InvalidSPDXAnalysisException on any error in conversion
     */
    public Snippet convertAndStore(org.spdx.library.model.v2.SpdxSnippet fromSnippet) throws InvalidSPDXAnalysisException {
        Optional<ModelObjectV3> existing = getExistingObject(fromSnippet.getObjectUri(), SpdxConstantsV3.SOFTWARE_SNIPPET);
        if (existing.isPresent()) {
            return (Snippet)existing.get();
        }
        String toObjectUri = defaultUriPrefix + toModelStore.getNextId(IdType.SpdxId);
        String existingUri = this.alreadyConverted.putIfAbsent(fromSnippet.getObjectUri(), toObjectUri);
        if (Objects.nonNull(existingUri)) {
            // small window if conversion occurred since the last check already converted
            return (Snippet)getExistingObject(fromSnippet.getObjectUri(), SpdxConstantsV3.SOFTWARE_SNIPPET).get();
        }
        Snippet toSnippet = (Snippet)SpdxModelClassFactoryV3.getModelObject(toModelStore,
                toObjectUri, SpdxConstantsV3.SOFTWARE_SNIPPET, copyManager, true, defaultUriPrefix);
        convertItemProperties(fromSnippet, toSnippet);
        StartEndPointer fromByteRange = fromSnippet.getByteRange();
        if (Objects.nonNull(fromByteRange)) {
            // noinspection DataFlowIssue
            ByteOffsetPointer startPointer = (ByteOffsetPointer) fromByteRange.getStartPointer();
            ByteOffsetPointer endPointer = (ByteOffsetPointer) fromByteRange.getEndPointer();
            if (Objects.nonNull(startPointer) && Objects.nonNull(endPointer)) {
                toSnippet.setByteRange(toSnippet
                        .createPositiveIntegerRange(toModelStore.getNextId(IdType.Anonymous))
                        .setBeginIntegerRange(startPointer.getOffset())
                        .setEndIntegerRange(endPointer.getOffset()).build());
            }
        }
        Optional<StartEndPointer> fromLineRange = fromSnippet.getLineRange();
        if (fromLineRange.isPresent()) {
            // noinspection DataFlowIssue
            LineCharPointer startPointer = (LineCharPointer) fromLineRange.get().getStartPointer();
            LineCharPointer endPointer = (LineCharPointer) fromLineRange.get().getEndPointer();
            if (Objects.nonNull(startPointer) && Objects.nonNull(endPointer)) {
                toSnippet.setLineRange(toSnippet
                        .createPositiveIntegerRange(toModelStore.getNextId(IdType.Anonymous))
                        .setBeginIntegerRange(startPointer.getLineNumber())
                        .setEndIntegerRange(endPointer.getLineNumber()).build());
            }
        }
        toSnippet.setSnippetFromFile(convertAndStore(Objects.requireNonNull(fromSnippet.getSnippetFromFile())));
        return toSnippet;
    }
}