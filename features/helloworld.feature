Feature: Helloworld



  Scenario: Posting and checking a response
    When I GET "/hello"
    Then I should receive a hello-world response

  Given I set the "Authorization" header to "something"
    When I GET "/hello"
    Then I should receive the following JSON response with status "200":
    """
    { "message": "Hello, World!"}
    """