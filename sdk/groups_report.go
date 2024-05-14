package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	apiError "github.com/ONSdigital/dp-identity-api/sdk/errors"
)

type GroupsReportResponse struct {
	Group string `json:"group"`
	User  string `json:"user"`
}

func (cli *Client) GetGroupsReport(ctx context.Context) (*[]GroupsReportResponse, apiError.Error) {
	path := fmt.Sprintf("%s/groups-report", cli.hcCli.URL)

	respInfo, apiErr := cli.callIdentityAPI(ctx, path, http.MethodGet, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var groupsReportResponse []GroupsReportResponse

	if err := json.Unmarshal(respInfo.Body, &groupsReportResponse); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal groupsReportResponse - error is: %v", err),
		}
	}

	return &groupsReportResponse, nil
}
