package response

import (
	"encoding/json"
	"mutating-trace-admission-controller/pkg/util/patch"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildDeamonSetPatch(raw []byte, newAnnotations map[string]string) *v1beta1.AdmissionResponse {
	var deamonSet appv1.DaemonSet
	err := json.Unmarshal(raw, &deamonSet)
	if err != nil {
		glog.Errorf("unmarshal deamonset raw failed: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := patch.Encode(patch.WithAnnotations(deamonSet.Annotations, newAnnotations))
	if err != nil {
		glog.Errorf("encode deamonset patch failed: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
