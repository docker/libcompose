package app

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/spf13/cobra"
)

// ProjectAction is an adapter to allow the use of ordinary functions as libcompose actions.
// Any function that has the appropriate signature can be register as an action on a codegansta/cli command.
//
// cli.Command{
//		Name:   "ps",
//		Usage:  "List containers",
//		Action: app.WithProject(factory, app.ProjectPs),
//	}
type ProjectAction func(project *project.Project, c *cobra.Command)

// BeforeApp is an action that is executed before any cli command.
func BeforeApp(c *cobra.Command) error {
	verbose, _ := c.Flags().GetBool("verbose")
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Warning("Note: This is an experimental alternate implementation of the Compose CLI (https://github.com/docker/compose)")
	return nil
}

// WithProject is a helper function to create a cli.Command action with a ProjectFactory.
func WithProject(factory ProjectFactory, action ProjectAction) func(context *cobra.Command, args []string) {
	return func(context *cobra.Command, args []string) {
		p, err := factory.Create(context)
		if err != nil {
			logrus.Fatalf("Failed to read project: %v", err)
		}
		action(p, context)
	}
}

// ProjectPs lists the containers.
func ProjectPs(p *project.Project, c *cobra.Command) {
	allInfo := project.InfoSet{}
	qFlag, _ := c.Flags().GetBool("quite")
	for _, name := range p.Configs.Keys() {
		service, err := p.CreateService(name)
		if err != nil {
			logrus.Fatal(err)
		}

		info, err := service.Info(qFlag)
		if err != nil {
			logrus.Fatal(err)
		}

		allInfo = append(allInfo, info...)
	}

	os.Stdout.WriteString(allInfo.String(!qFlag))
}

// ProjectPort prints the public port for a port binding.
func ProjectPort(p *project.Project, c *cobra.Command) {
	if len(c.Flags().Args()) != 2 {
		logrus.Fatalf("Please pass arguments in the form: SERVICE PORT")
	}

	index, _ := c.Flags().GetInt("index")
	protocol, _ := c.Flags().GetString("protocol")

	service, err := p.CreateService(c.Flags().Args()[0])
	if err != nil {
		logrus.Fatal(err)
	}

	containers, err := service.Containers()
	if err != nil {
		logrus.Fatal(err)
	}

	if index < 1 || index > len(containers) {
		logrus.Fatalf("Invalid index %d", index)
	}

	output, err := containers[index-1].Port(fmt.Sprintf("%s/%s", c.Flags().Args()[1], protocol))
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(output)
}

// ProjectStop stops all services.
func ProjectStop(p *project.Project, c *cobra.Command) {
	err := p.Stop(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectDown brings all services down (stops and clean containers).
func ProjectDown(p *project.Project, c *cobra.Command) {
	err := p.Down(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectBuild builds or rebuilds services.
func ProjectBuild(p *project.Project, c *cobra.Command) {
	err := p.Build(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectCreate creates all services but do not start them.
func ProjectCreate(p *project.Project, c *cobra.Command) {
	err := p.Create(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectUp brings all services up.
func ProjectUp(p *project.Project, c *cobra.Command) {
	err := p.Up(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}

	if background, _ := c.Flags().GetBool("detach"); !background {
		exitOnSignal(p, c)
	}
}

// ProjectRun runs a given command within a service's container.
func ProjectRun(p *project.Project, c *cobra.Command) {
	if len(c.Flags().Args()) == 1 {
		logrus.Fatal("No service specified")
	}

	serviceName := c.Flags().Args()[0]
	commandParts := c.Flags().Args()[1:]

	if !p.Configs.Has(serviceName) {
		logrus.Fatalf("%s is not defined in the template", serviceName)
	}

	exitCode, err := p.Run(serviceName, commandParts)
	if err != nil {
		logrus.Fatal(err)
	}

	os.Exit(exitCode)
}

// ProjectStart starts services.
func ProjectStart(p *project.Project, c *cobra.Command) {
	err := p.Start(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectRestart restarts services.
func ProjectRestart(p *project.Project, c *cobra.Command) {
	err := p.Restart(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectLog gets services logs.
func ProjectLog(p *project.Project, c *cobra.Command) {
	err := p.Log(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
	if follow, _ := c.Flags().GetBool("follow"); follow {
		wait()
	}
}

// ProjectPull pulls images for services.
func ProjectPull(p *project.Project, c *cobra.Command) {
	err := p.Pull(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectDelete deletes services.
func ProjectDelete(p *project.Project, c *cobra.Command) {
	stoppedContainers, err := p.ListStoppedContainers(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
	if len(stoppedContainers) == 0 {
		fmt.Println("No stopped containers")
		return
	}
	force, _ := c.Flags().GetBool("force")
	if !force {
		fmt.Printf("Going to remove %v\nAre you sure? [yN]\n", strings.Join(stoppedContainers, ", "))
		var answer string
		_, err := fmt.Scanln(&answer)
		if err != nil {
			logrus.Fatal(err)
		}
		if answer != "y" && answer != "Y" {
			return
		}
	}
	err = p.Delete(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectKill forces stop service containers.
func ProjectKill(p *project.Project, c *cobra.Command) {
	err := p.Kill(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectPause pauses service containers.
func ProjectPause(p *project.Project, c *cobra.Command) {
	err := p.Pause(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectUnpause unpauses service containers.
func ProjectUnpause(p *project.Project, c *cobra.Command) {
	err := p.Unpause(c.Flags().Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

// ProjectScale scales services.
func ProjectScale(p *project.Project, c *cobra.Command) {
	// This code is a bit verbose but I wanted to parse everything up front
	order := make([]string, 0, 0)
	serviceScale := make(map[string]int)
	services := make(map[string]project.Service)

	for _, arg := range c.Flags().Args() {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) != 2 {
			logrus.Fatalf("Invalid scale parameter: %s", arg)
		}

		name := kv[0]

		count, err := strconv.Atoi(kv[1])
		if err != nil {
			logrus.Fatalf("Invalid scale parameter: %v", err)
		}

		if !p.Configs.Has(name) {
			logrus.Fatalf("%s is not defined in the template", name)
		}

		service, err := p.CreateService(name)
		if err != nil {
			logrus.Fatalf("Failed to lookup service: %s: %v", service, err)
		}

		order = append(order, name)
		serviceScale[name] = count
		services[name] = service
	}

	for _, name := range order {
		scale := serviceScale[name]
		logrus.Infof("Setting scale %s=%d...", name, scale)
		err := services[name].Scale(scale)
		if err != nil {
			logrus.Fatalf("Failed to set the scale %s=%d: %v", name, scale, err)
		}
	}
}

func wait() {
	<-make(chan interface{})
}

func exitOnSignal(p *project.Project, c *cobra.Command) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range signalChan {
			fmt.Printf("\nGracefully stopping...\n")
			ProjectStop(p, c)
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
