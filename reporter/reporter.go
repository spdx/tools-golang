// Package reporter contains functions to generate a basic license count
// report from an in-memory SPDX Package section whose Files have been
// analyzed.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package reporter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/spdx/tools-golang/spdx"
)

// Generate takes a Package whose Files have been analyzed and an
// io.Writer, and outputs to the io.Writer a tabulated count of
// the number of Files for each unique LicenseConcluded in the set.
func Generate(pkg *spdx.Package, w io.Writer) error {
	if pkg.FilesAnalyzed == false {
		return fmt.Errorf("Package FilesAnalyzed is false")
	}
	totalFound, totalNotFound, foundCounts := countLicenses(pkg)

	wr := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintf(wr, "%d\t  License found\n", totalFound)
	fmt.Fprintf(wr, "%d\t  License not found\n", totalNotFound)
	fmt.Fprintf(wr, "%d\t  TOTAL\n", totalFound+totalNotFound)
	fmt.Fprintf(wr, "\n")

	counts := []struct {
		lic   string
		count int
	}{}
	for k, v := range foundCounts {
		var entry struct {
			lic   string
			count int
		}
		entry.lic = k
		entry.count = v
		counts = append(counts, entry)
	}

	sort.Slice(counts, func(i, j int) bool { return counts[i].count > counts[j].count })

	for _, c := range counts {
		fmt.Fprintf(wr, "%d\t  %s\n", c.count, c.lic)
	}
	fmt.Fprintf(wr, "%d\t  TOTAL FOUND\n", totalFound)

	wr.Flush()
	return nil
}

func countLicenses(pkg *spdx.Package) (int, int, map[string]int) {
	if pkg == nil || pkg.Files == nil {
		return 0, 0, nil
	}

	totalFound := 0
	totalNotFound := 0
	foundCounts := map[string]int{}
	for _, f := range pkg.Files {
		if f.LicenseConcluded == "" || f.LicenseConcluded == "NOASSERTION" {
			totalNotFound++
		} else {
			totalFound++
			foundCounts[f.LicenseConcluded]++
		}
	}

	return totalFound, totalNotFound, foundCounts
}
