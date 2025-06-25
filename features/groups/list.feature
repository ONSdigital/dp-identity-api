@Groups @GroupsList
Feature: Groups - List
  Scenario: GET /v1/groups with 0 groups and checking the response status 200
    Given there are 0 groups in the database
    And I am an admin user
    When I GET "/v1/groups"
    Then I should receive the following JSON response with status "200":
      """
      {
        "groups": null,
        "count": 0,
        "next_token": null
      }
      """

  Scenario: GET /v1/groups with >0 groups and checking the response status 200
    Given there are 2 groups in the database
    And I am an admin user
    When I GET "/v1/groups"
    Then the response code should be 200
    And the response should match the following json for listgroups
      """
      {
        "count": 2,
        "groups": [
          {
            "name": "group name description 1",
            "id": "group_name_1",
            "precedence": 55
          },
          {
            "name": "group name description 2",
            "id": "group_name_2",
            "precedence": 55
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups with >60 groups and checking the response status 200
    Given there are 100 groups in the database
    And I am an admin user
    When I GET "/v1/groups"
    Then the response code should be 200
    And the response should match the following json for listgroups
      """
      {
        "count": 100,
        "groups": [
          {
            "name": "group name description 1",
            "id": "group_name_1",
            "precedence": 55
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups?sort=name:asc and checking the response status 200
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=name:asc"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 3,
        "groups": [
          {
            "name": "A Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "B Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "C Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups?sort=name:desc and checking the response status 200
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=name:desc"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 3,
        "groups": [
          {
            "name": "C Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "B Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "A Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups?sort=name and checking the response status 200
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=name"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 3,
        "groups": [
          {
            "name": "A Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "B Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "C Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups?sort=created and checking the response status 200
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=created"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 3,
        "groups": [
          {
            "name": "B Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "A Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          },
          {
            "name": "C Group",
            "id": "",
            "creation_date": "2010-01-01T00:00:00Z",
            "last_modified_date": "2010-01-01T00:00:00Z",
            "precedence": 1,
            "role_arn": "",
            "user_pool_id": ""
          }
        ],
        "next_token": null
      }
      """

  Scenario: GET /v1/groups?sort=abc and checking the response status 400
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=abc"
    Then the HTTP status code should be "400"

  Scenario: GET /v1/groups?sort=name:xyz and checking the response status 400
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=name:xyz"
    Then the HTTP status code should be "400"

  Scenario: GET /v1/groups?sort=abc:asc and checking the response status 400
    Given group "B Group" exists in a list in the database
    And group "A Group" exists in a list in the database
    And group "C Group" exists in a list in the database
    And I am an admin user
    When I GET "/v1/groups?sort=abc:asc"
    Then the HTTP status code should be "400"
