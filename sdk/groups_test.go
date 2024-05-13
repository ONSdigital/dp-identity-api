package sdk

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetGroups(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupsEndpoint := "/groups"
	groupsBody := `{
		"groups":[
		   {
			"id":"1",
			"name":"Group A",
			"precedence":1,
			"created":"2024-05-13T12:00:00Z"
		   },
		   {
			"id":"1",
			"name":"Group B",
			"precedence":1,
			"created":"2024-05-13T12:00:00Z"
		   },
		   {
			"id":"1",
			"name":"Group C",
			"precedence":1,
			"created":"2024-05-13T12:00:00Z"
		   }
		],
		"count":1
	 }`

	Convey("Given GetGroups is called successfully", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(groupsBody))),
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroups is called", func() {
			groupsResponse, err := identityAPIClient.GetGroups(ctx, nil)
			So(err, ShouldBeNil)

			Convey("Then the expected groups response is returned", func() {
				expectedGroupsResponse := GroupsResponse{
					Groups: []Group{
						{
							ID:         "1",
							Name:       "Group A",
							Precedence: 1,
							Created:    time.Date(2024, 05, 13, 12, 00, 00, 00, time.UTC),
						},
						{
							ID:         "1",
							Name:       "Group B",
							Precedence: 1,
							Created:    time.Date(2024, 05, 13, 12, 00, 00, 00, time.UTC),
						},
						{
							ID:         "1",
							Name:       "Group C",
							Precedence: 1,
							Created:    time.Date(2024, 05, 13, 12, 00, 00, 00, time.UTC),
						},
					},
					Count: 1,
				}
				So(groupsResponse, ShouldResemble, &expectedGroupsResponse)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, groupsEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
					})
				})
			})
		})
	})

	Convey("Given GetGroups returns an error", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroups is called", func() {
			_, err := identityAPIClient.GetGroups(ctx, nil)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And client.Do should be called once with the expected parameters", func() {
					doCalls := httpClient.DoCalls()
					So(doCalls, ShouldHaveLength, 1)
					So(doCalls[0].Req.URL.Path, ShouldEqual, groupsEndpoint)
					So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
				})
			})
		})
	})
}

func TestGetGroup(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	groupBody := `{"group":{"id":"1","name":"Group A","precedence":1,"created":"2024-05-13T12:00:00Z"}}`
	groupID := "1"
	groupsEndpoint := "/groups/" + groupID

	Convey("Given GetGroup is called successfully", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(groupBody))),
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroup is called", func() {
			groupResponse, err := identityAPIClient.GetGroup(ctx, groupID)
			So(err, ShouldBeNil)

			Convey("Then the expected group response is returned", func() {
				expectedGroupResponse := Group{
					ID:         "1",
					Name:       "Group A",
					Precedence: 1,
					Created:    time.Date(2024, 05, 13, 12, 00, 00, 00, time.UTC),
				}
				So(groupResponse, ShouldResemble, &expectedGroupResponse)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, groupsEndpoint)
						So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
					})
				})
			})
		})
	})

	Convey("Given GetGroup returns an error", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			nil)

		identityAPIClient := newIdentityAPIClient(t, httpClient)

		Convey("When GetGroup is called", func() {
			_, err := identityAPIClient.GetGroup(ctx, groupID)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And client.Do should be called once with the expected parameters", func() {
					doCalls := httpClient.DoCalls()
					So(doCalls, ShouldHaveLength, 1)
					So(doCalls[0].Req.URL.Path, ShouldEqual, groupsEndpoint)
					So(doCalls[0].Req.Method, ShouldEqual, http.MethodGet)
				})
			})
		})
	})
}
