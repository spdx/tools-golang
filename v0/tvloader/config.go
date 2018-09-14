// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package tvloader

// Config contains data for configuring the tag-value reader.
type Config struct {
	// FIXME NOT YET IMPLEMENTED
	// StopOnError determines whether the reader and parser will attempt to
	// continue after encountering an error. If set to false, it will continue
	// unless it encounters an unrecoverable error.
	StopOnError bool

	// FIXME add config for whether repeated tags (with cardinality one)
	// should be treated as errors, or should silently overwrite
}
