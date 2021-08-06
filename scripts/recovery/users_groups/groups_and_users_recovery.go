package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/ONSdigital/dp-identity-api/utilities"
	"github.com/ONSdigital/log.go/log"
)

// A POC users and groups recovery golang script
// Restores user pool groups and theirs users from backup csv
// Built around golang's `exec` package

const (
	usersGroupsFilePath = "/Users/mk-dv-m0047/workspace/test/"
	usersGroupsFileName = "users_groups_export_2021-08-09_15_47_01.csv"
)

func main() {
	ctx := context.Background()

	allUsersGroupsDataBackup, err := utilities.LoadDataFromCSV(ctx, usersGroupsFilePath+usersGroupsFileName)
	if err != nil {
		log.Event(ctx, "fatal runtime error", log.Error(err), log.FATAL)
		os.Exit(1)
	}
	//list all groups for user pool
	listUsersAndGroups(ctx, allUsersGroupsDataBackup)
	listUsersInGroups(ctx, allUsersGroupsDataBackup)
	createGroup(ctx)
}

// list groups and user data from aws cognito using data from backup on disk
func listUsersAndGroups(ctx context.Context, allUsersGroupsDataBackup [][]string) {
	for _, userGroupData := range allUsersGroupsDataBackup {
		println("GroupName: "+userGroupData[0]+" | users: "+strconv.Itoa(len(userGroupData[1:])))
		i := 0
		for _, userName := range userGroupData[1:] {
			println("["+strconv.Itoa(i)+"]    Username: "+userName)
			i++
		}
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
