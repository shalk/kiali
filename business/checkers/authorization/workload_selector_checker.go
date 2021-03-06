package authorization

import (
	"k8s.io/apimachinery/pkg/labels"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
)

type WorkloadSelectorChecker struct {
	AuthorizationPolicy kubernetes.IstioObject
	WorkloadList        models.WorkloadList
}

func (wsc WorkloadSelectorChecker) Check() ([]*models.IstioCheck, bool) {
	checks, valid := make([]*models.IstioCheck, 0), true
	if selectorSpec, found := wsc.AuthorizationPolicy.GetSpec()["selector"]; found {
		if ml, ok := selectorSpec.(map[string]interface{}); ok {
			if matchLabels, found := ml["matchLabels"]; found {
				if selectors, ok := matchLabels.(map[string]interface{}); ok {
					labelSelectors := make(map[string]string, len(selectors))
					for k, v := range selectors {
						labelSelectors[k] = v.(string)
					}
					if !wsc.hasMatchingWorkload(labelSelectors) {
						check := models.Build("authorizationpolicy.selector.workloadnotfound", "spec/selector")
						checks = append(checks, &check)
					}
				}
			}
		}

	}
	return checks, valid
}

func (wsc WorkloadSelectorChecker) hasMatchingWorkload(labelSelector map[string]string) bool {
	selector := labels.SelectorFromSet(labels.Set(labelSelector))

	for _, wl := range wsc.WorkloadList.Workloads {
		wlLabelSet := labels.Set(wl.Labels)
		if selector.Matches(wlLabelSet) {
			return true
		}
	}
	return false
}
