package helm

import (
	"github.com/lorislab/dev/pkg/api"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

//Upgrade upgrade application
func Upgrade(app *api.App, version string, wait bool) (*release.Release, error) {

	cfg, err := newActionConfig(app.Namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewUpgrade(cfg)
	client.Namespace = app.Namespace
	client.Version = app.Helm.Version
	client.Wait = wait

	// chart request
	chartRequested, err := chartRequest(client.ChartPathOptions, app)

	// values
	vals, err := chartValues(app)
	if err != nil {
		return nil, err
	}

	return client.Run(app.Name, chartRequested, vals)
}
