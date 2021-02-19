package helm

import (
	"time"

	"github.com/lorislab/dev/pkg/api"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

//Uninstall uninstall application
func Uninstall(app *api.App, wait bool) (*release.UninstallReleaseResponse, error) {

	cfg, err := newActionConfig(app.Namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUninstall(cfg)
	client.DryRun = false
	client.DisableHooks = false
	client.KeepHistory = false
	client.Timeout = 300 * time.Second
	client.Description = ""
	return client.Run(app.Name)

}
