package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/confirmation"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/groups"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/users"
	"github.com/ONSdigital/log.go/v2/log"
	"os"
)

func main() {
	confirmation := confirmation.AskForConfirmation()

	if !confirmation {
		os.Exit(0)
	}

	ctx := context.Background()
	config := config.GetConfig()

	fmt.Printf("Config: %+v", config)

	var err error

	err = users.ImportUsersFromS3(ctx, config)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to import group"), err)
	}

	err = groups.ImportGroupsFromS3(ctx, config)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to import group"), err)
	}

	err = groups.ImportGroupsMembersFromS3(ctx, config)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to import group"), err)
	}

}
