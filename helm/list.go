package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

//List execute helm list
func List() ([]*release.Release, error) {
	cfg, err := newActionConfig("")
	if err != nil {
		return nil, err
	}
	client := action.NewList(cfg)
	client.AllNamespaces = true
	client.All = true
	client.SetStateMask()

	return client.Run()
}
