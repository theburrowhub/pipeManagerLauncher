package pipelineparser

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// MixPipelineFiles reads the directory and its subdirectories, merging YAML files into a single map
func MixPipelineFiles(dir string) (error, map[string]interface{}) {
	combinedData := make(map[string]interface{})

	err := readDirRecursively(dir, combinedData)
	if err != nil {
		return err, nil
	}

	return nil, combinedData
}

// readDirRecursively reads the directory and its subdirectories, merging YAML files into combinedData
func readDirRecursively(dir string, combinedData map[string]interface{}) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			err = readDirRecursively(path, combinedData)
			if err != nil {
				return err
			}
		} else if filepath.Ext(entry.Name()) == ".yaml" || filepath.Ext(entry.Name()) == ".yml" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var content map[string]interface{}
			err = yaml.Unmarshal(data, &content)
			if err != nil {
				return err
			}

			mergeMaps(combinedData, content)
		}
	}
	return nil
}

// mergeMaps merges two maps
func mergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if mv, ok := v.(map[string]interface{}); ok {
			if dv, ok := dst[k]; ok {
				if dvMap, ok := dv.(map[string]interface{}); ok {
					mergeMaps(dvMap, mv)
					continue
				}
			}
		}
		dst[k] = v
	}
}
