package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"text/template"


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
		}
		default:
			fmt.Println("Not implemented", cmd.Rpc)
		}
		if cmd.Result != "" {
			vars[cmd.Result] = "x"
		}
	}
}