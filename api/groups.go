package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"github.com/ONSdigital/dp-identity-api/v2/models"
	dplogs "github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	GroupsCreatePermission = "groups:create"
	GroupsReadPermission   = "groups:read"
	GroupsEditPermission   = "groups:update"
	GroupsDeletePermission = "groups:delete"
	sortOrderAsc           = "asc"
	sortOrderDesc          = "desc"
)

// CreateGroupHandler creates a new group
func (api *API) CreateGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	createGroup := models.CreateUpdateGroup{}
	err = json.Unmarshal(body, &createGroup)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	// no groupname in body, set UUID
	if createGroup.ID == nil {
		UUID := uuid.NewString()
		createGroup.ID = &UUID
	}

	createGroup.GroupsList, err = api.GetListGroups(ctx)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListGroups request from update group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	validationErrs := createGroup.ValidateCreateUpdateGroupRequest(ctx, true)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	// build create group input
	input := createGroup.BuildCreateGroupInput(&api.UserPoolID)
	_, err = api.CognitoClient.CreateGroup(ctx, input)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito CreateGroup request from create a new group endpoint")
		if cognitoErr.Code == models.GroupExistsError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	jsonResponse, responseErr := createGroup.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return createGroup.NewSuccessResponse(jsonResponse, http.StatusCreated, nil), nil
}

// UpdateGroupHandler update group details for a given group by id (GroupName)
func (api *API) UpdateGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)

	id := vars["id"]
	updateGroup := models.CreateUpdateGroup{
		ID: &id,
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	err = json.Unmarshal(body, &updateGroup)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	validationErrs := updateGroup.ValidateCreateUpdateGroupRequest(ctx, false)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	input := updateGroup.BuildUpdateGroupInput(api.UserPoolID)
	_, err = api.CognitoClient.UpdateGroup(ctx, input)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito UpdateGroup request from update a group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	jsonResponse, responseErr := updateGroup.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return updateGroup.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// AddUserToGroupHandler adds a user to the specified group
func (api *API) AddUserToGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}

	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolID)
	_, err := api.CognitoClient.GetGroup(ctx, groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from Get group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	var bodyJSON map[string]string
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}
	userID := bodyJSON["user_id"]

	validationErrs := group.ValidateAddRemoveUser(ctx, userID)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	response, responseErr := api.AddUserToGroup(ctx, group, userID)

	if responseErr != nil {
		cognitoErr := models.NewCognitoError(ctx, responseErr, "Cognito AddUserToGroup request from add user to group endpoint")
		if cognitoErr.Code == models.UserNotFoundError || cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}
	jsonResponse, responseErr := response.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// ListUsersInGroupHandler list the users in the specified group
func (api *API) ListUsersInGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}

	listUsers, err := api.getUsersInAGroup(ctx, group)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListUsersInGroup request from list users in group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	listOfUsers := models.UsersList{}
	listOfUsers.MapCognitoUsers(&listUsers)

	if err = req.ParseForm(); err != nil {
		dplogs.Error(ctx, "error parsing form", err)
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, err)
	}
	users := listOfUsers.Users
	sortBy := strings.Split(req.Form.Get("sort"), ":")
	sortUsers(ctx, users, sortBy)

	jsonResponse, responseErr := listOfUsers.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func sortUsers(ctx context.Context, users []models.UserParams, sortBy []string) bool {
	if sortBy[0] == "created" || sortBy[0] == "" {
		return true
	}
	switch sortBy[0] {
	case "forename":
		switch sortBy[1] {
		case sortOrderAsc:
			sortByUserNameAsc := func(i, j int) bool {
				return users[i].Forename < users[j].Forename
			}
			sort.Slice(users, sortByUserNameAsc)
			return true
		case sortOrderDesc:
			sortByUserNameDesc := func(i, j int) bool {
				return users[i].Forename > users[j].Forename
			}
			sort.Slice(users, sortByUserNameDesc)
			return true
		default:
			dplogs.Info(ctx, "groups.sortUsers: Not a correct sort by value. Users not sorted.", dplogs.Data{"sortBy": sortBy})
			return false
		}
	default:
		dplogs.Info(ctx, "groups.sortUsers: Not a correct sort by value. Users not sorted.", dplogs.Data{"sortBy": sortBy})
		return false
	}
}

func (api *API) getUsersInAGroup(ctx context.Context, group models.Group) (listOfUsers []types.UserType, err error) {
	firstTimeCheck := false
	var nextToken string
	for !firstTimeCheck || nextToken != "" {
		firstTimeCheck = true

		groupMembersRequest := group.BuildListUsersInGroupRequestWithNextToken(api.UserPoolID, nextToken)
		groupMembersResponse, err := api.CognitoClient.ListUsersInGroup(ctx, groupMembersRequest)
		if err != nil {
			return nil, err
		}

		listOfUsers = append(listOfUsers, groupMembersResponse.Users...)
		nextToken = ""
		if groupMembersResponse.NextToken != nil {
			nextToken = *groupMembersResponse.NextToken
		}
	}
	return listOfUsers, nil
}

// RemoveUserFromGroupHandler adds a user to the specified group
func (api *API) RemoveUserFromGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}
	userID := vars["user_id"]

	validationErrs := group.ValidateAddRemoveUser(ctx, userID)
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
	}

	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolID)
	_, err := api.CognitoClient.GetGroup(ctx, groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from Get group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	response, responseErr := api.RemoveUserFromGroup(ctx, group, userID)
	if responseErr != nil {
		cognitoErr := models.NewCognitoError(ctx, responseErr, "Cognito RemoveUserFromGroupEndpoint request from add user to group endpoint")

		switch cognitoErr.Code {
		case models.UserNotFoundError:
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		case models.NotFoundError:
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		default:
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
		}
	}
	jsonResponse, responseErr := response.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// GetListGroups - List Groups pagination allows first call and then any other call if nextToken is not ""
func (api *API) GetListGroups(ctx context.Context) (*cognitoidentityprovider.ListGroupsOutput, error) {
	firstTimeCheck := false
	var nextToken string
	group := models.ListUserGroupType{}

	listOfGroups := cognitoidentityprovider.ListGroupsOutput{}

	for !firstTimeCheck || nextToken != "" {
		firstTimeCheck = true

		listGroupsRequest := group.BuildListGroupsRequest(api.UserPoolID, nextToken)
		listGroupsResponse, err := api.CognitoClient.ListGroups(ctx, listGroupsRequest)
		if err != nil {
			return nil, err
		}

		listOfGroups.Groups = append(listOfGroups.Groups, listGroupsResponse.Groups...)
		nextToken = ""
		if listGroupsResponse.NextToken != nil {
			nextToken = *listGroupsResponse.NextToken
		}
	}
	return &listOfGroups, nil
}

// ListGroupsHandler lists the users in the user pool
func (api *API) ListGroupsHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	finalGroupsResponse := models.ListUserGroups{}

	listOfGroups, err := api.GetListGroups(ctx)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListofUserGroups request from list user groups endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	if err = req.ParseForm(); err != nil {
		dplogs.Error(ctx, "error parsing form", err)
		return nil, models.NewErrorResponse(http.StatusBadRequest, nil, err)
	}
	query := req.Form.Get("sort")
	if query != "" && query != "created" {
		sortCriteria, err := validateQuery(query)
		if err != nil {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, err)
		}
		if err := sortGroups(listOfGroups, sortCriteria); err != nil {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, err)
		}
	}

	jsonResponse, responseErr := finalGroupsResponse.BuildListGroupsSuccessfulJSONResponse(ctx, listOfGroups)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// GetGroupHandler gets group details for given groups
func (api *API) GetGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}
	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolID)
	groupGetResponse, err := api.CognitoClient.GetGroup(ctx, groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from Get group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	group.MapCognitoDetails(*groupGetResponse.Group)

	jsonResponse, responseErr := group.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// DeleteGroupHandler deletes the group for the given group id
func (api *API) DeleteGroupHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	group := models.Group{ID: vars["id"]}
	groupDeleteRequest := group.BuildDeleteGroupRequest(api.UserPoolID)
	_, err := api.CognitoClient.DeleteGroup(ctx, groupDeleteRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito DeleteGroup request from Delete group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}
	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

// SetGroupUsersHandler adds a user to the specified group
func (api *API) SetGroupUsersHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)

	group := models.Group{ID: vars["id"]}

	groupGetRequest := group.BuildGetGroupRequest(api.UserPoolID)
	_, err := api.CognitoClient.GetGroup(ctx, groupGetRequest)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito GetGroup request from Get group endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, handleBodyReadError(ctx, err)
	}

	var bodyJSON []map[string]string
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}
	listOfUsers := models.UsersList{}

	for _, s1 := range bodyJSON {
		validationErrs := group.ValidateAddRemoveUser(ctx, s1["user_id"])
		if len(validationErrs) != 0 {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, validationErrs...)
		}
		userType := types.UserType{
			Username:   aws.String(s1["user_id"]),
			Enabled:    true,
			UserStatus: types.UserStatusTypeConfirmed,
		}
		listOfUsers.Users = append(listOfUsers.Users, models.UserParams{}.MapCognitoDetails(userType))
	}

	setResponse, setErr := api.SetGroupUsers(ctx, group, listOfUsers)
	if setErr != nil {
		return nil, setErr
	}

	jsonResponse, responseErr := setResponse.BuildSuccessfulJSONResponse(ctx)
	if responseErr != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}
	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

func (api *API) SetGroupUsers(ctx context.Context, group models.Group, users models.UsersList) (*models.UsersList, *models.ErrorResponse) {
	keep := false
	successResponse := &models.UsersList{}

	listUsers, err := api.getUsersInAGroup(ctx, group)
	if err != nil {
		cognitoErr := models.NewCognitoError(ctx, err, "Cognito ListUsersInGroup request from set group membership endpoint")
		if cognitoErr.Code == models.NotFoundError {
			return nil, models.NewErrorResponse(http.StatusBadRequest, nil, cognitoErr)
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, cognitoErr)
	}

	for _, s1 := range listUsers {
		keep = false
		for i := range users.Users {
			s2 := &users.Users[i]
			if s2.ID == *s1.Username {
				keep = true
			}
		}
		if !keep {
			successResponse, _ = api.RemoveUserFromGroup(ctx, group, *s1.Username)
		}
	}

	for i := range users.Users {
		s1 := &users.Users[i] // Use reference instead of copying
		keep = false
		for _, s2 := range listUsers {
			if *s2.Username == s1.ID {
				keep = true
			}
		}
		if !keep {
			successResponse, _ = api.AddUserToGroup(ctx, group, s1.ID)
		}
	}

	return successResponse, nil
}

// AddUserToGroup adds a user to the specified group
func (api *API) AddUserToGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, error) {
	userAddToGroupInput := group.BuildAddUserToGroupRequest(api.UserPoolID, userID)
	_, err := api.CognitoClient.AdminAddUserToGroup(ctx, userAddToGroupInput)
	if err != nil {
		return nil, err
	}

	listUsers, err := api.getUsersInAGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	listOfUsers := models.UsersList{}
	listOfUsers.MapCognitoUsers(&listUsers)

	return &listOfUsers, nil
}

// RemoveUserFromGroup adds a user to the specified group
func (api *API) RemoveUserFromGroup(ctx context.Context, group models.Group, userID string) (*models.UsersList, error) {
	userRemoveFromGroupInput := group.BuildRemoveUserFromGroupRequest(api.UserPoolID, userID)
	_, err := api.CognitoClient.AdminRemoveUserFromGroup(ctx, userRemoveFromGroupInput)
	if err != nil {
		return nil, err
	}

	listUsers, err := api.getUsersInAGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	listOfUsers := models.UsersList{}
	listOfUsers.MapCognitoUsers(&listUsers)

	return &listOfUsers, nil
}

// ListGroupsUsersHandler produces a user requested report of all groups with members including groups that act as roles
// output by default is json but if request header accept == text/csv then the output is csv format
// each line consists of the group description and user email
func (api *API) ListGroupsUsersHandler(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	var (
		GroupsUsersList *[]models.ListGroupUsersType
	)
	listOfGroups, err := api.GetListGroups(ctx)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}
	GroupsUsersList, err = api.GetTeamsReportLines(ctx, listOfGroups)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}

	if req.Header.Get("Accept") == "text/csv" {
		header := map[string]string{"Content-type": "text/csv"}
		return models.NewSuccessResponse(api.ListGroupsUsersCSV(GroupsUsersList).Bytes(), http.StatusOK, header), nil
	}

	jsonResponse, err := json.Marshal(GroupsUsersList)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil
}

// ListGroupsUsersCSV converts the groupsUsersList output to csv
func (api *API) ListGroupsUsersCSV(groupsUsersList *[]models.ListGroupUsersType) *bytes.Buffer {
	var csvHeader = models.ListGroupUsersType{
		GroupName: "Group",
		UserEmail: "User",
	}
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	rows := [][]string{{csvHeader.GroupName, csvHeader.UserEmail}}
	for _, record := range *groupsUsersList {
		rows = append(rows, []string{record.GroupName, record.UserEmail})
	}

	if err := w.WriteAll(rows); err != nil {
		dplogs.Error(context.Background(), "failed to write CSV rows", err)
	}
	return buf
}

// GetTeamsReportLines  from the listOfGroups for each group gets the list of members and produces output
// group description user email for each group member
func (api *API) GetTeamsReportLines(ctx context.Context, listOfGroups *cognitoidentityprovider.ListGroupsOutput) (*[]models.ListGroupUsersType, error) {
	var GroupsUsersList []models.ListGroupUsersType
	for _, ListGroup := range listOfGroups.Groups {
		inputGroup := models.Group{ID: *ListGroup.GroupName}
		listUsers, err := api.getUsersInAGroup(ctx, inputGroup)
		if err != nil {
			return nil, err
		}
		for _, user := range listUsers {
			for _, attribute := range user.Attributes {
				if strings.EqualFold(*attribute.Name, "email") {
					GroupsUsersList = append(GroupsUsersList, models.ListGroupUsersType{
						GroupName: *ListGroup.Description,
						UserEmail: *attribute.Value,
					})
				}
			}
		}
	}

	if GroupsUsersList == nil {
		GroupsUsersList = []models.ListGroupUsersType{}
	}

	return &GroupsUsersList, nil
}

// sortGroups sorts groups in alphabetical order based on the specified sorting criteria
func sortGroups(listGroupOutput *cognitoidentityprovider.ListGroupsOutput, sortBy []string) error {
	groups := listGroupOutput.Groups

	switch {
	case len(sortBy) == 1 && sortBy[0] == "name":
		sortByGroupName(groups, true)
		return nil
	case len(sortBy) == 2 && sortBy[0] == "name":
		switch sortBy[1] {
		case sortOrderAsc:
			sortByGroupName(groups, true)
			return nil
		case sortOrderDesc:
			sortByGroupName(groups, false)
			return nil
		default:
			return fmt.Errorf("incorrect sort value: %v Groups not sorted", sortBy)
		}
	default:
		return fmt.Errorf("incorrect sort value: %v Groups not sorted", sortBy)
	}
}

// sortByGroupName determines the sorting criteria and sorts groups in either ascending or descending order
func sortByGroupName(groups []types.GroupType, ascending bool) {
	sort.Slice(groups, func(i, j int) bool {
		if ascending {
			return strings.ToLower(*groups[i].Description) < strings.ToLower(*groups[j].Description)
		}
		return strings.ToLower(*groups[i].Description) > strings.ToLower(*groups[j].Description)
	})
}

// validateQuery validates the incoming "sort" query string
func validateQuery(query string) ([]string, error) {
	if strings.Contains(query, ":") {
		return strings.Split(query, ":"), nil
	}

	if query == "name" {
		return []string{"name"}, nil
	}

	return nil, fmt.Errorf("invalid query string: %v", query)
}
