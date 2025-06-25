@Groups @GroupsGet
Feature: Groups - Get Group
  Scenario: GET /v1/groups and checking the response status 200
    Given group "test-group" exists in the database
    And I am an admin user
    When I GET "/v1/groups/test-group"
    Then I should receive the following JSON response with status "200":
      """
      {
        "id": "test-group",
        "name": "A test group",
        "precedence": 100,
        "created": "2010-01-01T00:00:00Z"
      }
      """

  Scenario: GET /v1/groups for unknown group and checking the response status 404
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
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: GET /v1/groups/{id} without a JWT token and checking the response status 401
    When I GET "/v1/groups/test-group"
    Then the HTTP status code should be "401"

  Scenario: GET /v1/groups/{id} as a publisher user and checking the response status 403
    Given I am a publisher user
    When I GET "/v1/groups/test-group"
    Then the HTTP status code should be "403"

