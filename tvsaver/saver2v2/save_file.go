// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package saver2v2

import (
	"fmt"
	"io"
	"sort"

	"github.com/spdx/tools-golang/spdx"
)

func renderFile2_2(f *spdx.File2_2, w io.Writer) error {
	if f.Name != "" {
		fmt.Fprintf(w, "FileName: %s\n", f.Name)
	}
	if f.SPDXIdentifier != "" {
		fmt.Fprintf(w, "SPDXID: %s\n", spdx.RenderElementID(f.SPDXIdentifier))
	}
	for _, s := range f.Type {
		fmt.Fprintf(w, "FileType: %s\n", s)
	}
	if f.ChecksumSHA1 != "" {
		fmt.Fprintf(w, "FileChecksum: SHA1: %s\n", f.ChecksumSHA1)
	}
	if f.ChecksumSHA256 != "" {
		fmt.Fprintf(w, "FileChecksum: SHA256: %s\n", f.ChecksumSHA256)
	}
	if f.ChecksumMD5 != "" {
		fmt.Fprintf(w, "FileChecksum: MD5: %s\n", f.ChecksumMD5)
	}
	if f.LicenseConcluded != "" {
		fmt.Fprintf(w, "LicenseConcluded: %s\n", f.LicenseConcluded)
	}
	for _, s := range f.LicenseInfoInFile {
		fmt.Fprintf(w, "LicenseInfoInFile: %s\n", s)
	}
	if f.LicenseComments != "" {
		fmt.Fprintf(w, "LicenseComments: %s\n", f.LicenseComments)
	}
	if f.CopyrightText != "" {
		fmt.Fprintf(w, "FileCopyrightText: %s\n", textify(f.CopyrightText))
	}
	for _, aop := range f.ArtifactOfProjects {
		fmt.Fprintf(w, "ArtifactOfProjectName: %s\n", aop.Name)
		if aop.HomePage != "" {
			fmt.Fprintf(w, "ArtifactOfProjectHomePage: %s\n", aop.HomePage)
		}
		if aop.URI != "" {
			fmt.Fprintf(w, "ArtifactOfProjectURI: %s\n", aop.URI)
		}
	}
	if f.Comment != "" {
		fmt.Fprintf(w, "FileComment: %s\n", f.Comment)
	}
	if f.Notice != "" {
		fmt.Fprintf(w, "FileNotice: %s\n", f.Notice)
	}
	for _, s := range f.Contributor {
		fmt.Fprintf(w, "FileContributor: %s\n", s)
	}
	for _, s := range f.AttributionTexts {
		fmt.Fprintf(w, "FileAttributionText: %s\n", textify(s))
	}
	for _, s := range f.Dependencies {
		fmt.Fprintf(w, "FileDependency: %s\n", s)
	}

	fmt.Fprintf(w, "\n")

	// also render any snippets for this file
	// get slice of Snippet identifiers so we can sort them
	snippetKeys := []string{}
	for k := range f.Snippets {
		snippetKeys = append(snippetKeys, string(k))
	}
	sort.Strings(snippetKeys)
	for _, sID := range snippetKeys {
		s := f.Snippets[spdx.ElementID(sID)]
		renderSnippet2_2(s, w)
	}

	return nil
}
