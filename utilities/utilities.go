package utilities

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"os"

	"github.com/ONSdigital/log.go/log"
)

func ComputeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// load group data from backup file on disk
func LoadDataFromCSV(ctx context.Context, fileName string) ([][]string, error) {
    f, err := os.Open(fileName);
	if err != nil {
        log.Event(ctx, "Unable to read input file <"+fileName+">", log.Error(err), log.ERROR)
		return nil, err
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	
    groupsData, err := csvReader.ReadAll()
	if err != nil {
        log.Event(ctx, "Unable to parse <"+fileName+"> as CSV", log.Error(err), log.ERROR)
		return nil, err
    }

    return groupsData, nil
}
