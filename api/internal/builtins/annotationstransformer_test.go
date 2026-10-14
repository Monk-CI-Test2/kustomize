// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package builtins

import (
	"testing"

	"sigs.k8s.io/kustomize/api/types"
)

func TestAnnotationsTransformerEmptyOptionsUsesDefaults(t *testing.T) {
	p := &AnnotationsTransformerPlugin{Options: &types.AnnotationOptions{}}
	if err := p.Transform(nil); err != nil {
		t.Fatalf("empty annotation options should use defaults: %v", err)
	}
}
