package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type config struct {
	teamsDir,
	collectionDir,
	collectionCopyDir string
	migration bool
}

type teamFile struct {
	id   string
	name string
}

func readConfig() *config {
	conf := &config{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "teamsDir":
			missingVariables("teamsDir", pair[1])
			conf.teamsDir = pair[1]
		case "collectionDir":
			missingVariables("collectionDir", pair[1])
			conf.collectionDir = pair[1]
		case "collectionCopyDir":
			missingVariables("collectionCopyDir", pair[1])
			conf.collectionCopyDir = pair[1]
		case "migration":
			missingVariables("migration", pair[1])
			boolValue, err := strconv.ParseBool(pair[1])
			if err != nil {
				log.Fatal(err)
			}
			conf.migration = boolValue
		}
	}
	return conf
}

func missingVariables(envValue string, value string) bool {
	if len(value) == 0 {
		log.Fatalf("Please set Environment Variable: %s", envValue)
	}
	return true
}

func getDirList(dir string) ([]fs.FileInfo, string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return nil, "", err
	}

	tmpdirFilesList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, path, err
	}
	var dirFilesList []fs.FileInfo

	for _, f := range tmpdirFilesList {
		if strings.Contains(f.Name(), ".json") {
			dirFilesList = append(dirFilesList, f)
		}
	}

	return dirFilesList, path, nil
}

func getTeamsList(dir string) (map[string]string, map[string]string) {
	teamNameByID := make(map[string]string)
	teamIDByName := make(map[string]string)

	var result teamFile
	teamFiles, path, err := getDirList(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range teamFiles {
		teamFile, _ := ioutil.ReadFile(filepath.Join(path, file.Name()))
		err = json.Unmarshal(teamFile, &result)
		if err != nil {
			continue
		}
		if result.name != "" && result.id != "" {
			teamIDByName[result.name] = result.id
			teamNameByID[result.id] = result.name
		}
	}
	return teamIDByName, teamNameByID
}

func amendCollectionTeams(dir, copyDir string, teamsList map[string]string) {
	collectionFiles, path, err := getDirList(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range collectionFiles {
		var saveFile = new(bool)
		*saveFile = false
		collectionFile, _ := ioutil.ReadFile(filepath.Join(path, file.Name()))

		//copy collectionFile to copy dir
		if err := os.Mkdir(copyDir, 0777); err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		_ = ioutil.WriteFile(filepath.Join(copyDir, file.Name()), collectionFile, 0777)

		var result map[string]interface{}
		_ = json.Unmarshal(collectionFile, &result)
		if teams, ok := result["teams"].([]interface{}); ok {
			for ind, team := range teams {
				if val, ok := teamsList[team.(string)]; ok {
					teams[ind] = val
					*saveFile = true
				}
			}
		}

		if *saveFile {
			fmt.Println("File Amended ", file.Name())
			amendedFile, _ := json.Marshal(result)
			_ = ioutil.WriteFile(filepath.Join(dir, file.Name()), amendedFile, 0644)
		}
	}

}

func main() {
	conf := readConfig()
	teamIDByName, teamNameByID := getTeamsList(conf.teamsDir)
	if conf.migration {
		amendCollectionTeams(conf.collectionDir, conf.collectionCopyDir, teamIDByName)
	} else {
		amendCollectionTeams(conf.collectionDir, conf.collectionCopyDir, teamNameByID)
	}
}
