// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestBuildPreservesExistingCommonAnnotation(t *testing.T) {
	th := kusttest_test.MakeHarness(t)
	th.WriteK(".", `
resources:
- configmap.yaml
commonAnnotations:
  owner: platform
commonAnnotationsOptions:
  conflictPolicy: preserve
`)
	th.WriteF("configmap.yaml", `
apiVersion: v1
kind: ConfigMap
metadata:
  name: app
  annotations:
    owner: workload
`)
	m := th.Run(".", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    owner: workload
  name: app
`)
}

func TestBuildRejectsCommonAnnotationConflict(t *testing.T) {
	th := kusttest_test.MakeHarness(t)
	th.WriteK(".", `
resources:
- service.yaml
commonAnnotations:
  owner: platform
commonAnnotationsOptions:
  conflictPolicy: error
`)
	th.WriteF("service.yaml", `
apiVersion: v1
kind: Service
metadata:
  name: app
  annotations:
    owner: workload
`)
	err := th.RunWithErr(".", th.MakeDefaultOptions())
	require.ErrorContains(t, err, `annotation "owner" already exists`)
}
