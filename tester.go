package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"text/template"


	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type command struct {
	Node string
	Rpc string
	Js string
	Result string
}

type commandList struct {
	Execute []command
}

var dockerIDCache = make(map[string]string)

func dockerID(name string) string {
	if id, ok := dockerIDCache[name]; ok {
		return id
	}
	ctx := context.Background()
	opts := types.ContainerListOptions{All: true}
	opts.Filters = filters.NewArgs()
	opts.Filters.Add("name", name)
	dockerCli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	containers, err := dockerCli.ContainerList(ctx, opts)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if len(containers) != 1 {
		return ""
	}
	id := containers[0].ID
	dockerIDCache[name] = id
	return id
}

func main() {
	testPath := flag.String("t", "test.yml", "Full path to the test yaml file")
	flag.Parse()
	data, _ := ioutil.ReadFile(*testPath)
	var cl commandList
	yaml.Unmarshal(data, &cl)	
	vars := make(map[string]string)
	for _, cmd := range cl.Execute {
		switch {
		case cmd.Js != "": {
			tmpl, _ := template.New("test").Parse(cmd.Js)
			var b bytes.Buffer
			tmpl.Execute(&b, vars)
			js := b.String()
			fmt.Println("exec js", js)
			containerID := dockerID(cmd.Node)
			fmt.Println(containerID)
		}
		default:
			fmt.Println("Not implemented", cmd.Rpc)
		}
		if cmd.Result != "" {
			vars[cmd.Result] = "x"
		}
	}
	fmt.Println("vars", vars)
}
