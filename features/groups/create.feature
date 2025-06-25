@Groups @GroupsCreate
Feature: Groups - Create
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
            "code": "InvalidGroupName",
            "description": "the group name was not found"
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
            "code": "InvalidGroupPrecedence",
            "description": "the group precedence was not found"
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
            "code": "InvalidGroupName",
            "description": "a group name cannot start with 'role-' or 'ROLE-'"
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
            "code": "InvalidGroupName",
            "description": "a group name cannot start with 'role-' or 'ROLE-'"
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
            "code": "InvalidGroupPrecedence",
            "description": "the group precedence needs to be a minumum of 10 and maximum of 100"
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
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """
