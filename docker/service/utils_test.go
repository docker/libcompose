package service

import (
	"testing"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/stretchr/testify/assert"
)

func TestDefaultDependentServices(t *testing.T) {
	p := project.NewProject(&project.Context{
		ServiceFactory: NewFactory(&ctx.Context{}),
	}, nil, nil)
	serviceConfigs := config.NewServiceConfigs()
	p.ServiceConfigs = serviceConfigs

	serviceConfigs.Add("serviceA", &config.ServiceConfig{})
	serviceConfigs.Add("serviceB", &config.ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceConfigs.Add("serviceC", &config.ServiceConfig{
		VolumesFrom: []string{"serviceB"},
	})

	serviceConfigs.Add("serviceD", &config.ServiceConfig{
		VolumesFrom: []string{"serviceB"},
		DependsOn:   []string{"serviceC"},
	})

	serviceConfigs.Add("serviceE", &config.ServiceConfig{
		NetworkMode: "service:serviceA",
	})

	serviceA, err := p.CreateService("serviceA")
	assert.Nil(t, err)
	serviceADeps := DefaultDependentServices(p, serviceA)
	assert.Equal(t, 0, len(serviceADeps))

	serviceB, err := p.CreateService("serviceB")
	assert.Nil(t, err)
	serviceBDeps := DefaultDependentServices(p, serviceB)
	assert.Equal(t, true, containsService(serviceBDeps, "serviceA"))
	assert.Equal(t, 1, len(serviceBDeps))

	serviceC, err := p.CreateService("serviceC")
	assert.Nil(t, err)
	serviceCDeps := DefaultDependentServices(p, serviceC)
	assert.Equal(t, true, containsService(serviceCDeps, "serviceB"))
	assert.Equal(t, 1, len(serviceCDeps))

	serviceD, err := p.CreateService("serviceD")
	assert.Nil(t, err)
	serviceDDeps := DefaultDependentServices(p, serviceD)
	assert.Equal(t, true, containsService(serviceDDeps, "serviceB"))
	assert.Equal(t, true, containsService(serviceDDeps, "serviceC"))
	assert.Equal(t, 2, len(serviceDDeps))

	serviceE, err := p.CreateService("serviceE")
	assert.Nil(t, err)
	serviceEDeps := DefaultDependentServices(p, serviceE)
	assert.Equal(t, true, containsService(serviceEDeps, "serviceA"))
	assert.Equal(t, 1, len(serviceEDeps))
}

func TestDefaultDependentServicesInvalid(t *testing.T) {
	p := project.NewProject(&project.Context{
		ServiceFactory: NewFactory(&ctx.Context{}),
	}, nil, nil)
	serviceConfigs := config.NewServiceConfigs()
	p.ServiceConfigs = serviceConfigs

	serviceConfigs.Add("serviceA", &config.ServiceConfig{})
	serviceConfigs.Add("serviceB", &config.ServiceConfig{
		Links: []string{"foobar"},
	})

	serviceA, err := p.CreateService("serviceA")
	assert.Nil(t, err)
	serviceADeps := DefaultDependentServices(p, serviceA)
	assert.Equal(t, 0, len(serviceADeps))

	serviceB, err := p.CreateService("serviceB")
	assert.Nil(t, err)
	serviceBDeps := DefaultDependentServices(p, serviceB)
	assert.Equal(t, true, containsService(serviceBDeps, "foobar"))
	assert.Equal(t, 1, len(serviceBDeps))
}

func TestDependentServicesRecursive(t *testing.T) {
	p := project.NewProject(&project.Context{
		ServiceFactory: NewFactory(&ctx.Context{}),
	}, nil, nil)
	serviceConfigs := config.NewServiceConfigs()
	p.ServiceConfigs = serviceConfigs

	serviceConfigs.Add("serviceA", &config.ServiceConfig{})
	serviceConfigs.Add("serviceB", &config.ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceConfigs.Add("serviceC", &config.ServiceConfig{
		VolumesFrom: []string{"serviceB"},
	})

	serviceConfigs.Add("serviceD", &config.ServiceConfig{
		VolumesFrom: []string{"serviceB"},
		DependsOn:   []string{"serviceC"},
	})

	serviceConfigs.Add("serviceE", &config.ServiceConfig{
		NetworkMode: "service:serviceA",
	})

	serviceA, err := p.CreateService("serviceA")
	assert.Nil(t, err)
	serviceADeps := DefaultDependentServices(p, serviceA)
	assert.Equal(t, 0, len(serviceADeps))

	serviceB, err := p.CreateService("serviceB")
	assert.Nil(t, err)
	serviceBDeps := RecursiveDependentServices(p, serviceB)
	assert.Equal(t, true, containsService(serviceBDeps, "serviceA"))
	assert.Equal(t, 1, len(serviceBDeps))

	serviceC, err := p.CreateService("serviceC")
	assert.Nil(t, err)
	serviceCDeps := RecursiveDependentServices(p, serviceC)
	assert.Equal(t, true, containsService(serviceCDeps, "serviceA"))
	assert.Equal(t, true, containsService(serviceCDeps, "serviceB"))
	assert.Equal(t, 2, len(serviceCDeps))

	serviceD, err := p.CreateService("serviceD")
	assert.Nil(t, err)
	serviceDDeps := RecursiveDependentServices(p, serviceD)
	assert.Equal(t, true, containsService(serviceDDeps, "serviceA"))
	assert.Equal(t, true, containsService(serviceDDeps, "serviceB"))
	assert.Equal(t, true, containsService(serviceDDeps, "serviceC"))
	assert.Equal(t, 4, len(serviceDDeps))

	serviceE, err := p.CreateService("serviceE")
	assert.Nil(t, err)
	serviceEDeps := RecursiveDependentServices(p, serviceE)
	assert.Equal(t, true, containsService(serviceEDeps, "serviceA"))
	assert.Equal(t, 1, len(serviceEDeps))
}

func TestDependentServicesRecursiveInvalid(t *testing.T) {
	p := project.NewProject(&project.Context{
		ServiceFactory: NewFactory(&ctx.Context{}),
	}, nil, nil)
	serviceConfigs := config.NewServiceConfigs()
	p.ServiceConfigs = serviceConfigs

	serviceConfigs.Add("serviceA", &config.ServiceConfig{
		Links: []string{"foobar"},
	})
	serviceConfigs.Add("serviceB", &config.ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceA, err := p.CreateService("serviceA")
	assert.Nil(t, err)
	serviceADeps := RecursiveDependentServices(p, serviceA)
	assert.Equal(t, true, containsService(serviceADeps, "foobar"))
	assert.Equal(t, 1, len(serviceADeps))

	serviceB, err := p.CreateService("serviceB")
	assert.Nil(t, err)
	serviceBDeps := RecursiveDependentServices(p, serviceB)
	assert.Equal(t, true, containsService(serviceBDeps, "serviceA"))
	assert.Equal(t, true, containsService(serviceBDeps, "foobar"))
	assert.Equal(t, 2, len(serviceBDeps))
}

func TestDependentServicesRecursiveLoop(t *testing.T) {
	// Loops are invalid configuration, however its good to know
	// the code wont loop
	p := project.NewProject(&project.Context{
		ServiceFactory: NewFactory(&ctx.Context{}),
	}, nil, nil)
	serviceConfigs := config.NewServiceConfigs()
	p.ServiceConfigs = serviceConfigs

	serviceConfigs.Add("serviceA", &config.ServiceConfig{
		Links: []string{"serviceB"},
	})
	serviceConfigs.Add("serviceB", &config.ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceA, err := p.CreateService("serviceA")
	assert.Nil(t, err)
	serviceADeps := RecursiveDependentServices(p, serviceA)
	assert.Equal(t, true, containsService(serviceADeps, "serviceA"))
	assert.Equal(t, true, containsService(serviceADeps, "serviceB"))
	assert.Equal(t, 2, len(serviceADeps))

	serviceB, err := p.CreateService("serviceB")
	assert.Nil(t, err)
	serviceBDeps := RecursiveDependentServices(p, serviceB)
	assert.Equal(t, true, containsService(serviceBDeps, "serviceA"))
	assert.Equal(t, true, containsService(serviceBDeps, "serviceB"))
	assert.Equal(t, 2, len(serviceBDeps))
}

func containsService(set []project.ServiceRelationship, str string) bool {
	for _, r := range set {
		if str == r.Target {
			return true
		}
	}
	return false
}
