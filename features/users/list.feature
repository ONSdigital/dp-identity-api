@Users @UsersList
Feature: Users - List
  Scenario: GET /v1/users and checking the response status 200
    Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
    And a user with non-verified email "new_email@ons.gov.uk" and password "TeMpPassw0rd!"
    And I am an admin user
    When I GET "/v1/users"
    Then I should receive the following JSON response with status "200":
      """
      {
        "users": [
          {
            "id": "aaaabbbbcccc",
            "forename": "Bob",
            "lastname": "Smith",
            "email": "email@ons.gov.uk",
            "groups": [],
            "status": "CONFIRMED",
            "active": true,
            "status_notes": ""
          },
          {
            "id": "aaaabbbbcccc",
            "forename": "Bob",
            "lastname": "Smith",
            "email": "new_email@ons.gov.uk",
            "groups": [],
            "status": "FORCE_CHANGE_PASSWORD",
            "active": true,
            "status_notes": ""
          }
        ],
        "count": 2,
        "PaginationToken": ""
      }
      """

  Scenario: GET /v1/users without a JWT token and checking the response status 401
    When I GET "/v1/users"
    Then the HTTP status code should be "401"

  Scenario: GET /v1/users as a publisher user and checking the response status 403
    Given I am a publisher user
    When I GET "/v1/users"
    Then the HTTP status code should be "403"

  Scenario: GET /v1/users with more than 60 active users and  30 inactive users checking the response status 200 with the correct number of users
    Given there are "70" active users and "30" inactive users in the database
    And I am an admin user
    When I GET "/v1/users"
    Then the HTTP status code should be "200"
    And the list response should contain "100" entries

  Scenario: GET /v1/users?active=true with more than 60 active users and  30 inactive users checking the response status 200 with the correct number of users
    Given there are "70" active users and "30" inactive users in the database
    And I am an admin user
    When I GET "/v1/users?active=true"
    Then the HTTP status code should be "200"
    And the list response should contain "70" entries

  Scenario: GET /v1/users?active=false with more than 60 active users and  30 inactive users checking the response status 200 with the correct number of users
    Given there are "70" active users and "30" inactive users in the database
    And I am an admin user
    When I GET "/v1/users?active=false"
    Then the HTTP status code should be "200"
    And the list response should contain "30" entries

  Scenario: GET /v1/users?active=anything with more than 60 active users and checking the response status 400 with the correct number of users
    Given there are "70" active users and "30" inactive users in the database
    And I am an admin user
    When I GET "/v1/users?active=anything"
    Then the HTTP status code should be "400"

  Scenario: GET /v1/user?active=false with more than 60 active users and checking the response status 404 with the correct number of users
    Given there are "70" active users and "30" inactive users in the database
    And I am an admin user
    When I GET "/v1/user?active=false"
    Then the HTTP status code should be "404"

  Scenario: GET /v1/users unexpected server error and checking the response status 500
    Given a user with email "internal.error@ons.gov.uk" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users"
    Then I should receive the following JSON response with status "500":
      """
      {
        "code": "InternalServerError",
        "description": "Internal Server Error"
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort default
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 3,
        "users": [
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort by id
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users?sort=id"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 3,
        "users": [
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort by id:asc
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users?sort=id:asc"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 3,
        "users": [
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort by id:desc
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users?sort=id:desc"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 3,
        "users": [
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """
  @get-users-list
  Scenario: GET /v1/users and checking the response sort by dog
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And I am an admin user
    When I GET "/v1/users?sort=dog"
    Then I should receive the following JSON response with status "400":
      """
      {
        "errors": [
          {}
        ]
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort by forename,lastname
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And a user with forename "Adam", lastname "Williams", email "email10@ons.gov.uk", id "id_4" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Adams", email "email11@ons.gov.uk", id "id_5" and password "Passw0rd!" exists in the database

    And I am an admin user
    When I GET "/v1/users?sort=forename,lastname"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 5,
        "users": [
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email10@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_4",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email11@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_5",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """

  @get-users-list
  Scenario: GET /v1/users and checking the response status 200 with sort by forename,lastname:desc
    Given a user with forename "Adam", lastname "Adams", email "email5@ons.gov.uk", id "id_2" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Williams", email "email9@ons.gov.uk", id "id_1" and password "Passw0rd!" exists in the database
    And a user with forename "Mary", lastname "Martin", email "email7@ons.gov.uk", id "id_3" and password "Passw0rd!" exists in the database
    And a user with forename "Adam", lastname "Williams", email "email10@ons.gov.uk", id "id_4" and password "Passw0rd!" exists in the database
    And a user with forename "William", lastname "Adams", email "email11@ons.gov.uk", id "id_5" and password "Passw0rd!" exists in the database

    And I am an admin user
    When I GET "/v1/users?sort=forename,lastname:desc"
    Then I should receive the following JSON response with status "200":
      """
      {
        "PaginationToken": "",
        "count": 5,
        "users": [
          {
            "active": true,
            "email": "email10@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_4",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email5@ons.gov.uk",
            "forename": "Adam",
            "groups": [],
            "id": "id_2",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email7@ons.gov.uk",
            "forename": "Mary",
            "groups": [],
            "id": "id_3",
            "lastname": "Martin",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email9@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_1",
            "lastname": "Williams",
            "status": "CONFIRMED",
            "status_notes": ""
          },
          {
            "active": true,
            "email": "email11@ons.gov.uk",
            "forename": "William",
            "groups": [],
            "id": "id_5",
            "lastname": "Adams",
            "status": "CONFIRMED",
            "status_notes": ""
          }
        ]
      }
      """
