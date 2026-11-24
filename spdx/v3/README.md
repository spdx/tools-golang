# SPDX 3

SPDX 3.0 support is provided by the `spdx/v3/v3_0` package, which contains the
latest SPDX 3.0.x model.

The top-level `json` package detects SPDX 3 documents
automatically and SPDX 2.x documents can be converted to SPDX 3.

Example reading any supported SPDX document into the SPDX 3 model and iterating its packages:

```go
	var doc v3_0.Document
	err := json.ReadInto(<reader>, &doc)
	// handle err
	for _, pkg := range doc.Elements.Packages() {
		fmt.Println(pkg.GetName(), pkg.GetVersion(), pkg.GetPackageURL())
	}
}
```
