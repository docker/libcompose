package libcompose

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

// a minimalist and uncomplete Compose YAML parser we used for prototyping and demos
// it converts a Compose YAML file into a list of Services

// TODO (gdevillele) check for duplicated port use on the same host

// ParseServicesYml parses a compose.yml file content
// and returns an ordered list of services to create
func ParseServicesYml(ymlString []byte) ([]*Service, error) {

	// create a map of services and fill it with YAML structure
	servicesMap := make(map[string]Service)
	err := yaml.Unmarshal(ymlString, &servicesMap)
	if err != nil {
		return nil, err
	}

	// construct an array of services in which we will put the services from servicesMap
	var services []*Service

	// loop over the services map[string]Service
	// and set the Name field of each service structure
	for name, _ := range servicesMap {
		// create a instance we can keep (copy)
		service := servicesMap[name]
		service.Name = name
		services = append(services, &service)
	}

	// validate YML
	//-------------------------------
	fmt.Println("[libcompose] Validating YAML...")

	// loop over the services to check for dependency errors
	// TODO: return the error to transmit it to the libcompose caller
	success := true
	for _, service := range services {
		var expandedServices []*Service
		success = CheckDependencies(service, services, expandedServices)
		if success == false {
			break
		}
	}
	if success == false {
		return services, errors.New("[libcompose] YAML validation error")
	}

	// Service Reordering (for containers creation)
	//-------------------------------
	fmt.Println("[libcompose] Reordering services for creation...")
	// we are doing several passes over the services and detecting
	// which service we can create at each pass
	// if not a single service can be build during a pass,
	// then it is impossible to build the stack described
	var servicesToOrder []*Service
	var orderedServices []*Service

	// populate servicesToOrder with all the services
	for _, s := range services {
		servicesToOrder = append(servicesToOrder, s)
	}

	// at this point, servicesToOrder is full and orderedServices is empty

	// do passes while we have remaining services to order
	serviceOrderedDuringThePass := true
	for len(servicesToOrder) > 0 && serviceOrderedDuringThePass {

		serviceOrderedDuringThePass = false

		// do a pass (we find no more than 1 service per pass)
		for i, s := range servicesToOrder {

			// check if s has either no dependencies or all its dependencies are already in orderedServices
			if s.HasNoDependencies() || s.AreDependenciesIn(orderedServices) {
				// add the current service in orderedServices
				orderedServices = append(orderedServices, s)

				// and remove it from servicesToOrder
				if i < len(servicesToOrder)-1 {
					servicesToOrder = append(servicesToOrder[:i], servicesToOrder[i+1:]...)
				} else {
					servicesToOrder = servicesToOrder[:i]
				}

				// and increment the per-pass ordered services counter
				serviceOrderedDuringThePass = true
				break
			}
		}
	}
	// end of the passes, if services remain in "servicesToOrder"
	// it means we cannot define a build order for the services
	if len(servicesToOrder) > 0 {
		// ERROR: cannot build this stack
		fmt.Println("[libcompose] ERROR: can't build this stack")
		return services, errors.New("cannot reorder services for creation creation order for this stack")
	}

	// print unordered services
	fmt.Printf("    unordered services : ")
	for _, s := range services {
		fmt.Printf("%s ", s.Name)
	}
	fmt.Printf("\n")

	// print ordered services
	fmt.Printf("    ordered services   : ")
	for _, s := range orderedServices {
		fmt.Printf("%s ", s.Name)
	}
	fmt.Printf("\n")

	return orderedServices, nil
}

// Checks the dependency tree of a service (recursive function)
//
// params :
// - service : service whose we check the dependencies
// - allServices : all the services of the stack
// - unavailableServices : contains the services which have been expanded already during the recursion
//
// returns whether dependency check succeeded
func CheckDependencies(service *Service, allServices []*Service, expandedServices []*Service) bool {

	// retrieve the names of the direct dependencies for this service
	dependencies := service.GetDependenciesNames()

	// the service has dependencies
	if len(dependencies) > 0 {

		// we will recursively expand the dependecy tree
		// we add the service into the expandedServices slice
		expandedServices = append(expandedServices, service)

		result := true
		for _, dependencyName := range dependencies {

			// get the service corresponding to the link // TODO: check that the link exists (it is in all services)
			dependencyService, err := findServiceByName(allServices, dependencyName)
			if err != nil {
				// ERROR : link does not refer to a known service
				fmt.Println("[libcompose] ERROR: link does not reference a known service")
				result = false
				break
			}

			// check that the dependency is not in the unavailable services list
			_, err = findServiceByName(expandedServices, dependencyName)
			if err == nil { // link has been found
				// ERROR : cyclic inclusion
				fmt.Println("[libcompose] ERROR: cyclic inclusion")
				result = false
				break
			}

			// link look fine, we expand it and check for its dependencies
			result = CheckDependencies(dependencyService, allServices, expandedServices)
			if result == false {
				break
			}
		}

		return result

	} else {

		// the service has no dependency
		// so we consider the dependency check a success
		return true
	}
}

// returns a Service or nil
func findServiceByName(services []*Service, name string) (*Service, error) {
	for _, s := range services {
		if s.Name == name {
			return s, nil
		}
	}
	return &Service{}, errors.New("service not found")
}
