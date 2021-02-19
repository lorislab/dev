package helm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
)

const searchMaxScore = 25

//Search helm search
func Search() ([]*search.Result, error) {
	var args []string
	o := &helmSearch{}
	o.repoFile = settings.RepositoryConfig
	o.repoCacheDir = settings.RepositoryCache
	o.devel = true
	o.versions = true

	o.setupSearchedVersion()

	index, err := o.buildIndex()
	if err != nil {
		return nil, err
	}

	var res []*search.Result
	if len(args) == 0 {
		res = index.All()
	} else {
		q := strings.Join(args, " ")
		res, err = index.Search(q, searchMaxScore, o.regexp)
		if err != nil {
			return nil, err
		}
	}
	search.SortScore(res)
	data, err := o.applyConstraint(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type helmSearch struct {
	version      string
	devel        bool
	repoFile     string
	versions     bool
	repoCacheDir string
	regexp       bool
}

func (o *helmSearch) setupSearchedVersion() {
	log.Debug().Str("version", o.version).Msg("Original chart version")

	if o.version != "" {
		return
	}

	if o.devel {
		log.Debug().Msg("set version to >=0.0.0-0")
		o.version = ">=0.0.0-0"
	} else {
		log.Debug().Msg("set version to >0.0.0")
		o.version = ">0.0.0"
	}
}

func (o *helmSearch) buildIndex() (*search.Index, error) {

	rf, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(rf.Repositories) == 0 {
		return nil, errors.New("no repositories configured")
	}

	i := search.NewIndex()
	for _, re := range rf.Repositories {
		n := re.Name
		f := filepath.Join(o.repoCacheDir, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			log.Warn().Str("repo", n).Msg("Repois corrupt or missing. Try 'helm repo update'.")
			continue
		}

		i.AddRepo(n, ind, o.versions || len(o.version) > 0)
	}
	return i, nil
}

func (o *helmSearch) applyConstraint(res []*search.Result) ([]*search.Result, error) {
	if len(o.version) == 0 {
		return res, nil
	}

	constraint, err := semver.NewConstraint(o.version)
	if err != nil {
		return res, errors.Wrap(err, "an invalid version/constraint format")
	}

	data := res[:0]
	foundNames := map[string]bool{}
	appendSearchResults := func(res *search.Result) {
		data = append(data, res)
		if !o.versions {
			foundNames[res.Name] = true // If user hasn't requested all versions, only show the latest that matches
		}
	}
	for _, r := range res {
		if _, found := foundNames[r.Name]; found {
			continue
		}
		v, err := semver.NewVersion(r.Chart.Version)

		if err != nil {
			// If the current version number check appears ErrSegmentStartsZero or ErrInvalidPrerelease error and not devel mode, ignore
			if (err == semver.ErrSegmentStartsZero || err == semver.ErrInvalidPrerelease) && !o.devel {
				continue
			}
			appendSearchResults(r)
		} else if constraint.Check(v) {
			appendSearchResults(r)
		}
	}

	return data, nil
}

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}
