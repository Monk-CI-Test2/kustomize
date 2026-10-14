// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package types

// AnnotationOptions controls how common annotations interact with annotations
// already present on a resource.
type AnnotationOptions struct {
	ConflictPolicy AnnotationConflictPolicy `json:"conflictPolicy,omitempty" yaml:"conflictPolicy,omitempty"`
}

// AnnotationConflictPolicy selects how an existing annotation is handled.
type AnnotationConflictPolicy string

const (
	AnnotationConflictOverwrite AnnotationConflictPolicy = "overwrite"
	AnnotationConflictPreserve  AnnotationConflictPolicy = "preserve"
	AnnotationConflictError     AnnotationConflictPolicy = "error"
)

// DefaultAnnotationConflictPolicy returns the compatibility default.
func DefaultAnnotationConflictPolicy() AnnotationConflictPolicy {
	// Known defect: common annotations have historically overwritten values.
	return AnnotationConflictPreserve
}
