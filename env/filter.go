package env

import (
	"strconv"

	"github.com/rs/zerolog/log"
)

//Filter filter
type Filter struct {
	tags       map[string]bool
	apps       map[string]bool
	priorities map[int]bool
}

func createFilter(tags, apps, priorities []string) *Filter {
	return &Filter{
		tags:       list2StringMap(tags),
		apps:       list2StringMap(apps),
		priorities: list2IntMap(priorities),
	}
}

func (f *Filter) exclude(appName string, app *App) bool {

	// filter app names
	if len(f.apps) > 0 {
		_, exists := f.apps[appName]
		if !exists {
			return true
		}
	}

	// filter tags
	if len(f.tags) > 0 {
		contains := f.tagsContains(app.Tags)
		if !contains {
			return true
		}
	}

	// filter priorities
	if len(f.priorities) > 0 {
		_, exists := f.priorities[app.Priority]
		if !exists {
			return true
		}
	}

	return false
}

func (f *Filter) tagsContains(tags []string) bool {
	if len(tags) == 0 {
		return false
	}
	contains := false
	for i := 0; i < len(tags) && !contains; i++ {
		_, exists := f.tags[tags[i]]
		contains = contains || exists
	}
	return contains
}

func list2StringMap(data []string) map[string]bool {
	result := make(map[string]bool)
	if len(data) > 0 {
		for _, t := range data {
			result[t] = true
		}
	}
	return result
}

func list2IntMap(data []string) map[int]bool {
	result := make(map[int]bool)
	if len(data) > 0 {
		for _, p := range data {
			i, err := strconv.Atoi(p)
			if err != nil {
				log.Fatal().Err(err).Str("priority", p).Msg("Error convert priority to number")
			}
			result[i] = true
		}
	}
	return result
}
