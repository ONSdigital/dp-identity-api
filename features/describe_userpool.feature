Feature: DescribeUserPool

  Scenario: The user pool exists and can be reached
    Given user pool with id "us-west-2_aaaaaaaaa" exists
    When I GET "/userpool/us-west-2_aaaaaaaaa"
    Then I should receive the following JSON response with status "200":
    """
    { "message": "API and Cognito healthy"}
    """

  Scenario: The user pool does not exist
      When I GET "/userpool/us-west-2_bbbbbbbb"
      Then the HTTP status code should be "500"
