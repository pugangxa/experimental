/*
Copyright 2021 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"regexp"

	"github.com/google/cel-go/cel"
)

const (
	pluginsDir = "/var/plugins"
	celLibName = "CustomLib"
)

func listFiles(dir, pattern string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, err
		}
		if matched {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles, nil
}

func loadPlugins(env *cel.Env) (*cel.Env, error) {
	if _, err := os.Stat(pluginsDir); err != nil {
		return nil, err
	}

	plugins, err := listFiles(pluginsDir, `.*_plugin.so`)
	if err != nil {
		return nil, err
	}

	for _, celPlugin := range plugins {
		plug, err := plugin.Open(path.Join(pluginsDir, celPlugin.Name()))
		if err != nil {
			fmt.Printf("failed to open plugin %s: %v\n", celPlugin.Name(), err)
			continue
		}
		celLibSymbol, err := plug.Lookup(celLibName)
		if err != nil {
			fmt.Printf("plugin %s does not export symbol \"%s\"\n",
				celPlugin.Name(), celLibName)
			continue
		}
		celLib, ok := celLibSymbol.(cel.Library)
		if !ok {
			fmt.Printf("Symbol %s (from %s) does not implement cel.Library interface\n",
				celLibName, celPlugin.Name())
			continue
		}
		env, _ = env.Extend(cel.Lib(celLib))
	}
	return env, nil
}
