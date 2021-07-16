Feature: Groups

#   Add user to group scenarios
    Scenario: POST /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
        And there are "0" users in group "test-group"
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name": "test-group",
                    "description": "A test group",
                    "precedence": 100,
                    "members": [
                        {
                            "id": "abcd1234",
                            "forename": "Bob",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members with no user Id submitted and checking the response status 400
        Given group "test-group" exists in the database
        And there are "0" users in group "test-group"
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": ""
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code": "InvalidUserId",
                            "description": "the user id was missing"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members add user to group, group not found returns 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code": "NotFound",
                            "description": "the group could not be found"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members add user to group, user not found returns 400
        Given group "test-group" exists in the database
        And there are "0" users in group "test-group"
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code": "UserNotFound",
                            "description": "the user could not be found"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members add user to group, internal server error returns 500
        Given group "internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/internal-error/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "500":
            """
                {
                    "errors": [
                        {
                            "code": "InternalServerError",
                            "description": "Something went wrong"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members get group, internal server error returns 500
        Given group "get-group-internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/get-group-internal-error/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "500":
            """
                {
                    "errors": [
                        {
                            "code": "InternalServerError",
                            "description": "Something went wrong"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members get group, group not found returns 500
        Given group "get-group-not-found" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/get-group-not-found/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "500":
            """
                {
                    "errors": [
                        {
                            "code": "NotFound",
                            "description": "get group - group not found"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members get group, internal server error returns 500
        Given group "list-group-users-internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/list-group-users-internal-error/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "500":
            """
                {
                    "errors": [
                        {
                            "code": "InternalServerError",
                            "description": "Something went wrong"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups/{id}/members get group, group not found returns 500
        Given group "list-group-users-not-found" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/list-group-users-not-found/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "500":
            """
                {
                    "errors": [
                        {
                            "code": "NotFound",
                            "description": "list members - group not found"
                        }
                    ]
                }
            """
