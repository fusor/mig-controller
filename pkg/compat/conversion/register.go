package conversion

import (
	"github.com/konveyor/mig-controller/pkg/compat/conversion/appsv1"
	"github.com/konveyor/mig-controller/pkg/compat/conversion/batchv1beta"
	"k8s.io/apimachinery/pkg/runtime"
)

func RegisterConversions(s *runtime.Scheme) error {
	err := appsv1.RegisterConversions(s)
	if err != nil {
		return err
	}

	err = batchv1beta.RegisterConversions(s)
	if err != nil {
		return err
	}
	return nil
}
