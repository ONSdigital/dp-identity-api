package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/log.go/log"
)

// A POC groups recovery golang script
// Restores user pool groups from backup csv
// Built around golang's `exec` package

const (
	filePath     = "/Users/mk-dv-m0047/workspace/test/"
	fileName     = "groups_export_2021-08-06_13_12_49.csv"
)

func main() {
	ctx := context.Background()

	allGroupsDataBackup, err := utilities.LoadDataFromCSV(ctx, filePath+fileName)
	if err != nil {
		log.Event(ctx, "fatal runtime error", log.Error(err), log.FATAL)
		os.Exit(1)
	}
	//list all groups for user pool
	listGroups(ctx, allGroupsDataBackup)
	// listUsersInGroups(ctx, allGroupsDataBackup)
	// createGroup(ctx)
}

// list group data from aws cognito using data from backup on disk
func listGroups(ctx context.Context, allGroupsDataBackup [][]string) {
	for _, groupData := range allGroupsDataBackup {
		cmd := exec.Command(
			"aws",
			"cognito-idp",
			"get-group",
			"--user-pool-id",
			groupData[5],
			"--group-name",
			groupData[0],
		)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			log.Event(ctx, "Error returned from aws congnito: ", log.Error(err), log.ERROR)
		}
		println("GroupName -> "+groupData[0]+" : group data => "+out.String())
	}
}

// create a new group - using dummy data for POC
func createGroup(ctx context.Context) {
	cmd := exec.Command(
		"aws",
		"cognito-idp",
		"create-group",
		"--user-pool-id",
		"eu-west-1_Rnma9lp2q",
		"--group-name",
		"backup-test-group-1",
	)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Event(ctx, "Error returned from aws congnito: ", log.Error(err), log.ERROR)
	}
	println("GroupName -> `backup-test-group-1` : group data => "+out.String())
}

// list users in groups data from aws cognito using data from backup on disk
func listUsersInGroups(ctx context.Context, allGroupsDataBackup [][]string) {
	for _, usersInGroupsData := range allGroupsDataBackup {
		cmd := exec.Command(
			"aws",
			"cognito-idp",
			"list-users-in-group",
			"--user-pool-id",
			usersInGroupsData[5],
			"--group-name",
			usersInGroupsData[0],
		)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			log.Event(ctx, "Error returned from aws congnito: ", log.Error(err), log.ERROR)
		}
		println("GroupName -> "+usersInGroupsData[0]+" : users in group data => "+out.String())
	}
}