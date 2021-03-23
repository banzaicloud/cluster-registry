package v1alpha1

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/banzaicloud/operator-tools/pkg/resources"
)

func (m SyncRuleMatch) Match(obj runtime.Object) (bool, error) {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}

	// objectkey
	if m.ObjectKey.Name != "" && m.ObjectKey.Name != objMeta.GetName() {
		return false, nil
	}
	if m.ObjectKey.Namespace != "" && m.ObjectKey.Namespace != objMeta.GetNamespace() {
		return false, nil
	}

	// namespace
	if len(m.Namespaces) > 0 && objMeta.GetNamespace() != "" {
		found := false
		for _, name := range m.Namespaces {
			if objMeta.GetNamespace() == name {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}

	// labels
	if len(m.Labels) > 0 {
		found := false
		for _, labelSelector := range m.Labels {
			labelSelector := labelSelector
			matcher, err := metav1.LabelSelectorAsSelector(&labelSelector)
			if err != nil {
				return false, err
			}
			if matcher.Matches(labels.Set(objMeta.GetLabels())) {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	// annotations
	if len(m.Annotations) > 0 {
		found := false
		for _, annotationSelector := range m.Annotations {
			labelSelectorFromAnnotations := metav1.LabelSelector{
				MatchExpressions: annotationSelector.convertMatchExpressions(),
				MatchLabels:      annotationSelector.MatchAnnotations,
			}

			matcher, err := metav1.LabelSelectorAsSelector(&labelSelectorFromAnnotations)
			if err != nil {
				return false, err
			}
			if matcher.Matches(labels.Set(objMeta.GetAnnotations())) {
				found = true
			}
		}
		if !found {
			return false, nil
		}
	}

	// content
	if len(m.Content) > 0 {
		j, err := json.Marshal(obj)
		if err != nil {
			return false, err
		}

		for _, content := range m.Content {
			res := gjson.GetBytes(j, content.Key)
			if !res.Exists() {
				return false, nil
			}
			switch content.Value.Type {
			case intstr.Int:
				if res.Int() != int64(content.Value.IntVal) {
					return false, nil
				}
			case intstr.String:
				if res.String() != content.Value.StrVal {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

type MatchedRules []SyncRule

func (r MatchedRules) GetOverrides() []resources.K8SResourceOverlayPatch {
	overrides := make([]resources.K8SResourceOverlayPatch, 0)
	for _, matchedRule := range r {
		overrides = append(overrides, matchedRule.Mutations.Overrides...)
	}

	return overrides
}

func (r MatchedRules) GetSyncStatus() bool {
	for _, matchedRule := range r {
		if matchedRule.Mutations.SyncStatus == true {
			return true
		}
	}

	return false
}

func (r ResourceSyncRuleSpec) Match(obj runtime.Object) (bool, MatchedRules, error) {
	matchedRules := make(MatchedRules, 0)

	if resources.ConvertGVK(obj.GetObjectKind().GroupVersionKind()) != r.GVK {
		return false, matchedRules, nil
	}

	for _, rule := range r.Rules {
		ok, err := rule.Match(obj)
		if err != nil {
			return false, matchedRules, err
		}
		if ok {
			matchedRules = append(matchedRules, rule)
		}
	}

	return len(matchedRules) > 0, matchedRules, nil
}

func (s *ResourceSyncRule) Match(obj runtime.Object) (bool, MatchedRules, error) {
	return s.Spec.Match(obj)
}

func (r *SyncRule) Match(obj runtime.Object) (bool, error) {
	if len(r.Matches) == 0 {
		return true, nil
	}

	for _, m := range r.Matches {
		if ok, err := m.Match(obj); ok && err == nil {
			return ok, nil
		} else if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (s AnnotationSelector) convertMatchExpressions() []metav1.LabelSelectorRequirement {
	reqs := make([]metav1.LabelSelectorRequirement, 0)
	for _, r := range s.MatchExpressions {
		values := make([]string, len(r.Values))
		for i, v := range r.Values {
			values[i] = string(v)
		}
		reqs = append(reqs, metav1.LabelSelectorRequirement{
			Key:      r.Key,
			Operator: r.Operator,
			Values:   values,
		})
	}

	return reqs
}
