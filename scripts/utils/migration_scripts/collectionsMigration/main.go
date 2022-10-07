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
	"time"
)

type config struct {
	teamsDir,
	collectionDir,
	collectionCopyDir string
	migration bool
}

type teamCollection struct {
	id   string
	name string
}

type collection struct {
	InProgressUris        []interface{} `json:"inProgressUris"`
	CompleteUris          []interface{} `json:"completeUris"`
	ReviewedUris          []interface{} `json:"reviewedUris"`
	PendingDeletes        []interface{} `json:"pendingDeletes"`
	PublishResults        []interface{} `json:"publishResults"`
	TimeseriesImportFiles []interface{} `json:"timeseriesImportFiles"`
	PublishTransactionIds struct {
	} `json:"publishTransactionIds"`
	EventsByUri struct {
	} `json:"eventsByUri"`
	Datasets        []interface{} `json:"datasets"`
	Interactives    []interface{} `json:"interactives"`
	DatasetVersions []interface{} `json:"datasetVersions"`
	ApprovalStatus  string        `json:"approvalStatus"`
	PublishComplete bool          `json:"publishComplete"`
	IsEncrypted     bool          `json:"isEncrypted"`
	Events          []struct {
		Date  time.Time `json:"date"`
		Type  string    `json:"type"`
		Email string    `json:"email"`
	} `json:"events"`
	Id    string   `json:"id"`
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Teams []string `json:"teams"`
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

func getTeamsList(dir string) []teamCollection {
	var teamsList []teamCollection
	teamFiles, path, err := getDirList(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range teamFiles {
		teamFile, _ := ioutil.ReadFile(filepath.Join(path, file.Name()))
		var result map[string]interface{}
		_ = json.Unmarshal(teamFile, &result)
		var team teamCollection
		f, ok := result["id"].(string)
		if ok {
			team = teamCollection{
				id:   f,
				name: result["name"].(string),
			}
		} else {
			team = teamCollection{
				id:   fmt.Sprintf("%.0f", result["id"]),
				name: result["name"].(string),
			}
		}
		teamsList = append(teamsList, team)

	}
	return teamsList
}

func amendCollectionTeams(dir, copyDir string, teamsList []teamCollection) {
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

		var result collection
		_ = json.Unmarshal(collectionFile, &result)
		for ind, collectionTeam := range result.Teams {
			for _, team := range teamsList {
				if team.name == collectionTeam {
					result.Teams[ind] = team.id
					*saveFile = true
				}
			}
		}
		if *saveFile {
			fmt.Println("File Amended ", filepath.Join(dir, file.Name()))
			amendedFile, _ := json.Marshal(result)
			_ = ioutil.WriteFile(filepath.Join(dir, file.Name()), amendedFile, 0644)
		}

	}
}

func revertCollectionTeams(dir, copyDir string, teamsList []teamCollection) {
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
		_ = ioutil.WriteFile(filepath.Join(copyDir, file.Name()), collectionFile, 0644)

		var result collection
		_ = json.Unmarshal(collectionFile, &result)
		for ind, collectionTeam := range result.Teams {
			for _, tmpTeam := range teamsList {
				if tmpTeam.id == collectionTeam {
					result.Teams[ind] = tmpTeam.name
					*saveFile = true
				}
			}
		}
		if *saveFile {
			fmt.Println("File Reverted ", filepath.Join(dir, file.Name()))
			revertFile, _ := json.Marshal(result)
			_ = ioutil.WriteFile(filepath.Join(dir, file.Name()), revertFile, 0644)
		}
	}
}

func main() {
	conf := readConfig()
	teamsList := getTeamsList(conf.teamsDir)
	for _, v := range teamsList {
		fmt.Println(v)
	}
	if conf.migration {
		amendCollectionTeams(conf.collectionDir, conf.collectionCopyDir, teamsList)
	} else {
		revertCollectionTeams(conf.collectionDir, conf.collectionCopyDir, teamsList)
	}
}
