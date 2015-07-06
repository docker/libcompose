package main

import (
	"fmt"
	"github.com/docker/libcompose"
	"io/ioutil"
)

func main() {

	data, err := ioutil.ReadFile("../samples/compose.yml")
	if err != nil {
		fmt.Println("Error: unable to read ../samples/compose.yml file")
		return
	}

	services, err := libcompose.ParseServicesYml(data)
	fmt.Printf("%+v\n", services)
}
