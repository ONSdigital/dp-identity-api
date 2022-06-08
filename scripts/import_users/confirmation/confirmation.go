package confirmation

import (
	"fmt"
	"github.com/ONSdigital/dp-identity-api/scripts/import_users/config"
	"strings"
)

func AskForConfirmation() bool {
	configs := config.GetConfig()

	groupFilename := configs.GroupsFilename
	userFilename := configs.UserFileName
	groupUserFilename := configs.GroupUsersFilename

	var s string

	fmt.Printf("Importing groups from %s \n", groupFilename)
	fmt.Printf("Importing users from %s \n", userFilename)
	fmt.Printf("Importing group users from %s \n", groupUserFilename)

	fmt.Printf("If everything is correct please proceed with (y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
