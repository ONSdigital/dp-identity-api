@Groups @GroupsAddUser
Feature: Groups - Add User
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
            "email": "email@ons.gov.uk",
            "forename": "Bob",
            "groups": [],
            "id": "abcd1234",
            "lastname": "Smith",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ],
        "count": 1,
        "PaginationToken": ""
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
        "PaginationToken": ""
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
            "active": true,
            "email": "email@ons.gov.uk",
            "forename": "Bob",
            "groups": [],
            "id": "abcd1234",
            "lastname": "Smith",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ],
        "count": 2,
        "PaginationToken": ""
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
        "PaginationToken": ""
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
