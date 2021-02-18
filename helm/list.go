package helm

import (
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

type ListConfig struct {
	AllNamespaces bool
	All           bool
}

func List() ([]*release.Release, error) {
	config := ListConfig{
		AllNamespaces: true,
		All:           true,
	}
	cfg := &action.Configuration{}
	client := action.NewList(cfg)
	client.AllNamespaces = config.AllNamespaces
	client.All = config.AllNamespaces

	if client.AllNamespaces {
		if err := cfg.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), debug); err != nil {
			return nil, err
		}
	}
	client.SetStateMask()

	results, err := client.Run()
	if err != nil {
		return nil, err
	}
	return results, nil
}
