@Groups @GroupsReport
Feature: Groups - Report
  Scenario: GET /v1/groups-report checking the response status 200 got an empty report no groups
    Given I am an admin user
    When I GET "/v1/groups-report"
    Then the response header "Content-Type" should contain "application/json"
    And I should receive the following JSON response with status "200":
      """
      []
      """

  Scenario: GET /v1/groups-report checking the response status 200 an empty report with one groups but no users
    Given I am an admin user
    And group "test-group" and description "test group description" exists in the database
    And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    When I GET "/v1/groups-report"
    Then the response code should be 200
    And the response header "Content-Type" should contain "application/json"

  Scenario: GET /v1/groups-report as a publisher user and checking the response status 403
    Given I am a publisher user
    When I GET "/v1/groups-report"
    Then the HTTP status code should be "403"

  Scenario: GET /v1/groups-report without a JWT token and checking the response status 401
    When I GET "/v1/groups-report"
    Then the HTTP status code should be "401"

  Scenario: GET /v1/groups-report checking the response status 200 one group with member
    Given I am an admin user
    And group "test-group" and description "test group description" exists in the database
    And a user with username "abcd1234" and email "email@ons.gov.uk" exists in the database
    And user "abcd1234" is a member of group "test-group"
    When I GET "/v1/groups-report"
    Then the response header "Content-Type" should contain "application/json"
    And I should receive the following JSON response with status "200":
      """
      [
        {
          "group": "test group description",
          "user": "email@ons.gov.uk"
        }
      ]
      """

  Scenario: GET /v1/groups-report checking the response status 200 many groups with members
    Given I am an admin user
    And group "test-group_1" and description "test group_1 description" exists in the database
    And group "test-group_2" and description "test group_2 description" exists in the database
    And group "test-group_3" and description "test group_3 description" exists in the database
    And group "test-group_4" and description "test group_4 description" exists in the database
    And group "test-group_5" and description "test group_5 description" exists in the database

    And a user with username "abcd1234" and email "email1@ons.gov.uk" exists in the database
    And user "abcd1234" is a member of group "test-group_1"
    And user "abcd1234" is a member of group "test-group_2"
    And user "abcd1234" is a member of group "test-group_4"
    And user "abcd1234" is a member of group "test-group_5"

    And a user with username "abcd1235" and email "email2@ons.gov.uk" exists in the database
    And user "abcd1235" is a member of group "test-group_1"
    And user "abcd1235" is a member of group "test-group_2"

    And a user with username "abcd1236" and email "email3@ons.gov.uk" exists in the database
    And user "abcd1236" is a member of group "test-group_4"
    And user "abcd1236" is a member of group "test-group_5"

    When I GET "/v1/groups-report"
    Then I should receive the following JSON response with status "200":
      """
      [
        {
          "group": "test group_1 description",
          "user": "email1@ons.gov.uk"
        },
        {
          "group": "test group_1 description",
          "user": "email2@ons.gov.uk"
        },
        {
          "group": "test group_2 description",
          "user": "email1@ons.gov.uk"
        },
        {
          "group": "test group_2 description",
          "user": "email2@ons.gov.uk"
        },
        {
          "group": "test group_4 description",
          "user": "email1@ons.gov.uk"
        },
        {
          "group": "test group_4 description",
          "user": "email3@ons.gov.uk"
        },
        {
          "group": "test group_5 description",
          "user": "email1@ons.gov.uk"
        },
        {
          "group": "test group_5 description",
          "user": "email3@ons.gov.uk"
        }
      ]
      """


  Scenario: GET /v1/groups-report checking the response status 200 many groups with members
    Given I am an admin user

    And group "test-group_1" and description "test group_1 description" exists in the database
    And group "test-group_2" and description "test group_2 description" exists in the database
    And group "test-group_3" and description "test group_3 description" exists in the database
    And group "test-group_4" and description "test group_4 description" exists in the database
    And group "test-group_5" and description "test group_5 description" exists in the database

    And a user with username "abcd1234" and email "email1@ons.gov.uk" exists in the database
    And user "abcd1234" is a member of group "test-group_1"
    And user "abcd1234" is a member of group "test-group_2"
    And user "abcd1234" is a member of group "test-group_4"
    And user "abcd1234" is a member of group "test-group_5"

    And a user with username "abcd1235" and email "email2@ons.gov.uk" exists in the database
    And user "abcd1235" is a member of group "test-group_1"
    And user "abcd1235" is a member of group "test-group_2"

    And a user with username "abcd1236" and email "email3@ons.gov.uk" exists in the database
    And user "abcd1236" is a member of group "test-group_4"
    And user "abcd1236" is a member of group "test-group_5"
    And request header Accept is "text/csv"

    When I GET "/v1/groups-report"
    Then the HTTP status code should be "200"
    And the response should match the following csv:
      """
      Group,User
      test group_1 description,email1@ons.gov.uk
      test group_1 description,email2@ons.gov.uk
      test group_2 description,email1@ons.gov.uk
      test group_2 description,email2@ons.gov.uk
      test group_4 description,email1@ons.gov.uk
      test group_4 description,email3@ons.gov.uk
      test group_5 description,email1@ons.gov.uk
      test group_5 description,email3@ons.gov.uk
      """
    And the response header "Content-Type" should contain "text/csv"
