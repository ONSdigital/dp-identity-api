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
                    "created": "2010-01-01T00:00:00Z",
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

#   Remove user from group scenarios
    Scenario: DELETE /v1/groups/{id}/members/{user_id} and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And there are "1" users in group "test-group"
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name": "test-group",
                    "description": "A test group",
                    "precedence": 100,
                    "created": "2010-01-01T00:00:00Z",
                    "members": []
                }
            """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} and checking the response status 200 with other members listed
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And a user with username "efgh5678" and email "other-email@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And user "efgh5678" is a member of group "test-group"
        And there are "2" users in group "test-group"
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name": "test-group",
                    "description": "A test group",
                    "precedence": 100,
                    "created": "2010-01-01T00:00:00Z",
                    "members": [
                        {
                            "id": "efgh5678",
                            "forename": "Bob",
                            "lastname": "Smith",
                            "email": "other-email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ]
                }
            """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} remove user from group, group not found returns 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/test-group/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} remove user from group, user not found returns 400
        Given group "test-group" exists in the database
        And there are "0" users in group "test-group"
        When I DELETE "/v1/groups/test-group/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} remove user from group, internal server error returns 500
        Given group "internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/internal-error/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, internal server error returns 500
        Given group "get-group-internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/get-group-internal-error/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, group not found returns 500
        Given group "get-group-not-found" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/get-group-not-found/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, internal server error returns 500
        Given group "list-group-users-internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/list-group-users-internal-error/members/abcd1234"
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, group not found returns 500
        Given group "list-group-users-not-found" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I DELETE "/v1/groups/list-group-users-not-found/members/abcd1234"
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
#   Get users from group scenarios        
    Scenario: GET /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And there are "1" users in group "test-group"
        When I GET "/v1/groups/test-group/members"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
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
                    ],
                    "count": 1,
                    "PaginationToken":""
                }
            """

    Scenario: GET /v1/groups/{id}/members, group not found returns 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I GET "/v1/groups/test-group/members"
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
 Scenario: GET /v1/groups/{id}/members, internal server error returns 500
        Given group "internal-error" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I GET "/v1/groups/internal-error/members"
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