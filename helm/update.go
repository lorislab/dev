package helm

import (
	"sync"

	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func Update() {

	o := &helmUpdate{
		repoFile:  settings.RepositoryConfig,
		repoCache: settings.RepositoryCache,
		update:    updateCharts,
	}

	f, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(f.Repositories) == 0 {
		log.Fatal().Err(err).Msg("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.Fatal().Err(err).Msg("Erorr load chart repository before update")
		}
		if o.repoCache != "" {
			r.CachePath = o.repoCache
		}
		repos = append(repos, r)
	}

	o.update(repos)
}

type helmUpdate struct {
	update    func([]*repo.ChartRepository)
	repoFile  string
	repoCache string
}

func updateCharts(repos []*repo.ChartRepository) {
	log.Debug().Msg("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Warn().Err(err).Str("url", re.Config.URL).Str("repository", re.Config.Name).Msg("...Unable to get an update from the chart repository")
			} else {
				log.Debug().Str("repository", re.Config.Name).Msg("...Successfully got an update from the chart repository")
			}
		}(re)
	}
	wg.Wait()
	log.Debug().Msg("Update Complete. ⎈Happy Helming!⎈")
}
