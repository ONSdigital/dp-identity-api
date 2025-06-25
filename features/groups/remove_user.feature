@Groups @GroupsRemoveUser
Feature: Groups - Remove User
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
        "PaginationToken": ""
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
        "PaginationToken": ""
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
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, internal server error returns 500
    Given group "get-group-internal-error" exists in the database
    And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I DELETE "/v1/groups/get-group-internal-error/members/abcd1234"
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: DELETE /v1/groups/{id}/members/{user_id} get group, group not found returns 404
    Given group "get-group-not-found" exists in the database
    And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I DELETE "/v1/groups/get-group-not-found/members/abcd1234"
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
