package helm

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lorislab/dev/pkg/api"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/util/homedir"
)

//Install install application
func Install(app *api.App, version string, wait bool) (*release.Release, error) {

	client, err := createClient(app.Namespace)
	if err != nil {
		return nil, err
	}
	client.Namespace = app.Namespace
	client.Version = version
	client.Wait = wait
	client.ReleaseName = app.Name
	client.CreateNamespace = true

	// chart request
	chartRequested, err := chartRequest(client.ChartPathOptions, app)
	if err != nil {
		return nil, err
	}

	// values
	vals, err := chartValues(app)
	if err != nil {
		return nil, err
	}

	return client.Run(chartRequested, vals)
}

func createClient(namespace string) (*action.Install, error) {
	cfg, err := newActionConfig(namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(cfg)
	client.Namespace = namespace
	client.CreateNamespace = false
	client.DryRun = false
	client.DisableHooks = false
	client.Replace = false
	client.Timeout = 300 * time.Second
	client.WaitForJobs = false
	client.GenerateName = false
	client.NameTemplate = ""
	client.Description = ""
	client.Devel = false
	client.DependencyUpdate = false
	client.DisableOpenAPIValidation = false
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = false

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

func defaultKeyring() string {
	if v, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(v, "pubring.gpg")
	}
	return filepath.Join(homedir.HomeDir(), ".gnupg", "pubring.gpg")
}

func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
