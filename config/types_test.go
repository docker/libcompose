package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependentServices(t *testing.T) {
	serviceConfigs := NewServiceConfigs()
	serviceConfigs.Add("serviceA", &ServiceConfig{})
	serviceConfigs.Add("serviceB", &ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceConfigs.Add("serviceC", &ServiceConfig{
		VolumesFrom: []string{"serviceB"},
	})

	serviceConfigs.Add("serviceD", &ServiceConfig{
		VolumesFrom: []string{"serviceB"},
		DependsOn:   []string{"serviceC"},
	})

	serviceConfigs.Add("serviceE", &ServiceConfig{
		NetworkMode: "service:serviceA",
	})

	serviceADeps := serviceConfigs.DependentServices("serviceA")
	assert.Equal(t, true, containsString(serviceADeps, "serviceA"))
	assert.Equal(t, false, containsString(serviceADeps, "serviceB"))
	assert.Equal(t, false, containsString(serviceADeps, "serviceC"))
	assert.Equal(t, false, containsString(serviceADeps, "serviceD"))
	assert.Equal(t, false, containsString(serviceADeps, "serviceE"))
	assert.Equal(t, 1, len(serviceADeps))

	serviceBDeps := serviceConfigs.DependentServices("serviceB")
	assert.Equal(t, true, containsString(serviceBDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceBDeps, "serviceB"))
	assert.Equal(t, false, containsString(serviceBDeps, "serviceC"))
	assert.Equal(t, false, containsString(serviceBDeps, "serviceD"))
	assert.Equal(t, false, containsString(serviceBDeps, "serviceE"))
	assert.Equal(t, 2, len(serviceBDeps))

	serviceCDeps := serviceConfigs.DependentServices("serviceC")
	assert.Equal(t, true, containsString(serviceCDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceCDeps, "serviceB"))
	assert.Equal(t, true, containsString(serviceCDeps, "serviceC"))
	assert.Equal(t, false, containsString(serviceCDeps, "serviceD"))
	assert.Equal(t, false, containsString(serviceCDeps, "serviceE"))
	assert.Equal(t, 3, len(serviceCDeps))

	serviceDDeps := serviceConfigs.DependentServices("serviceD")
	assert.Equal(t, true, containsString(serviceDDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceDDeps, "serviceB"))
	assert.Equal(t, true, containsString(serviceDDeps, "serviceC"))
	assert.Equal(t, true, containsString(serviceDDeps, "serviceD"))
	assert.Equal(t, false, containsString(serviceDDeps, "serviceE"))
	assert.Equal(t, 4, len(serviceDDeps))

	serviceEDeps := serviceConfigs.DependentServices("serviceE")
	assert.Equal(t, true, containsString(serviceEDeps, "serviceA"))
	assert.Equal(t, false, containsString(serviceEDeps, "serviceB"))
	assert.Equal(t, false, containsString(serviceEDeps, "serviceC"))
	assert.Equal(t, false, containsString(serviceEDeps, "serviceD"))
	assert.Equal(t, true, containsString(serviceEDeps, "serviceE"))
	assert.Equal(t, 2, len(serviceEDeps))
}

func TestDependentServicesInvalid(t *testing.T) {
	serviceConfigs := NewServiceConfigs()
	serviceConfigs.Add("serviceA", &ServiceConfig{})
	serviceConfigs.Add("serviceB", &ServiceConfig{
		Links: []string{"foobar"},
	})

	serviceADeps := serviceConfigs.DependentServices("serviceA")
	assert.Equal(t, true, containsString(serviceADeps, "serviceA"))
	assert.Equal(t, false, containsString(serviceADeps, "serviceB"))
	assert.Equal(t, false, containsString(serviceADeps, "foobar"))
	assert.Equal(t, 1, len(serviceADeps))

	serviceBDeps := serviceConfigs.DependentServices("serviceB")
	assert.Equal(t, false, containsString(serviceBDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceBDeps, "serviceB"))
	assert.Equal(t, false, containsString(serviceBDeps, "foobar"))
	assert.Equal(t, 1, len(serviceBDeps))
}

func TestDependentServicesDuplicate(t *testing.T) {
	serviceConfigs := NewServiceConfigs()
	serviceConfigs.Add("serviceA", &ServiceConfig{})
	serviceConfigs.Add("serviceB", &ServiceConfig{
		Links: []string{"serviceA"},
	})
	serviceConfigs.Add("serviceC", &ServiceConfig{
		Links:       []string{"serviceA"},
		VolumesFrom: []string{"serviceA"},
		DependsOn:   []string{"serviceB"},
	})

	serviceCDeps := serviceConfigs.DependentServices("serviceC")
	assert.Equal(t, true, containsString(serviceCDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceCDeps, "serviceB"))
	assert.Equal(t, true, containsString(serviceCDeps, "serviceC"))
	assert.Equal(t, 3, len(serviceCDeps))
}

func TestDependentServicesLoop(t *testing.T) {
	// Loops are invalid configuration, however its good to know
	// the code wont loop
	serviceConfigs := NewServiceConfigs()
	serviceConfigs.Add("serviceA", &ServiceConfig{
		Links: []string{"serviceB"},
	})
	serviceConfigs.Add("serviceB", &ServiceConfig{
		Links: []string{"serviceA"},
	})

	serviceADeps := serviceConfigs.DependentServices("serviceA")
	assert.Equal(t, true, containsString(serviceADeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceADeps, "serviceB"))
	assert.Equal(t, 2, len(serviceADeps))

	serviceBDeps := serviceConfigs.DependentServices("serviceB")
	assert.Equal(t, true, containsString(serviceBDeps, "serviceA"))
	assert.Equal(t, true, containsString(serviceBDeps, "serviceB"))
	assert.Equal(t, 2, len(serviceBDeps))
}

func containsString(set []string, str string) bool {
	for _, s := range set {
		if str == s {
			return true
		}
	}
	return false
}
