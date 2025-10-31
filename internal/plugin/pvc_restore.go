package plugin

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
)

// PVCRemoveSelectorPlugin removes spec.selector from PVCs during restore
type PVCRemoveSelectorPlugin struct {
	log logrus.FieldLogger
}

func NewPVCRemoveSelectorPlugin(log logrus.FieldLogger) *PVCRemoveSelectorPlugin {
	return &PVCRemoveSelectorPlugin{log: log}
}

func (p *PVCRemoveSelectorPlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"persistentvolumeclaims"},
	}, nil
}

func (p *PVCRemoveSelectorPlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.log.Info("Processing PVC for selector removal")

	unstruct, ok := input.Item.(*unstructured.Unstructured)
	if !ok {
		p.log.Error("Failed to convert input.Item to *unstructured.Unstructured")
		return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
	}

	item := unstruct.DeepCopy()

	ns := item.GetNamespace()
	name := item.GetName()
	p.log.Infof("Processing PVC: %s/%s", ns, name)

	if selector, found, _ := unstructured.NestedMap(item.Object, "spec", "selector"); found && selector != nil {
		p.log.Info("Removing spec.selector from PVC")
		unstructured.RemoveNestedField(item.Object, "spec", "selector")
	}

	// Remove status to allow Kubernetes to regenerate it
	unstructured.RemoveNestedField(item.Object, "status")

	return velero.NewRestoreItemActionExecuteOutput(item), nil
}

