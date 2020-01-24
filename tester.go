package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
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

func dockerExec(id string, _ []string) (string, error) {
	//osCmd := exec.Command("docker", cmd...)
	containerID := "1fbb79d1d767" 
	datadir := "/root/.ethereum/goerli"
	quotedJs := "'admin.nodeInfo'"
	args := []string{"exec", "-it", containerID, "geth", "--datadir", datadir, "attach", "--exec", quotedJs}
	//args := []string{"exec", "-t", containerID, "ls"}
	osCmd := exec.Command("docker", args...)
	out, err := osCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func main() {
	testPath := flag.String("t", "test.yml", "Full path to the test yaml file")
	flag.Parse()
	data, _ := ioutil.ReadFile(*testPath)
	var cl commandList
	yaml.Unmarshal(data, &cl)	
	vars := make(map[string]string)
	for _, cmd := range cl.Execute[0:1] {
		switch {
		case cmd.Js != "": {
			tmpl, _ := template.New("test").Parse(cmd.Js)
			var b bytes.Buffer
			tmpl.Execute(&b, vars)
			js := b.String()
			fmt.Println("exec js", js)
			containerID := dockerID(cmd.Node)
			datadir := "/root/.ethereum/goerli"
			quotedJs := "'" + js + "'"
			args := []string{"exec", "-t", containerID, "geth", "--datadir", datadir, "attach", "--exec", quotedJs}
			result, err := dockerExec(containerID, args)
			if err != nil {
				fmt.Println("dockerExec", err)
			}
			fmt.Println("result", result)
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
