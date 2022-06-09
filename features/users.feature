Feature: Users

#   Create User
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
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

#   List Users
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
                "PaginationToken":""
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
                {"code":"InternalServerError", "description":"Internal Server Error"}
            """

#   Get User
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

#   Update User
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
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
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

#   Change password - auth challenge
    Scenario: PUT /v1/users/self/password and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then the HTTP status code should be "202"
        And the response header "Authorization" should be "Bearer accessToken"
        And the response header "ID" should be "idToken"
        And the response header "Refresh" should be "refreshToken"

    Scenario: PUT /v1/users/self/password missing type and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "UnknownRequestType",
                        "description": "unknown password change type received"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password new password required type with verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "Password2",
                "session": "auth-challenge-session"
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

    Scenario: PUT /v1/users/self/password missing password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "Password2",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing password and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "",
                "password": "",
                "session": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidEmail",
                        "description": "the submitted email could not be validated"
                    },
                    {
                        "code": "InvalidChallengeSession",
                        "description": "no valid auth challenge session was provided"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito internal error
        Given an internal server error is returned from Cognito
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "internalerrorException",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: PUT /v1/users/self/password Cognito invalid password
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "invalidpassword",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "password does not meet requirements"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito user not found
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "NewPasswordRequired",
                "email": "email@ons.gov.uk",
                "password": "Password",
                "session": "auth-challenge-session"
            }
            """
        Then the HTTP status code should be "202"

#   Change password - forgotten password
    Scenario: PUT /v1/users/self/password forgotten password type and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then the HTTP status code should be "202"

    Scenario: PUT /v1/users/self/password missing type and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "UnknownRequestType",
                        "description": "unknown password change type received"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password forgotten password type with challenge session and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "session": "auth-challenge-session"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "Password2",
                "verification_token": "verification-token"
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

    Scenario: PUT /v1/users/self/password missing password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidUserId",
                        "description": "the user id was missing"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "Password2",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidUserId",
                        "description": "the user id was missing"
                    },
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password missing email and password and verification token and checking the response status 400
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "",
                "password": "",
                "verification_token": ""
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "the submitted password could not be validated"
                    },
                    {
                        "code": "InvalidUserId",
                        "description": "the user id was missing"
                    },
                    {
                        "code": "InvalidToken",
                        "description": "the submitted token could not be validated"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito internal error for ForgottenPassword
        Given an internal server error is returned from Cognito
        And I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "internalerrorException",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: PUT /v1/users/self/password Cognito invalid password for forgottenPassword
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "invalidpassword",
                "verification_token": "verification-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidPassword",
                        "description": "password does not meet requirements"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito invalid token
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "invalid-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "InvalidCode",
                        "description": "verification token does not meet requirements"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito expired token
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password2",
                "verification_token": "expired-token"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "ExpiredCode",
                        "description": "verification token has expired"
                    }
                ]
            }
            """

    Scenario: PUT /v1/users/self/password Cognito user not found
        Given I am an admin user
        When I PUT "/v1/users/self/password"
            """
            {
                "type": "ForgottenPassword",
                "email": "email@ons.gov.uk",
                "password": "Password",
                "verification_token": "verification-token"
            }
            """
        Then the HTTP status code should be "202"

#   Request password reset
    Scenario: POST /v1/password-reset for an existing user and checking the response status 202
        Given a user with email "email@ons.gov.uk" and password "Passw0rd!" exists in the database
        When I POST "/v1/password-reset"
            """
            {
                "email": "email@ons.gov.uk"
            }
            """
        Then the HTTP status code should be "202"

    Scenario: POST /v1/password-reset missing email and checking the response status 400
        When I POST "/v1/password-reset"
            """
            {
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

    Scenario: POST /v1/password-reset non ONS email address and checking the response status 202
        When I POST "/v1/password-reset"
            """
            {
                "email": "email@gmail.com"
            }
            """
        Then the HTTP status code should be "202"

    Scenario: POST /v1/password-reset Cognito internal error
        Given an internal server error is returned from Cognito
        When I POST "/v1/password-reset"
            """
            {
                "email": "internal.error@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "500":
            """
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: POST /v1/password-reset Cognito too many requests error
        Given an internal server error is returned from Cognito
        When I POST "/v1/password-reset"
            """
            {
                "email": "too.many@ons.gov.uk"
            }
            """
        Then I should receive the following JSON response with status "400":
            """
            {
                "errors": [
                    {
                        "code": "TooManyRequests",
                        "description": "Slow down"
                    }
                ]
            }
            """

    Scenario: POST /v1/password-reset Cognito user not found
        When I POST "/v1/password-reset"
            """
            {
                "email": "email@ons.gov.uk"
            }
            """
        Then the HTTP status code should be "202"

#   List get users for user
    Scenario: GET /v1/users/{id}/groups and checking the response status 200
        Given a user with username "listgrouptestuser" and email "email@ons.gov.uk" exists in the database
        And I am an admin user
        And there 1 groups exists in the database that username "listgrouptestuser" is a member
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
        And there 0 groups exists in the database that username "listgrouptestuser2" is a member
        When I GET "/v1/users/listgrouptestuser2/groups"
        Then I should receive the following JSON response with status "200":
            """
            {
                "count":0,
                "groups":null,
                "next_token":null
            }
            """

    Scenario: GET /v1/users/{id}/groups  user not found returns 500
        Given a user with username "get-user-not-found" and email "email@ons.gov.uk" exists in the database
        And I am an admin user
        And there 0 groups exists in the database that username "get-user-not-found" is a member
        When I GET "/v1/users/get-user-not-found/groups"
        Then I should receive the following JSON response with status "500":
            """
            {"code":"InternalServerError", "description":"Internal Server Error"}
            """

    Scenario: GET /v1/users/{id}/groups without a JWT token and checking the response status 401
        When I GET "/v1/users/listgrouptestuser/groups"
        Then the HTTP status code should be "401"

    Scenario: GET /v1/users/{id}/groups as a publisher user and checking the response status 403
        Given I am a publisher user
        When I GET "/v1/users/listgrouptestuser/groups"
        Then the HTTP status code should be "403"
