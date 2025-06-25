@Users @UsersGroups
Feature: Users - Groups
  Scenario: GET /v1/users/{id}/groups and checking the response status 200
    Given a user with username "listgrouptestuser" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    And 1 groups exist in the database that username "listgrouptestuser" is a member
    When I GET "/v1/users/listgrouptestuser/groups"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 1,
        "groups": [
          {
            "creation_date": null,
            "name": "group name description 0",
            "id": "group_name_0",
            "last_modified_date": null,
            "precedence": 13,
            "role_arn": null,
            "user_pool_id": null
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/users/{id}/groups  for 0 groups and checking the response status 200
    Given a user with username "listgrouptestuser2" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    And 0 groups exist in the database that username "listgrouptestuser2" is a member
    When I GET "/v1/users/listgrouptestuser2/groups"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 0,
        "groups": null,
        "next_token": null
      }
      """

  Scenario: GET /v1/users/{id}/groups  user not found returns 500
    Given a user with username "get-user-not-found" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    And 0 groups exist in the database that username "get-user-not-found" is a member
    When I GET "/v1/users/get-user-not-found/groups"
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: GET /v1/users/{id}/groups without a JWT token and checking the response status 401
    When I GET "/v1/users/listgrouptestuser/groups"
    Then the HTTP status code should be "401"

  Scenario: GET /v1/users/{id}/groups as a publisher user and checking the response status 403
    Given I am a publisher user
    When I GET "/v1/users/listgrouptestuser/groups"
    Then the HTTP status code should be "403"
