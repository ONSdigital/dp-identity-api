@Users @UsersCreate
Feature: Users - Create
  Scenario: POST /v1/users and checking the response status 201
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "smileons",
        "lastname": "bobbings",
        "email": "emailx@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "201":
      """
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "forename": "smileons",
        "lastname": "bobbings",
        "email": "emailx@ons.gov.uk",
        "groups": [],
        "status": "FORCE_CHANGE_PASSWORD",
        "active": true,
        "status_notes": ""
      }
      """

  Scenario: POST /v1/users without a JWT token and checking the response status 401
    Given I POST "/v1/users"
      """
      """
    Then the HTTP status code should be "401"

  Scenario: POST /v1/users as a publisher user and checking the response status 403
    Given I am a publisher user
    When I POST "/v1/users"
      """
      """
    Then the HTTP status code should be "403"

  Scenario: POST /v1/users missing email and checking the response status 400
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "smileons",
        "lastname": "bobbings",
        "email": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          }
        ]
      }
      """

  Scenario: POST /v1/users missing forename and checking the response status 400
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "",
        "lastname": "bobbings",
        "email": "emailx@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidForename",
            "description": "the submitted user's forename could not be validated"
          }
        ]
      }
      """

  Scenario: POST /v1/users missing lastname and checking the response status 400
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "smileons",
        "lastname": "",
        "email": "emailx@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidSurname",
            "description": "the submitted user's lastname could not be validated"
          }
        ]
      }
      """

  Scenario: POST /v1/users and checking the response status 400
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "",
        "lastname": "",
        "email": ""
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidForename",
            "description": "the submitted user's forename could not be validated"
          },
          {
            "code": "InvalidSurname",
            "description": "the submitted user's lastname could not be validated"
          },
          {
            "code": "InvalidEmail",
            "description": "the submitted email could not be validated"
          }
        ]
      }
      """

  Scenario: POST /v1/users and checking the response status 500
    Given I am an admin user
    When I POST "/v1/users"
      """

      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: POST /v1/users unexpected server error and checking the response status 500
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "bob",
        "lastname": "bobbings",
        "email": "emailx@ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: POST /v1/users duplicate email found and checking the response status 400
    Given I am an admin user
    When I POST "/v1/users"
      """
      {
        "forename": "bob",
        "lastname": "bobbings",
        "email": "email@ext.ons.gov.uk"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidEmail",
            "description": "account using email address found"
          }
        ]
      }
      """
