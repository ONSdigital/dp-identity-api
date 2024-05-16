package sdk

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetGroupsReport(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupsReportEndpoint := "/groups-report"
	groupsReportBody := `[
		{
			"group": "Group A",
			"user": "user1@ons.gov.uk"
		},
		{
			"group": "Group B",
			"user": "user2@ons.gov.uk"
		},
		{
			"group": "Group C",
			"user": "user3@ons.gov.uk"
		}
		]`

	Convey("Given GetGroupReports is called successfully", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(groupsReportBody))),
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroupReports is called", func() {
			groupsReportResponse, err := identityAPIClient.GetGroupsReport(ctx)
			So(err, ShouldBeNil)

			Convey("Then the expected groupsreport response is returned", func() {
				expectedGroupsReportResponse := GroupsReportResponse{
					{
						Group: "Group A",
						User:  "user1@ons.gov.uk",
					},
					{
						Group: "Group B",
						User:  "user2@ons.gov.uk",
					},
					{
						Group: "Group C",
						User:  "user3@ons.gov.uk",
					},
				}
				So(groupsReportResponse, ShouldNotBeNil)
				So(groupsReportResponse, ShouldResemble, &expectedGroupsReportResponse)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, groupsReportEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
					})
				})
			})
		})
	})

	Convey("Given GetGroupsReport returns an error", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroupsReport is called", func() {
			_, err := identityAPIClient.GetGroupsReport(ctx)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And client.Do should be called once with the expected parameters", func() {
					doCalls := httpClient.DoCalls()
					So(doCalls, ShouldHaveLength, 1)
					So(doCalls[0].Req.URL.Path, ShouldEqual, groupsReportEndpoint)
					So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
				})
			})
		})
	})
}
