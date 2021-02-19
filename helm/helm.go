package helm

import (
	"fmt"
	"os"

	"github.com/lorislab/dev/pkg/api"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	yaml2 "sigs.k8s.io/yaml"
)

var settings = cli.New()

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func newActionConfig(namespace string) (*action.Configuration, error) {

	cfg := &action.Configuration{}

	c, err := settings.RESTClientGetter().ToRESTConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error create kube rest config")
		return nil, err
	}
	restClientGetter := newConfigFlags(c, namespace)

	if err := cfg.Init(restClientGetter, namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Error().Err(err).Msg("Error create action config")
		return nil, err
	}
	return cfg, nil
}

func newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Debug().Msg(fmt.Sprintf(format, v...))
	}
}

func chartRequest(chartOptions action.ChartPathOptions, app *api.App) (*chart.Chart, error) {
	cp, err := chartOptions.LocateChart(app.Helm.Chart, settings)
	if err != nil {
		log.Error().Err(err).Str("app", app.Name).Msg("Error localte chart")
		return nil, err
	}
	log.Debug().Str("app", app.Name).Str("chart", cp).Msg("Helm chart")

	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Error().Err(err).Str("app", app.Name).Msg("Error merge values")
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		log.Error().Err(err).Str("app", app.Name).Msg("Error check if installable")
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		log.Warn().Str("app", app.Name).Msg("This chart is deprecated")
	}
	return chartRequested, nil
}

func chartValues(app *api.App) (map[string]interface{}, error) {
	p := getter.All(settings)
	vf := []string{}
	if len(app.Helm.ValuesFiles) > 0 {
		vf = append(vf, app.Helm.ValuesFiles...)
	}

	valueOpts := &values.Options{ValueFiles: vf, Values: []string{}, StringValues: []string{}, FileValues: []string{}}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Error().Err(err).Str("app", app.Name).Msg("Error merge values")
		return nil, err
	}

	// add embedded values
	if app.Helm.Values != nil {
		tmp := app.Helm.Values
		bytes, err := yaml.Marshal(&tmp)
		if err != nil {
			log.Error().Err(err).Str("app", app.Name).Msg("Error marschale values")
			return nil, err
		}
		em := map[string]interface{}{}
		if err := yaml2.Unmarshal(bytes, &em); err != nil {
			log.Error().Err(err).Str("app", app.Name).Msg("Error unmarschale values")
			return nil, err
		}
		vals = mergeMaps(vals, em)
	}

	return vals, nil
}
