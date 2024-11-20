package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	apiError "github.com/ONSdigital/dp-identity-api/v2/sdk/errors"
)

type sortParam string

var (
	SortCreated  sortParam = "created"
	SortName     sortParam = "name"
	SortNameAsc  sortParam = "name:asc"
	SortNameDesc sortParam = "name:desc"
)

type GroupsResponse struct {
	Groups []Group
	Count  int
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type Group struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Precedence int       `json:"precedence"`
	Created    time.Time `json:"created"`
}

// GetGroups gets a list of groups
func (cli *Client) GetGroups(ctx context.Context, sort *sortParam) (*GroupsResponse, apiError.Error) {
	path := fmt.Sprintf("%s/groups", cli.hcCli.URL)
	if sort != nil {
		path += "?sort=" + url.QueryEscape(string(*sort))
	}

	respInfo, apiErr := cli.callIdentityAPI(ctx, path, http.MethodGet, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var groupsResponse GroupsResponse

	if err := json.Unmarshal(respInfo.Body, &groupsResponse); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal groupsResponse - error is: %v", err),
		}
	}

	return &groupsResponse, nil
}

// GetGroup gets a single group by its ID
func (cli *Client) GetGroup(ctx context.Context, id string) (*GroupsResponse, apiError.Error) {
	path := fmt.Sprintf("%s/groups/%s", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callIdentityAPI(ctx, path, http.MethodGet, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var groupResponse GroupResponse

	if err := json.Unmarshal(respInfo.Body, &groupResponse); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal groupResponse - error is: %v", err),
		}
	}

	groupsResponse := &GroupsResponse{
		Groups: []Group{groupResponse.Group},
		Count:  1,
	}
	return groupsResponse, nil
}
