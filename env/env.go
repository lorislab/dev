package env

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/lorislab/dev/helm"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

//AppItem application item
type AppItem struct {
	Namespace      string
	AppName        string
	Declaration    *App
	CurrentVersion *semver.Version
	NextVersion    *semver.Version
	Action         AppAction
	Chart          string
	ChartRepo      string
	Cluster        *release.Release
	Repo           *search.Result
}

//Update update command
func Update() {
	helm.Update()
}

//LoadApps load applications for the environments
func LoadApps(env *LocalEnvironment, tags, apps, priorities []string) (map[int][]*AppItem, []int) {

	// list all releases in the cluster
	releases := listAllReleases()
	log.Debug().Int("count", len(releases)).Msg("Load apps releases from clusterConfig...")

	// search for all application in the helm repositories
	searchResults := searchAllApplications()
	log.Debug().Int("count", len(searchResults)).Msg("Load apps releases from helm repo...")

	// create filter
	filter := createFilter(tags, apps, priorities)

	result := make(map[int][]*AppItem)

	for appName, app := range env.Apps {

		// exclude application
		if filter.exclude(appName, app) {
			continue
		}

		// check chart repository from the definition
		chartRepo, local := chartRepository(app)

		// cluster version
		namespace := env.Namespace(appName)
		id := id(namespace, appName)

		var currentVersion *semver.Version
		clusterVersion, exists := releases[id]
		if exists {
			currentVersion = createSemVer(clusterVersion.Chart.Metadata.Version)
		}

		// repository version
		var nextVersion *semver.Version
		var repoVersion *search.Result
		if local {
			nextVersion = localChartVersion(chartRepo, app.Helm.Version)
		} else {
			repoVersions, exists := searchResults[chartRepo]
			if exists {
				nextVersion, repoVersion = findLatestBaseOnTheConstraints(repoVersions, app.Helm.Version)
			}
		}

		// create action
		action := createAction(currentVersion, nextVersion)

		// create application item
		appItem := &AppItem{
			AppName:        appName,
			Namespace:      namespace,
			Declaration:    app,
			Cluster:        clusterVersion,
			CurrentVersion: currentVersion,
			Action:         action,
			Chart:          app.Helm.Chart,
			ChartRepo:      chartRepo,
			NextVersion:    nextVersion,
			Repo:           repoVersion,
		}
		list, exists := result[app.Priority]
		if !exists {
			list = make([]*AppItem, 0)
		}
		list = append(list, appItem)
		result[app.Priority] = list
	}

	// sort priority keys
	keys := make([]int, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return result, keys
}

func listAllReleases() map[string]*release.Release {
	result := make(map[string]*release.Release)
	items, err := helm.List()
	if err != nil {
		log.Fatal().Err(err).Msg("Error load list of all releases")
	}
	for _, item := range items {
		result[id(item.Namespace, item.Name)] = item
	}
	return result
}

func searchAllApplications() map[string][]*search.Result {
	items, err := helm.Search()
	if err != nil {
		log.Fatal().Err(err).Msg("Error search all applications")
	}
	result := make(map[string][]*search.Result)
	for _, item := range items {
		list, e := result[item.Name]
		if !e {
			list = make([]*search.Result, 0)
		}
		list = append(list, item)
		result[item.Name] = list
	}
	return result
}

func findLatestBaseOnTheConstraints(items []*search.Result, constraints string) (*semver.Version, *search.Result) {

	c := createConstraints(constraints)

	vs := make([]*semver.Version, len(items))
	for i, r := range items {
		vs[i] = createSemVer(r.Chart.Version)
	}
	sort.Sort(sort.Reverse(semver.Collection(vs)))

	for i, ver := range vs {
		if c.Check(ver) {
			return ver, items[i]
		}
	}
	return nil, nil
}

func localChartVersion(repo string, constraints string) *semver.Version {

	c := createConstraints(constraints)

	chart, err := loader.Load(repo)
	if err != nil {
		log.Fatal().Err(err).Str("repo", repo).Msg("Error loading the local helm chart")
	}
	version := createSemVer(chart.Metadata.Version)

	if c.Check(version) {
		return version
	}
	return nil
}

func createConstraints(constraints string) *semver.Constraints {
	c, err := semver.NewConstraint(constraints)
	if err != nil {
		log.Fatal().Err(err).Str("constraints", constraints).Msg("Error create constrains for the SemVer version")
	}
	return c
}

func chartRepository(app *App) (string, bool) {
	chartRepo := app.Helm.Chart
	local := false
	if strings.HasPrefix(app.Helm.Repo, "alias:") {
		chartRepo = strings.TrimPrefix(app.Helm.Repo, "alias:") + "/" + chartRepo
	} else if strings.HasPrefix(app.Helm.Repo, "file://") {
		chartRepo = strings.TrimPrefix(app.Helm.Repo, "file://")
		local = true
	}
	return chartRepo, local
}

func createSemVer(version string) *semver.Version {
	tmp, err := semver.NewVersion(version)
	if err != nil {
		log.Fatal().Str("version", version).Err(err).Msg("Error parsing version")
	}
	return tmp
}

func createAction(currentVersion, nextVersion *semver.Version) AppAction {
	if nextVersion == nil {
		return AppActionNotFound
	}
	if currentVersion == nil {
		return AppActionInstall
	}
	if currentVersion.Equal(nextVersion) {
		return AppActionNothing
	}
	if currentVersion.LessThan(nextVersion) {
		return AppActionUpgrade
	}
	return AppActionDowngrade
}
