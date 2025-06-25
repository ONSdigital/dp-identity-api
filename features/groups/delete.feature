@Groups @GroupsDelete
Feature: Groups - Delete
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

  Scenario: DELETE /v1/groups/{id} for unknown group and checking the response status 404
    Given I am an admin user
    When I DELETE "/v1/groups/delete-group-not-found"
    Then the HTTP status code should be "404"

  Scenario: DELETE /v1/groups/{id} internal server error returns 500
    Given I am an admin user
    When I DELETE "/v1/groups/internal-error"
    Then the HTTP status code should be "500"
