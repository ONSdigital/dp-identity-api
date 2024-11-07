package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ONSdigital/dp-identity-api/jwks"
)

func main() {
	jwksHandler := &jwks.JWKS{}

	userPoolID, ok := os.LookupEnv("USER_POOL_ID")
	if !ok {
		fmt.Println("ensure the USER_POOL_ID environment variable is set.")
		os.Exit(2)
	}

	REGION, ok := os.LookupEnv("REGION")
	if !ok {
		fmt.Println("ensure the REGION environment variable is set.")
		os.Exit(2)
	}

	jwksRSAKeys, err := jwksHandler.GetJWKSRSAKeys(REGION, userPoolID)
	if err != nil {
		log.Fatal("could not retrieve the JWKS RSA public keys", err)
	}

	keys := []string{}
	for k, v := range jwksRSAKeys {
		key := fmt.Sprintf("%s:%s", k, v)
		keys = append(keys, key)
	}
	pubKeyString := strings.Join(keys, ",")
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, pubKeyString)
	fmt.Println(colored)
}
