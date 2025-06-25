@Users @UsersUpdate
Feature: Users - Update
  Scenario: PUT /v1/users/{id} to update users names and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Changed",
        "lastname": "Names",
        "active": true,
        "status_notes": ""
      }
      """
    Then I should receive the following JSON response with status "200":
      """
      {
        "id": "abcd1234",
        "forename": "Changed",
        "lastname": "Names",
        "email": "email@ons.gov.uk",
        "groups": [],
        "status": "CONFIRMED",
        "active": true,
        "status_notes": ""
      }
      """

  Scenario: PUT /v1/users/{id} set user disabled and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": false,
        "status_notes": "user disabled"
      }
      """
    Then I should receive the following JSON response with status "200":
      """
      {
        "id": "abcd1234",
        "forename": "Bob",
        "lastname": "Smith",
        "email": "email@ons.gov.uk",
        "groups": [],
        "status": "CONFIRMED",
        "active": false,
        "status_notes": "user disabled"
      }
      """

  Scenario: PUT /v1/users/{id} set user enabled and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And user "abcd1234" active is "false"
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": true,
        "status_notes": "user reactivated"
      }
      """
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
        "status_notes": "user reactivated"
      }
      """

  Scenario: PUT /v1/users/{id} set user disabled and change names and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Changed",
        "lastname": "Names",
        "active": false,
        "status_notes": "user suspended"
      }
      """
    Then I should receive the following JSON response with status "200":
      """
      {
        "id": "abcd1234",
        "forename": "Changed",
        "lastname": "Names",
        "email": "email@ons.gov.uk",
        "groups": [],
        "status": "CONFIRMED",
        "active": false,
        "status_notes": "user suspended"
      }
      """

  Scenario: PUT /v1/users/{id} set user enabled and change names and checking the response status 200
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And user "abcd1234" active is "false"
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Changed",
        "lastname": "Names",
        "active": true,
        "status_notes": "user reactivated"
      }
      """
    Then I should receive the following JSON response with status "200":
      """
      {
        "id": "abcd1234",
        "forename": "Changed",
        "lastname": "Names",
        "email": "email@ons.gov.uk",
        "groups": [],
        "status": "CONFIRMED",
        "active": true,
        "status_notes": "user reactivated"
      }
      """

  Scenario: PUT /v1/users/{id} without a JWT token and checking the response status 401
    When I PUT "/v1/users/abcd1234"
      """
      """
    Then the HTTP status code should be "401"

  Scenario: PUT /v1/users/{id} as a publisher user and checking the response status 403
    Given I am a publisher user
    When I PUT "/v1/users/abcd1234"
      """
      """
    Then the HTTP status code should be "403"

  Scenario: PUT /v1/users/{id} missing forename and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "",
        "lastname": "Smith",
        "active": true,
        "status_notes": ""
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

  Scenario: PUT /v1/users/{id} missing lastname and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "",
        "active": true,
        "status_notes": ""
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

  Scenario: PUT /v1/users/{id} invalid notes and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Stan",
        "lastname": "Smith",
        "active": true,
        "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidStatusNotes",
            "description": "the status notes are too long"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/{id} missing forename and lastname and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "",
        "lastname": "",
        "active": true,
        "status_notes": ""
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
          }
        ]
      }
      """

  Scenario: PUT /v1/users/{id} missing forename and invalid notes and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "",
        "lastname": "Smith",
        "active": true,
        "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
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
            "code": "InvalidStatusNotes",
            "description": "the status notes are too long"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/{id} missing lastname and invalid notes and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Stan",
        "lastname": "",
        "active": true,
        "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
      }
      """
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {
            "code": "InvalidSurname",
            "description": "the submitted user's lastname could not be validated"
          },
          {
            "code": "InvalidStatusNotes",
            "description": "the status notes are too long"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/{id} missing forename, lastname and invalid notes and checking the response status 400
    Given a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "",
        "lastname": "",
        "active": true,
        "status_notes": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras eu turpis libero. Sed convallis pharetra mollis. Mauris ex nisi, finibus in mi quis, tincidunt pulvinar risus. Ut iaculis lobortis nisl. Suspendisse venenatis ante congue erat posuere, eget mattis massa facilisis. Vivamus bibendum pharetra suscipit. Integer laoreet molestie velit, vitae euismod ligula dictum eu. Phasellus a fermentum metus, nec dignissim ex. Sed dolor lectus, sollicitudin sit amet imperdiet eget, fringilla nec felis. Morbi commodo diam massa, sed interdum tellus sit"
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
            "code": "InvalidStatusNotes",
            "description": "the status notes are too long"
          }
        ]
      }
      """

  Scenario: PUT /v1/users/{id} user not found and checking the response status 404
    Given I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": true,
        "status_notes": ""
      }
      """
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

  Scenario: PUT /v1/users/{id} unexpected server error disabling user and checking the response status 500
    Given a user with username "abcd1234" and email "disable.internalerror@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": false,
        "status_notes": "user suspended"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: PUT /v1/users/{id} unexpected server error enabling user and checking the response status 500
    Given a user with username "abcd1234" and email "enable.internalerror@ons.gov.uk" exists in the database
    And user "abcd1234" active is "false"
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": true,
        "status_notes": "user reactivated"
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: PUT /v1/users/{id} unexpected server error updating user and checking the response status 500
    Given a user with username "abcd1234" and email "update.internalerror@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": true,
        "status_notes": ""
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  Scenario: PUT /v1/users/{id} unexpected server error loading updated user and checking the response status 500
    Given a user with username "abcd1234" and email "internal.error@ons.gov.uk" exists in the database
    And I am an admin user
    When I PUT "/v1/users/abcd1234"
      """
      {
        "forename": "Bob",
        "lastname": "Smith",
        "active": true,
        "status_notes": ""
      }
      """
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """
