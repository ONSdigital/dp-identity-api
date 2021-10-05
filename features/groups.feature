Feature: Groups

#   Create new group scenarios
    Scenario: POST /v1/groups to create group, group created returns 201
    When I POST "/v1/groups"
        """
            {
                "description": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                "precedence": 49
            }
        """
    Then I should receive the following JSON response with status "201":
        """
            {
                "description": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                "precedence": 49,
                "GroupName": "thisisatestdescriptionforanewgroup"
            }
        """

    Scenario: POST /v1/groups to create group with no description in request, group created returns 400
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
                        "code":"InvalidGroupDescription",
                        "description":"the group description was not found"
                    }
                ]
            }
        """

    Scenario: POST /v1/groups to create group with no precedence in request, group created returns 400
    When I POST "/v1/groups"
        """
            {
                "description": "Thi$s is a te||st des$%£@^c ription for  a n ew group  $"
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
    When I POST "/v1/groups"
        """
            {
                "description": "role_Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                "precedence": 49
            }
        """
    Then I should receive the following JSON response with status "400":
        """
            {
                "errors": [
                    {
                        "code":"InvalidGroupDescription",
                        "description":"a group description cannot start with 'role_' or 'ROLE_'"
                    }
                ]
            }
        """

    Scenario: POST /v1/groups to create group with reserved pattern in description [upper case], group created returns 400
    When I POST "/v1/groups"
        """
            {
                "description": "ROLE_Thi$s is a te||st des$%£@^c ription for  a n ew group  $",
                "precedence": 49
            }
        """
    Then I should receive the following JSON response with status "400":
        """
            {
                "errors": [
                    {
                        "code":"InvalidGroupDescription",
                        "description":"a group description cannot start with 'role_' or 'ROLE_'"
                    }
                ]
            }
        """

    Scenario: POST /v1/groups to create group group precedence doesn't meet minimum of `3`, returns 400
    When I POST "/v1/groups"
        """
            {
                "description": "This is a test description",
                "precedence": 1
            }
        """
    Then I should receive the following JSON response with status "400":
        """
            {
                "errors": [
                    {
                        "code":"InvalidGroupPrecedence",
                        "description":"the group precedence needs to be a minumum of 3"
                    }
                ]
            }
        """

    Scenario: POST /v1/groups to create group an unexpected 500 error is returned from Cognito
    When I POST "/v1/groups"
        """
            {
                "description": "Internal Server Error",
                "precedence": 5
            }
        """
    Then I should receive the following JSON response with status "500":
        """
            {
                "errors": [
                    {
                        "code":"InternalServerError",
                        "description":"Something went wrong"
                    }
                ]
            }
        """
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

#   Get listgroups scenarios     
#   list for no groups found   
    Scenario: GET /v1/groups and checking the response status 200
        Given there "0" groups exists in the database
        When I GET "/v1/groups"
        Then the response code should be 200
        And the response should match the following json for listgroups
            """
                {
                    "groups":null,
                    "count":0,
                    "next_token":null
                }
            """  
#   list for one groups found  
    Scenario: GET /v1/groups and checking the response status 200
        Given there "2" groups exists in the database
        When I GET "/v1/groups"
        Then the response code should be 200
        And the response should match the following json for listgroups
            """
                {
                "count": 2,
                "groups": [
                    {
                    "description": "group name description 1",
                    "group_name": "group_name_1",
                    "precedence": 55
                    }
                ],
                "next_token": null
                }
            """  
#   list for many groups found   given blocks of 60 for one cognito call
    Scenario: GET /v1/groups and checking the response status 200
        Given there "100" groups exists in the database
        When I GET "/v1/groups"
        Then the response code should be 200
        And the response should match the following json for listgroups
            """
                {
                "count": 100,
                "groups": [
                    {
                    "description": "group name description 1",
                    "group_name": "group_name_1",
                    "precedence": 55
                    }
                ],
                "next_token": null
                }
            """  

#   Get getGroup scenarios     
#   successful return   
    Scenario: GET /v1/groups and checking the response status 200
        Given group "test-group" exists in the database
        When I GET "/v1/groups/test-group"
        Then I should receive the following JSON response with status "200":
            """
                {
                    "name":"test-group",
                    "description":"A test group",
                    "precedence": 100,
                    "created": "2010-01-01T00:00:00Z",
                    "members": null
                }
            """  
#   404 return   
    Scenario: GET /v1/groups and checking the response status 404
        Given group "get-group-not-found" exists in the database
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
        When I GET "/v1/groups/internal-error"
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