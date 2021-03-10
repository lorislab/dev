package helm

import (
	"time"

	"github.com/lorislab/dev/pkg/api"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

//Upgrade upgrade application
func Upgrade(hc *api.HelmCluster, app *api.App, version string, wait bool) (*release.Release, error) {
	client, err := createUpgradeClient(app.Namespace)
	if err != nil {
		return nil, err
	}
	client.Namespace = app.Namespace
	client.Version = app.Helm.Version
	client.Wait = wait

	// chart request
	chartRequested, err := chartRequest(client.ChartPathOptions, app)
	if err != nil {
		return nil, err
	}

	// values
	vals, err := chartValues(app, hc)
	if err != nil {
		return nil, err
	}

	return client.Run(app.Name, chartRequested, vals)
}

func createUpgradeClient(namespace string) (*action.Upgrade, error) {
	cfg, err := newActionConfig(namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewUpgrade(cfg)
	client.Install = false
	client.DryRun = false
	client.DisableHooks = false
	client.Timeout = 300 * time.Second
	client.WaitForJobs = false
	client.Description = ""
	client.Devel = false
	client.DisableOpenAPIValidation = false
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = false
	client.Wait = false
	client.Force = false
	client.ResetValues = false
	client.ReuseValues = false
	client.MaxHistory = settings.MaxHistory
	client.CleanupOnFail = false

	client.ChartPathOptions.Version = ""
	client.ChartPathOptions.Verify = false
	client.ChartPathOptions.Keyring = defaultKeyring()
	client.ChartPathOptions.RepoURL = ""
	client.ChartPathOptions.Username = ""
	client.ChartPathOptions.Password = ""
	client.ChartPathOptions.CertFile = ""
	client.ChartPathOptions.KeyFile = ""
	client.ChartPathOptions.InsecureSkipTLSverify = false
	client.ChartPathOptions.CaFile = ""
	return client, nil
}
