package main

import (
	"github.com/sirupsen/logrus"
	veleroplugin "github.com/vmware-tanzu/velero/pkg/plugin/framework"

	"github.com/mpryc/oadp-pvc-remove-selector/internal/plugin"
)

func main() {
	veleroplugin.NewServer().
		RegisterRestoreItemAction("mpryc.io/pvc-remove-selector", newPVCRestoreItemAction).
		Serve()
}

func newPVCRestoreItemAction(logger logrus.FieldLogger) (interface{}, error) {
	return plugin.NewPVCRemoveSelectorPlugin(logger), nil
}

