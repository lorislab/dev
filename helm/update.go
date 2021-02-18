package helm

import (
	"sync"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{"error": err}).Fatal("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.WithField("error", err).Fatal("Erorr load chart repository before update")
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
	log.Debug("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Warnf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				log.Debugf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	log.Debug("Update Complete. ⎈Happy Helming!⎈")
}
