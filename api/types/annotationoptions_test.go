// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package types_test

import (
	"testing"

	. "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

func TestDefaultAnnotationConflictPolicy(t *testing.T) {
	if got := DefaultAnnotationConflictPolicy(); got != AnnotationConflictOverwrite {
		t.Fatalf("default annotation conflict policy = %q, want %q", got, AnnotationConflictOverwrite)
	}
}

func TestAnnotationOptionsYAML(t *testing.T) {
	var k Kustomization
	err := yaml.Unmarshal([]byte("commonAnnotationsOptions:\n  conflictPolicy: preserve\n"), &k)
	if err != nil {
		t.Fatal(err)
	}
	if got := k.CommonAnnotationsOptions.ConflictPolicy; got != AnnotationConflictPreserve {
		t.Fatalf("conflict policy = %q, want %q", got, AnnotationConflictPreserve)
	}
}
