package project

import (
	"sort"

	"golang.org/x/net/context"
)

// Ps list containers for the specified services.
func (p *Project) Ps(ctx context.Context, services ...string) (InfoSet, error) {
	allInfo := InfoSet{}

	if services != nil {
		sort.Strings(services)
	}

	for _, name := range p.ServiceConfigs.Keys() {
		if services != nil { // apply filter
			index := sort.SearchStrings(services, name)
			// index hold the position where the data should be,
			// be it present or not
			if index > len(services) || services[index] != name {
				continue
			}
		}

		service, err := p.CreateService(name)
		if err != nil {
			return nil, err
		}

		info, err := service.Info(ctx)
		if err != nil {
			return nil, err
		}

		allInfo = append(allInfo, info...)
	}
	return allInfo, nil
}
