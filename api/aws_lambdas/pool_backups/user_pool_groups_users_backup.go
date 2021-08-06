package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
    backupFilePath      = "/tmp/"
    usersGroupsFileName = "users_groups_export_"
    backupFileName      = "groups_export_2021-08-06_13_27_11.csv"
    userPoolID          = "eu-west-1_Rnma9lp2q"
    fileName            = "groups_export_"
    dateLayout          = "2006-01-02_15_04_05"
    s3Bucket            = "ons-dp-develop-cognito-backup-poc"
    awsRegion           = "eu-west-1"
    localData           = "/tmp/users_groups_data.csv"
)

type cognitoClient struct {
	client *cognitoidentityprovider.CognitoIdentityProvider
}

type usersGroups struct {
    groupName string
    users []*cognitoidentityprovider.UserType
}

type MyEvent struct {
    Name string `json:"name"`
}

func main() {
    lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
    var (
        listUsersGroupsResp, result *cognitoidentityprovider.ListUsersInGroupOutput
        listUsersGroupsInput = buildListUsersGroupsRequest(userPoolID, "", "")
        userGroupsList = []usersGroups{}
        debug, fileDebug bool = false, false
    )

    // read data from backup - only a POC so breaking any modeling not too problematic
    allUsersGroupsDataBackup, err := loadDataFromCSV(ctx)
    if err != nil {
        os.Exit(1)
    }

    // build cognito client
    c := &cognitoClient{
		client: cognitoidentityprovider.New(session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})), &aws.Config{Region: aws.String(awsRegion)}),
	}

    // traverse all users in groups data
    for _, groupData := range allUsersGroupsDataBackup {
        // set group name for request
        listUsersGroupsInput.GroupName = &groupData[0]
        listUsersGroupsInput.NextToken = nil
        listUsersGroupsResp, _ = c.generateListUsersGroupsRequest(listUsersGroupsInput)
        // continue oif no users in group
        if len(listUsersGroupsResp.Users) == 0 {
            continue
        }
        if listUsersGroupsResp.NextToken != nil {
            listUsersGroupsInput.NextToken = listUsersGroupsResp.NextToken
            // set `loadingInProgress` to control requesting new list data
            loadingInProgress := true
            for loadingInProgress {
                result, _ = c.generateListUsersGroupsRequest(listUsersGroupsInput)
                listUsersGroupsResp.Users = append(listUsersGroupsResp.Users, result.Users...)
                if result.NextToken != nil {
                    listUsersGroupsInput.NextToken = result.NextToken
                    break
                } else {
                    loadingInProgress = false
                    break
                }
            }
        }
        // debug
        if debug {
            fmt.Fprintf(os.Stderr, "user count for group: "+groupData[0]+" is : "+strconv.Itoa(len(listUsersGroupsResp.Users)))
        }
        // add to userGroupsList
        userGroupsList = append(
            userGroupsList,
            usersGroups{
                groupName: groupData[0],
                users: listUsersGroupsResp.Users,
            },
        )
    }
    // write data to S3
    fileName := usersGroupsFileName+time.Now().Format(dateLayout)+".csv"
    file, _ := os.Create(backupFilePath+fileName)
 
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, userGroupsData := range userGroupsList {
		// build a string using user names
		var userNamesList string
		for _, user := range userGroupsData.users {
			userNamesList += ","+*user.Username
		}
		userGroupsData := userGroupsData.groupName+userNamesList
		_ = writer.Write(strings.Split(userGroupsData, ","))
	}
	file.Close()

	file, _ = os.Open(backupFilePath+fileName)
	defer file.Close()
	fileStats, _ := file.Stat()
	if fileDebug {
		// get the size etc
		println(fileName+" has a file length of: "+strconv.FormatInt(fileStats.Size(), 10)+" file name: "+fileStats.Name())
		readPrintFile(file)
	}

	// Create session object and S3 service client
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	fileBuffer := make([]byte, fileStats.Size())
    file.Read(fileBuffer)
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key: aws.String(fileName),
		Body: bytes.NewReader(fileBuffer),
		ContentLength: aws.Int64(fileStats.Size()),
		ContentType: aws.String("text/csv"),
	})
	
	// delete file
	os.Remove(backupFilePath+fileName)

	if err != nil {
		// Print the error and exit.
		return fmt.Sprintf("Unable to upload %q to %q, %v", fileName, s3Bucket, err), nil
	}

    return fmt.Sprintf("Request completed. Number of `userGroupsList` records written: ", strconv.Itoa(len(userGroupsList)), "fileName: ", fileName, "s3Bucket: ", s3Bucket), nil
}

// generateListUsersGroupsRequest - local routine to generate users in groups request
func (c cognitoClient) generateListUsersGroupsRequest(input *cognitoidentityprovider.ListUsersInGroupInput) (*cognitoidentityprovider.ListUsersInGroupOutput, error) {
    return c.client.ListUsersInGroup(input)
}

// buildListUsersGroupsRequest builds a correctly populated ListUsersInGroupInput object
func buildListUsersGroupsRequest(userPoolId, nextToken, groupName string) *cognitoidentityprovider.ListUsersInGroupInput {
    usersGroupsInput := &cognitoidentityprovider.ListUsersInGroupInput{
        UserPoolId: &userPoolId,
        GroupName: &groupName,
    }

    if nextToken != "" {
        usersGroupsInput.NextToken = &nextToken     
    }

    return usersGroupsInput
}

// load group data from backup file on disk
func loadDataFromCSV(ctx context.Context) ([][]string, error) {
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(awsRegion)},
    )
    
    downloader := s3manager.NewDownloader(sess)

    file, err := os.Create(localData)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to create file object for %q, %v", backupFileName, err)
        return nil, err
    }
    _, err = downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(s3Bucket),
            Key:    aws.String(backupFileName),
        },
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to download item %q, %v"+"\n", backupFileName, err)
        return nil, err
    }
    defer file.Close()

    csvReader := csv.NewReader(file)
    csvReader.FieldsPerRecord = -1
    
    groupsData, err := csvReader.ReadAll()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to parse <"+fileName+"> as CSV: ", err.Error())
        return nil, err
    }

    return groupsData, nil
}

// print file contents
func readPrintFile(file *os.File) {
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var text []string
  
    for scanner.Scan() {
        text = append(text, scanner.Text())
    }
    for _, each_ln := range text {
        fmt.Println(each_ln)
    }
}

//error handling function
func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}
