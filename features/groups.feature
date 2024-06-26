Feature: Groups

#   Create new group scenarios

    Scenario: POST /v1/groups to create group, group created returns 201
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                    "precedence": 49,
                    "id": "123e4567-e89b-12d3-a456-426614174000"
                }
            """
        Then I should receive the following JSON response with status "201":
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                    "precedence": 49,
                    "id": "123e4567-e89b-12d3-a456-426614174000"
                }
            """

    Scenario: POST /v1/groups without a JWT token and checking the response status 401
        When I POST "/v1/groups"
            """
            """
        Then the HTTP status code should be "401"

    Scenario: POST /v1/groups as a publisher user and checking the response status 403
        Given I am a publisher user
        When I POST "/v1/groups"
            """
            """
        Then the HTTP status code should be "403"

    Scenario: POST /v1/groups to create group with no description in request, group created returns 400
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"the group name was not found"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups to create group with no precedence in request, group created returns 400
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $"
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupPrecedence",
                            "description":"the group precedence was not found"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups to create group with reserved pattern in description [lower case], group created returns 400
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "role-Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"a group name cannot start with 'role-' or 'ROLE-'"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups to create group with reserved pattern in description [upper case], group created returns 400
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "ROLE-Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"a group name cannot start with 'role-' or 'ROLE-'"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups to create group group precedence doesn't meet minimum of `3`, returns 400
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "This is a test description",
                    "precedence": 1
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupPrecedence",
                            "description":"the group precedence needs to be a minumum of 10 and maximum of 100"
                        }
                    ]
                }
            """

    Scenario: POST /v1/groups to create group an unexpected 500 error is returned from Cognito
        Given I am an admin user
        When I POST "/v1/groups"
            """
                {
                    "name": "Internal Server Error",
                    "precedence": 12
                }
            """
        Then I should receive the following JSON response with status "500":
            """
               {"code":"InternalServerError", "description":"Internal Server Error"}
            """

#   Update group scenarios

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group, group updated returns 200
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  existing group  $",
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  existing group  $",
                    "precedence": 49,
                    "id": "123e4567-e89b-12d3-a456-426614174000"

                }
            """

    Scenario: PUT /v1/groups/{id} without a JWT token and checking the response status 401
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
            """
        Then the HTTP status code should be "401"

    Scenario: PUT /v1/groups/{id} as a publisher user and checking the response status 403
        Given I am a publisher user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
            """
        Then the HTTP status code should be "403"

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group with no description in request, group update returns 400
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"the group name was not found"
                        }
                    ]
                }
            """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group with no precedence in request, group update returns 200
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  updated group  $"
                }
            """
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name": "Thi$s is a te||st des$%£@^c ription for  updated group  $",
                    "id": "123e4567-e89b-12d3-a456-426614174000",
                    "precedence": null
                }
            """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group with reserved pattern in description [lower case], group update returns 400
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "role-Thi$s is a te||st des$%£@^c ription for  existing group  $",
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"a group name cannot start with 'role-' or 'ROLE-'"
                        }
                    ]
                }
             """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group with reserved pattern in description [upper case], group update returns 400
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "ROLE-Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                    "precedence": 49
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupName",
                            "description":"a group name cannot start with 'role-' or 'ROLE-'"
                        }
                    ]
                }
            """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group group precedence doesn't meet minimum of `10`, returns 400
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "This is a test description",
                    "precedence": 1
                }
            """
        Then I should receive the following JSON response with status "400":
            """
                {
                    "errors": [
                        {
                            "code":"InvalidGroupPrecedence",
                            "description":"the group precedence needs to be a minumum of 10 and maximum of 100"
                        }
                    ]
                }
            """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group an unexpected 500 error is returned from Cognito
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "Internal Server Error",
                    "precedence": 12
                }
            """
        Then I should receive the following JSON response with status "500":
            """
               {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: PUT /v1/groups/123e4567-e89b-12d3-a456-426614174000 to update group a resource not found 404 error is returned
        Given I am an admin user
        When I PUT "/v1/groups/123e4567-e89b-12d3-a456-426614174000"
            """
                {
                    "name": "resource not found",
                    "precedence": 12
                }
            """
        Then I should receive the following JSON response with status "404":
            """
                {
                    "errors": [
                        {
                            "code":"NotFound",
                            "description":"Resource not found"
                        }
                    ]
                }
            """

#   Add user to group scenarios

    Scenario: POST /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
            And there are 0 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
                        {
                            "active": true,
                            "email":  "email@ons.gov.uk",
                            "forename":  "Bob",
                            "groups": [],
                            "id": "abcd1234",
                            "lastname": "Smith",
                            "status": "CONFIRMED",
                            "status_notes": ""
                        }
                    ],
                    "count": 1,
                    "PaginationToken":""
                }
                """

    Scenario: POST /v1/groups/{id}/members without a JWT token and checking the response status 401
        When I POST "/v1/groups/test-group/members"
            """
            """
        Then the HTTP status code should be "401"

    Scenario: POST /v1/groups/{id}/members as a publisher user and checking the response status 403
        Given I am a publisher user
        When I POST "/v1/groups/test-group/members"
            """
            """
        Then the HTTP status code should be "403"

    Scenario: POST /v1/groups/{id}/members with no user Id submitted and checking the response status 400
        Given group "test-group" exists in the database
            And there are 0 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
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
            And I am an admin user
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234"
                }
            """
        Then I should receive the following JSON response with status "404":
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
            And there are 0 users in group "test-group"
            And I am an admin user
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

#   Remove user from group scenarios

    Scenario: DELETE /v1/groups/{id}/members/{user_id} and checking the response status 200
        Given group "test-group" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And user "abcd1234" is a member of group "test-group"
            And there are 1 users in group "test-group"
            And I am an admin user
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [],
                    "count": 0,
                    "PaginationToken":""
                }
            """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} without a JWT token and checking the response status 401
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then the HTTP status code should be "401"

    Scenario: DELETE /v1/groups/{id}/members/{user_id} as a publisher user and checking the response status 403
            Given I am a publisher user
            When I DELETE "/v1/groups/test-group/members/abcd1234"
            Then the HTTP status code should be "403"

    Scenario: DELETE /v1/groups/{id}/members/{user_id} and checking the response status 200 with other members listed
            Given group "test-group" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And a user with username "efgh5678" and email "other-email@ons.gov.uk" exists in the database
            And user "abcd1234" is a member of group "test-group"
            And user "efgh5678" is a member of group "test-group"
            And there are 2 users in group "test-group"
            And I am an admin user
            When I DELETE "/v1/groups/test-group/members/abcd1234"
            Then I should receive the following JSON response with status "200":
                """
                    {
                        "users": [
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
                        ],
                        "count": 1,
                        "PaginationToken":""
                    }
                """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} remove user from group, group not found returns 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then I should receive the following JSON response with status "404":
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

    Scenario: DELETE /v1/groups/{id}/members/{user_id} remove user from group, user not found returns 404
        Given group "test-group" exists in the database
            And there are 0 users in group "test-group"
            And I am an admin user
        When I DELETE "/v1/groups/test-group/members/abcd1234"
        Then I should receive the following JSON response with status "404":
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
             And I am an admin user
        When I DELETE "/v1/groups/internal-error/members/abcd1234"
        Then I should receive the following JSON response with status "500":
            """
                {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, internal server error returns 500
        Given group "get-group-internal-error" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I DELETE "/v1/groups/get-group-internal-error/members/abcd1234"
        Then I should receive the following JSON response with status "500":
            """
                {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, group not found returns 404
        Given group "get-group-not-found" exists in the database
             And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
             And I am an admin user
        When I DELETE "/v1/groups/get-group-not-found/members/abcd1234"
        Then I should receive the following JSON response with status "404":
            """
                {"errors":[{"code":"NotFound", "description":"get group - group not found"}]}
            """

#   Get users from group scenarios
    Scenario: GET /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And user "abcd1234" is a member of group "test-group"
            And there are 1 users in group "test-group"
            And I am an admin user
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

    Scenario: GET /v1/groups/{id}/members?sort=forename:asc and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" and forename "Andreas" exists in the database
        And a user with username "bcde1234" and email "email@ons.gov.uk" and forename "Dimitris" exists in the database
        And a user with username "cdef1234" and email "email@ons.gov.uk" and forename "Zeus" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And user "bcde1234" is a member of group "test-group"
        And user "cdef1234" is a member of group "test-group"
        And there are 3 users in group "test-group"
        And I am an admin user
        When I GET "/v1/groups/test-group/members?sort=forename:asc"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
                        {
                            "id": "abcd1234",
                            "forename": "Andreas",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "bcde1234",
                            "forename": "Dimitris",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "cdef1234",
                            "forename": "Zeus",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ],
                    "count": 3,
                    "PaginationToken":""
                }
            """

    Scenario: GET /v1/groups/{id}/members?sort=forename:desc and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" and forename "Andreas" exists in the database
        And a user with username "bcde1234" and email "email@ons.gov.uk" and forename "Dimitris" exists in the database
        And a user with username "cdef1234" and email "email@ons.gov.uk" and forename "Zeus" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And user "bcde1234" is a member of group "test-group"
        And user "cdef1234" is a member of group "test-group"
        And there are 3 users in group "test-group"
        And I am an admin user
        When I GET "/v1/groups/test-group/members?sort=forename:desc"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
                        {
                            "id": "cdef1234",
                            "forename": "Zeus",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "bcde1234",
                            "forename": "Dimitris",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "abcd1234",
                            "forename": "Andreas",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ],
                    "count": 3,
                    "PaginationToken":""
                }
            """

    Scenario: GET /v1/groups/{id}/members?sort=created and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" and forename "Andreas" exists in the database
        And a user with username "bcde1234" and email "email@ons.gov.uk" and forename "Dimitris" exists in the database
        And a user with username "cdef1234" and email "email@ons.gov.uk" and forename "Zeus" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And user "bcde1234" is a member of group "test-group"
        And user "cdef1234" is a member of group "test-group"
        And there are 3 users in group "test-group"
        And I am an admin user
        When I GET "/v1/groups/test-group/members?sort=created"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
                        {
                            "id": "abcd1234",
                            "forename": "Andreas",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "bcde1234",
                            "forename": "Dimitris",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "cdef1234",
                            "forename": "Zeus",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ],
                    "count": 3,
                    "PaginationToken":""
                }
            """

    Scenario: GET /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" and forename "Andreas" exists in the database
        And a user with username "bcde1234" and email "email@ons.gov.uk" and forename "Dimitris" exists in the database
        And a user with username "cdef1234" and email "email@ons.gov.uk" and forename "Zeus" exists in the database
        And user "abcd1234" is a member of group "test-group"
        And user "bcde1234" is a member of group "test-group"
        And user "cdef1234" is a member of group "test-group"
        And there are 3 users in group "test-group"
        And I am an admin user
        When I GET "/v1/groups/test-group/members"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "users": [
                        {
                            "id": "abcd1234",
                            "forename": "Andreas",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "bcde1234",
                            "forename": "Dimitris",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "id": "cdef1234",
                            "forename": "Zeus",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        }
                    ],
                    "count": 3,
                    "PaginationToken":""
                }
            """

    Scenario: GET /v1/groups/{id}/members without a JWT token and checking the response status 401
        When I GET "/v1/groups/test-group/members"
        Then the HTTP status code should be "401"

    Scenario: GET /v1/groups/{id}/members as a publisher user and checking the response status 403
        Given I am a publisher user
        When I GET "/v1/groups/test-group/members"
        Then the HTTP status code should be "403"

    Scenario: GET /v1/groups/{id}/members, group not found returns 400
        Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
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
        And I am an admin user
    When I GET "/v1/groups/internal-error/members"
    Then I should receive the following JSON response with status "500":
        """
            {"code":"InternalServerError", "description":"Internal Server Error"}
        """

#   Get listgroups scenarios
#   list for no groups found
    Scenario: GET /v1/groups and checking the response status 200
        Given there are 0 groups in the database
            And I am an admin user
        When I GET "/v1/groups"
        Then I should receive the following JSON response with status "200":
                """
                    {
                        "groups":null,
                        "count":0,
                        "next_token":null
                    }
        """

#   list for one groups found
    Scenario: GET /v1/groups and checking the response status 200
        Given there are 2 groups in the database
            And I am an admin user
        When I GET "/v1/groups"
        Then the response code should be 200
            And the response should match the following json for listgroups
                """
                    {
                        "count": 2,
                        "groups": [
                            {
                                "name": "group name description 1",
                                "id": "group_name_1",
                                "precedence": 55
                            },
                                                        {
                                "name": "group name description 2",
                                "id": "group_name_2",
                                "precedence": 55
                            }
                        ],
                        "next_token": null
                    }
                """

    Scenario: GET /v1/groups?sort=name:asc and checking the response status 200
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=name:asc"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "count": 3,
                    "groups": [
                        {
                            "name": "A Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                        {
                            "name": "B Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                                                {
                            "name": "C Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        }
                    ],
                    "next_token": null
                }
            """

    Scenario: GET /v1/groups?sort=name:desc and checking the response status 200
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=name:desc"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "count": 3,
                    "groups": [
                        {
                            "name": "C Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                        {
                            "name": "B Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                                                {
                            "name": "A Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        }
                    ],
                    "next_token": null
                }
            """

    Scenario: GET /v1/groups?sort=name and checking the response status 200
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=name"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "count": 3,
                    "groups": [
                        {
                            "name": "A Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                        {
                            "name": "B Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                                                {
                            "name": "C Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        }
                    ],
                    "next_token": null
                }
            """

    Scenario: GET /v1/groups?sort=created and checking the response status 200
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=created"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "count": 3,
                    "groups": [
                        {
                            "name": "B Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                        {
                            "name": "A Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        },
                                                {
                            "name": "C Group",
                            "id": "",
                            "creation_date": "2010-01-01T00:00:00Z",
                            "last_modified_date": "2010-01-01T00:00:00Z",
                            "precedence": 1,
                            "role_arn": "",
                            "user_pool_id": ""
                        }
                    ],
                    "next_token": null
                }
            """

    Scenario: GET /v1/groups?sort=abc and checking the response status 400
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=abc"
        Then the HTTP status code should be "400"
    
    Scenario: GET /v1/groups?sort=name:xyz and checking the response status 400
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=name:xyz"
        Then the HTTP status code should be "400"

    Scenario: GET /v1/groups?sort=abc:asc and checking the response status 400
        Given group "B Group" exists in a list in the database
        And group "A Group" exists in a list in the database
        And group "C Group" exists in a list in the database
        And I am an admin user
        When I GET "/v1/groups?sort=abc:asc"
        Then the HTTP status code should be "400"

#   list for many groups found   given blocks of 60 for one cognito call
    Scenario: GET /v1/groups and checking the response status 200
        Given there are 100 groups in the database
            And I am an admin user
        When I GET "/v1/groups"
        Then the response code should be 200
            And the response should match the following json for listgroups
                """
                    {
                        "count": 100,
                        "groups": [
                            {
                                "name": "group name description 1",
                                "id": "group_name_1",
                                "precedence": 55
                            }
                        ],
                        "next_token": null
                    }
                """

#   successful return
    Scenario: GET /v1/groups and checking the response status 200
            Given group "test-group" exists in the database
            And I am an admin user
            When I GET "/v1/groups/test-group"
            Then I should receive the following JSON response with status "200":
                """
                    {
                        "id":"test-group",
                        "name":"A test group",
                        "precedence": 100,
                        "created": "2010-01-01T00:00:00Z"
                    }
            """
#   404 return
    Scenario: GET /v1/groups and checking the response status 404
        Given group "get-group-not-found" exists in the database
            And I am an admin user
        When I GET "/v1/groups/get-group-not-found"
        Then I should receive the following JSON response with status "404":
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

    Scenario: GET /v1/groups/{id} internal server error returns 500
        Given group "internal-error" exists in the database
            And I am an admin user
        When I GET "/v1/groups/internal-error"
        Then I should receive the following JSON response with status "500":
            """
                {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: GET /v1/groups/{id} without a JWT token and checking the response status 401
        When I GET "/v1/groups/test-group"
        Then the HTTP status code should be "401"

    Scenario: GET /v1/groups/{id} as a publisher user and checking the response status 403
        Given I am a publisher user
        When I GET "/v1/groups/test-group"
        Then the HTTP status code should be "403"

#   Delete deleteGroup scenarios
#   successful return
    Scenario: DELETE /v1/groups/{id} and checking the response status 204
        Given group "test-group" exists in the database
            And I am an admin user
        When I DELETE "/v1/groups/test-group"
        Then the HTTP status code should be "204"

    Scenario: DELETE /v1/groups/{id} without a JWT token and checking the response status 401
        When I DELETE "/v1/groups/test-group"
            """
            """
        Then the HTTP status code should be "401"

    Scenario: DELETE /v1/groups/{id} as a publisher user and checking the response status 403
        Given I am a publisher user
        When I DELETE "/v1/groups/test-group"
            """
            """
        Then the HTTP status code should be "403"

#   404 return
    Scenario: DELETE /v1/groups/{id} and checking the response status 404
        Given I am an admin user
        When I DELETE "/v1/groups/delete-group-not-found"
        Then the HTTP status code should be "404"

    Scenario: DELETE /v1/groups/{id} internal server error returns 500
        Given I am an admin user
        When I DELETE "/v1/groups/internal-error"
        Then the HTTP status code should be "500"

#   Put SetGroupUsers scenarios
    Scenario: PUT /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
            And a user with username "user_1" and email "email@ons.gov.uk" exists in the database
            And user "user_1" is a member of group "test-group"
            And a user with username "user_2" and email "email@ons.gov.uk" exists in the database
            And user "user_2" is a member of group "test-group"
            And there are 2 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I PUT "/v1/groups/test-group/members"
            """
                [
                    {
                        "user_id": "abcd1234"
                    }
                ]
            """
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

    Scenario: PUT /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
            And a user with username "user_1" and email "email@ons.gov.uk" exists in the database
            And user "user_1" is a member of group "test-group"
            And a user with username "user_2" and email "email@ons.gov.uk" exists in the database
            And user "user_2" is a member of group "test-group"
            And there are 2 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I PUT "/v1/groups/test-group/members"
            """
                [
                    {
                        "user_id": "abcd1234"
                    },
                    {
                        "user_id": "user_1"
                    }
                ]
            """
        Then I should receive the following JSON response with status "200":
            """
               {
                    "users": [
                        {
                            "id": "user_1",
                            "forename": "Bob",
                            "lastname": "Smith",
                            "email": "email@ons.gov.uk",
                            "groups": [],
                            "status": "CONFIRMED",
                            "active": true,
                            "status_notes": ""
                        },
                        {
                            "active":true,
                            "email":"email@ons.gov.uk",
                            "forename":"Bob",
                            "groups":[],
                            "id":"abcd1234",
                            "lastname":"Smith",
                            "status":"CONFIRMED",
                            "status_notes":""
                        }
                    ],
                    "count": 2,
                    "PaginationToken":""
                }
            """

    Scenario: PUT /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
            And a user with username "user_1" and email "email@ons.gov.uk" exists in the database
            And user "user_1" is a member of group "test-group"
            And a user with username "user_2" and email "email@ons.gov.uk" exists in the database
            And user "user_2" is a member of group "test-group"
            And there are 2 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I PUT "/v1/groups/test-group/members"
            """
                []
            """
        Then I should receive the following JSON response with status "200":
            """
               {
                    "users": [],
                    "count": 0,
                    "PaginationToken":""
                }
            """

    Scenario: PUT /v1/groups/{id}/members and non-admin user
        Given group "test-group" exists in the database
            And a user with username "user_1" and email "email@ons.gov.uk" exists in the database
            And user "user_1" is a member of group "test-group"
            And a user with username "user_2" and email "email@ons.gov.uk" exists in the database
            And user "user_2" is a member of group "test-group"
            And there are 2 users in group "test-group"
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I PUT "/v1/groups/test-group/members"
            """
                []
            """
        Then the HTTP status code should be "401"

    Scenario: PUT /v1/groups/{id}/members and checking the response status 200
        Given a user with username "user_1" and email "email@ons.gov.uk" exists in the database
            And a user with username "user_2" and email "email@ons.gov.uk" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And I am an admin user
        When I PUT "/v1/groups/test-group/members"
            """
                [
                    {
                        "user_id": "abcd1234"
                    },
                    {
                        "user_id": "user_1"
                    }
                ]
            """
        Then the HTTP status code should be "404"

@groups-report
    Scenario: GET /v1/groups-report checking the response status 200 got an empty report no groups
        Given I am an admin user
        When I GET "/v1/groups-report"
        Then the response header "Content-Type" should contain "application/json"
        And I should receive the following JSON response with status "200":
            """
            []
            """
@groups-report
    Scenario: GET /v1/groups-report checking the response status 200 an empty report with one groups but no users
        Given I am an admin user
        And group "test-group" and description "test group description" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I GET "/v1/groups-report"
        Then the response code should be 200
        And the response header "Content-Type" should contain "application/json"


    @groups-report
    Scenario: GET /v1/groups-report as a publisher user and checking the response status 403
        Given I am a publisher user
        When I GET "/v1/groups-report"
        Then the HTTP status code should be "403"

@groups-report
    Scenario: GET /v1/groups-report without a JWT token and checking the response status 401
        When I GET "/v1/groups-report"
        Then the HTTP status code should be "401"

@groups-report
    Scenario: GET /v1/groups-report checking the response status 200 one group with member
        Given I am an admin user
            And group "test-group" and description "test group description" exists in the database
            And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
            And user "abcd1234" is a member of group "test-group"
        When I GET "/v1/groups-report"
        Then the response header "Content-Type" should contain "application/json"
            And I should receive the following JSON response with status "200":
                    """
                        [{"group":"test group description",
                        "user":"email@ons.gov.uk"}]
                    """

@groups-report
    Scenario: GET /v1/groups-report checking the response status 200 many groups with members
        Given I am an admin user
        And group "test-group_1" and description "test group_1 description" exists in the database
        And group "test-group_2" and description "test group_2 description" exists in the database
        And group "test-group_3" and description "test group_3 description" exists in the database
        And group "test-group_4" and description "test group_4 description" exists in the database
        And group "test-group_5" and description "test group_5 description" exists in the database

        And a user with username "abcd1234" and email "email1@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group_1"
        And user "abcd1234" is a member of group "test-group_2"
        And user "abcd1234" is a member of group "test-group_4"
        And user "abcd1234" is a member of group "test-group_5"

        And a user with username "abcd1235" and email "email2@ons.gov.uk" exists in the database
        And user "abcd1235" is a member of group "test-group_1"
        And user "abcd1235" is a member of group "test-group_2"

        And a user with username "abcd1236" and email "email3@ons.gov.uk" exists in the database
        And user "abcd1236" is a member of group "test-group_4"
        And user "abcd1236" is a member of group "test-group_5"

    When I GET "/v1/groups-report"
        Then I should receive the following JSON response with status "200":
                """
                    [
                        {"group": "test group_1 description",
                            "user": "email1@ons.gov.uk"},
                        {"group": "test group_1 description",
                            "user": "email2@ons.gov.uk"},
                        {"group": "test group_2 description",
                            "user": "email1@ons.gov.uk"},
                        {"group": "test group_2 description",
                            "user": "email2@ons.gov.uk"},
                        {"group": "test group_4 description",
                            "user": "email1@ons.gov.uk"},
                        {"group": "test group_4 description",
                            "user": "email3@ons.gov.uk"},
                        {"group": "test group_5 description",
                            "user": "email1@ons.gov.uk"},
                        {"group": "test group_5 description",
                            "user":  "email3@ons.gov.uk" }]
                """

@groups-report
    Scenario: GET /v1/groups-report checking the response status 200 many groups with members
        Given I am an admin user

        And group "test-group_1" and description "test group_1 description" exists in the database
        And group "test-group_2" and description "test group_2 description" exists in the database
        And group "test-group_3" and description "test group_3 description" exists in the database
        And group "test-group_4" and description "test group_4 description" exists in the database
        And group "test-group_5" and description "test group_5 description" exists in the database

        And a user with username "abcd1234" and email "email1@ons.gov.uk" exists in the database
        And user "abcd1234" is a member of group "test-group_1"
        And user "abcd1234" is a member of group "test-group_2"
        And user "abcd1234" is a member of group "test-group_4"
        And user "abcd1234" is a member of group "test-group_5"

        And a user with username "abcd1235" and email "email2@ons.gov.uk" exists in the database
        And user "abcd1235" is a member of group "test-group_1"
        And user "abcd1235" is a member of group "test-group_2"

        And a user with username "abcd1236" and email "email3@ons.gov.uk" exists in the database
        And user "abcd1236" is a member of group "test-group_4"
        And user "abcd1236" is a member of group "test-group_5"
        And request header Accept is "text/csv"

        When I GET "/v1/groups-report"
        Then the HTTP status code should be "200"
        And the response should match the following csv:
            """
            Group,User
            test group_1 description,email1@ons.gov.uk
            test group_1 description,email2@ons.gov.uk
            test group_2 description,email1@ons.gov.uk
            test group_2 description,email2@ons.gov.uk
            test group_4 description,email1@ons.gov.uk
            test group_4 description,email3@ons.gov.uk
            test group_5 description,email1@ons.gov.uk
            test group_5 description,email3@ons.gov.uk
            """
        And the response header "Content-Type" should contain "text/csv"






