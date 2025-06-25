@Groups @GroupsUpdate
Feature: Groups - Update
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
            "code": "InvalidGroupName",
            "description": "the group name was not found"
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
            "code": "InvalidGroupName",
            "description": "a group name cannot start with 'role-' or 'ROLE-'"
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
            "code": "InvalidGroupName",
            "description": "a group name cannot start with 'role-' or 'ROLE-'"
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
            "code": "InvalidGroupPrecedence",
            "description": "the group precedence needs to be a minumum of 10 and maximum of 100"
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
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
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
            "code": "NotFound",
            "description": "Resource not found"
          }
        ]
      }
      """
