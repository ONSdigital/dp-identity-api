package main

import (
	"context"
	"os"

	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/confirmation"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/groups"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/users"
	"github.com/ONSdigital/log.go/v2/log"
)

func main() {
	confirmationResponse := confirmation.AskForConfirmation()

	if !confirmationResponse {
		os.Exit(0)
	}

	ctx := context.Background()
	cfg := config.GetConfig()

	log.Info(ctx, "print out log", log.Data{"config": cfg})

	var err error

	err = users.ImportUsersFromS3(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to import group", err)
	}

	err = groups.ImportGroupsFromS3(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to import group", err)
	}

	err = groups.ImportGroupsMembersFromS3(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to import group", err)
	}
}
