@Groups @GroupsUsers
Feature: Groups - Get Users
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
        "PaginationToken": ""
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
        "PaginationToken": ""
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
        "PaginationToken": ""
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
        "PaginationToken": ""
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
        "PaginationToken": ""
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
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """
