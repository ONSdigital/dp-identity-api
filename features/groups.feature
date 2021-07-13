Feature: Groups

    # Add user to group scenarios

    Scenario: POST /v1/groups/{id}/members and checking the response status 200
        Given group "test-group" exists in the database
        And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
        When I POST "/v1/groups/test-group/members"
            """
                {
                    "user_id": "abcd1234",
                }
            """
        Then I should receive the following JSON response with status "200":
            """
                {
                    "id": "123e4567-e89b-12d3-a456-426614174000",
                    "forename": "smileons",
                    "lastname": "bobbings",
                    "email": "emailx@ons.gov.uk",
                    "groups": [],
                    "status": "FORCE_CHANGE_PASSWORD"
                }
            """