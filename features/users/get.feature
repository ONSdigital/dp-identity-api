@Users @UsersGet
Feature: Users - Get a User
  Scenario: GET /v1/users/{id} and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I GET "/v1/users/abcd1234"
    Then I should receive the following JSON response with status "200":
      """
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
      """

  Scenario: GET /v1/users/{id} without a JWT token and checking the response status 401
    When I GET "/v1/users/abcd1234"
    Then the HTTP status code should be "401"

  Scenario: GET /v1/users/{id} as a publisher user and checking the response status 403
    Given I am a publisher user
    When I GET "/v1/users/abcd1234"
    Then the HTTP status code should be "403"

  Scenario: GET /v1/users/{id} user not found and checking the response status 404
    Given I am an admin user
    When I GET "/v1/users/abcd1234"
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

  Scenario: GET /v1/users/{id} unexpected server error and checking the response status 500
    Given a user with username "abcd1234" and email "internal.error@ons.gov.uk" exists in the database
    And I am an admin user
    When I GET "/v1/users/abcd1234"
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """
