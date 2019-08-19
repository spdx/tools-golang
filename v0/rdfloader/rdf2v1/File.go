package rdf2v1

import (
	"github.com/deltamobile/goraptor"
)

type File struct {
	FileName                  ValueStr
	FileSPDXIdentifier        ValueStr
	FileChecksum              *Checksum
	LicenseInfoInFile         []ValueStr
	FileCopyrightText         ValueStr
	ExtractedLicensingInfo    *ExtractedLicensingInfo
	DisjunctiveLicenseSet     *DisjunctiveLicenseSet
	ConjunctiveLicenseSet     *ConjunctiveLicenseSet
	FileContributor           []ValueStr
	FileComment               ValueStr
	FileLicenseComments       ValueStr
	FileType                  []ValueStr
	FileNoticeText            ValueStr
	Annotation                []*Annotation
	Project                   []*Project
	SnippetLicense            *License
	FileDependency            []*File
	FileRelationship          *Relationship
	FileLicenseSPDXIdentifier ValueStr
}
type Project struct {
	HomePage ValueStr
	Name     ValueStr
	URI      ValueStr
}

func (p *Parser) requestFile(node goraptor.Term) (*File, error) {
	obj, err := p.requestElementType(node, TypeFile)
	if err != nil {
		return nil, err
	}
	return obj.(*File), err
}
func (p *Parser) requestFileChecksum(node goraptor.Term) (*Checksum, error) {
	obj, err := p.requestElementType(node, TypeChecksum)
	if err != nil {
		return nil, err
	}
	return obj.(*Checksum), err
}

func (p *Parser) requestProject(node goraptor.Term) (*Project, error) {
	obj, err := p.requestElementType(node, TypeProject)
	if err != nil {
		return nil, err
	}
	return obj.(*Project), err
}
func (p *Parser) MapFile(file *File) *builder {
	builder := &builder{t: TypeFile, ptr: file}
	file.FileSPDXIdentifier = SPDXIDFile
	file.FileLicenseSPDXIdentifier = SPDXIDLicense
	builder.updaters = map[string]updater{
		"fileName": update(&file.FileName),
		"checksum": func(obj goraptor.Term) error {
			cksum, err := p.requestChecksum(obj)
			file.FileChecksum = cksum
			return err
		},
		"fileType": updateList(&file.FileType),
		"licenseConcluded": func(obj goraptor.Term) error {
			lic, err := p.requestLicense(obj)
			file.SnippetLicense = lic
			if err != nil {
				dls, err := p.requestDisjunctiveLicenseSet(obj)
				file.DisjunctiveLicenseSet = dls
				if err != nil {
					eli, err := p.requestExtractedLicensingInfo(obj)
					file.ExtractedLicensingInfo = eli
					return err
				}
			}
			return nil
		},
		"licenseInfoInFile": updateList(&file.LicenseInfoInFile),
		"copyrightText":     update(&file.FileCopyrightText),
		"licenseComments":   update(&file.FileLicenseComments),
		"rdfs:comment":      update(&file.FileComment),
		"noticeText":        update(&file.FileNoticeText),
		"fileContributor":   updateList(&file.FileContributor),
		"annotation": func(obj goraptor.Term) error {
			an, err := p.requestAnnotation(obj)
			file.Annotation = append(file.Annotation, an)
			return err
		},
		"artifactOf": func(obj goraptor.Term) error {
			pro, err := p.requestProject(obj)
			if err != nil {
				return err
			}
			file.Project = append(file.Project, pro)
			return err
		},
		"fileDependency": func(obj goraptor.Term) error {
			f, err := p.requestFile(obj)
			file.FileDependency = append(file.FileDependency, f)
			return err
		},
		"relationship": func(obj goraptor.Term) error {
			rel, err := p.requestRelationship(obj)
			file.FileRelationship = rel
			return err
		},
	}
	return builder
}

func (p *Parser) MapProject(pro *Project) *builder {
	builder := &builder{t: TypeProject, ptr: pro}
	pro.URI = ProjectURI
	builder.updaters = map[string]updater{
		"doap:homepage": update(&pro.HomePage),
		"doap:name":     update(&pro.Name),
	}
	return builder
}
