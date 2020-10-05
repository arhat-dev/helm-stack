package cmd

import (
	"fmt"
	"sort"

	"arhat.dev/helm-stack/pkg/conf"
)

func GetEnvironmentsToRun(names []string, config *conf.ResolvedConfig) ([]*conf.Environment, error) {
	var toRun []*conf.Environment

	for _, name := range names {
		switch name {
		case "all":
			var all []*conf.Environment
			for n := range config.Environments {
				all = append(all, config.Environments[n])
			}
			toRun = all

			goto doSort
		default:
			e, ok := config.Environments[name]
			if !ok {
				return nil, fmt.Errorf("no such environment with name %q", name)
			}

			toRun = append(toRun, e)
		}
	}

doSort:
	sort.Slice(toRun, func(i, j int) bool {
		return toRun[i].Name < toRun[j].Name
	})

	return toRun, nil
}
