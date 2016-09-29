package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

const DOCKER_ROOT = "/var/lib/docker"

func main() {
	var (
		all  bool
		root string
	)
	flag.BoolVar(&all, "a", false, "Display all containers")
	flag.StringVar(&root, "r", DOCKER_ROOT, "Root of the Docker runtime")
	flag.Parse()

	dir, err := ioutil.ReadDir(filepath.Join(root, "containers"))
	if err != nil {
		fmt.Errorf("failed to read dir '%s': %v", root, err)
	}
	fmt.Printf("ID\t\tImage\t\tCreated\t\tcmd\t\tStatus\t\t\tName\n")
	for _, v := range dir {
		id := v.Name()
		c, err := ContainerFromDisk(id, filepath.Join(DOCKER_ROOT, "containers"))
		if err != nil {
			fmt.Errorf("failed load container %s: %v", id, err)
		}
		if c != nil {
			if !all {
				if c.IsRunning() {
					Display(c)
				}
			} else {
				Display(c)
			}
		}
	}
}

func Display(c *Container) {
	cmd := c.Path
	if len(c.Args) != 0 {
		cmd = c.Path + " " + fmt.Sprintf("%s", c.Args)
	}
	if len(cmd) > 12 {
		cmd = cmd[:12]
	}
	fmt.Printf("%s\t%s\t%s\t\t%s\t\t%s\t\t%s\n", c.ID[:12], c.ImageID[7:19], HumanDuration(time.Now().UTC().Sub(c.Created)), cmd, c.String(), c.Name)
}
