// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package annotations

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/filters/filtersutil"
	"sigs.k8s.io/kustomize/api/filters/fsslice"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type annoMap map[string]string

type Filter struct {
	// Annotations is the set of annotations to apply to the inputs
	Annotations annoMap `yaml:"annotations,omitempty"`

	// FsSlice contains the FieldSpecs to locate the namespace field
	FsSlice types.FsSlice

	// ConflictPolicy controls handling of annotations already present.
	ConflictPolicy types.AnnotationConflictPolicy

	trackableSetter filtersutil.TrackableSetter
}

var _ kio.Filter = Filter{}
var _ kio.TrackableFilter = &Filter{}

// WithMutationTracker registers a callback which will be invoked each time a field is mutated
func (f *Filter) WithMutationTracker(callback func(key, value, tag string, node *yaml.RNode)) {
	f.trackableSetter.WithMutationTracker(callback)
}

func (f Filter) Filter(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
	keys := yaml.SortedMapKeys(f.Annotations)
	_, err := kio.FilterAll(yaml.FilterFunc(
		func(node *yaml.RNode) (*yaml.RNode, error) {
			for _, k := range keys {
				setValue := f.trackableSetter.SetEntry(k, f.Annotations[k], yaml.NodeTagString)
				switch f.ConflictPolicy {
				case types.AnnotationConflictPreserve:
					setValue = f.trackableSetter.SetEntryIfEmpty(k, f.Annotations[k], yaml.NodeTagString)
				case types.AnnotationConflictError:
					setValue = rejectAnnotationConflict(k, f.Annotations[k])
				}
				if err := node.PipeE(fsslice.Filter{
					FsSlice:    f.FsSlice,
					SetValue:   setValue,
					CreateKind: yaml.MappingNode, // Annotations are MappingNodes.
					CreateTag:  yaml.NodeTagMap,
				}); err != nil {
					return nil, err
				}
			}
			return node, nil
		})).Filter(nodes)
	return nodes, err
}

func rejectAnnotationConflict(key, value string) filtersutil.SetFn {
	return func(node *yaml.RNode) error {
		if existing := node.Field(key); !existing.IsNilOrEmpty() {
			return fmt.Errorf("annotation %q already exists", key)
		}
		return filtersutil.SetEntry(key, value, yaml.NodeTagString)(node)
	}
}
