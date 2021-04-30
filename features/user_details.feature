Feature: DescribeUserPool

  Scenario: The user pool exists and can be reached
    Given a user with username "aaaa-bbbb-cccc-dddd" exists
    When I GET "/users/aaaa-bbbb-cccc-dddd"
    Then I should receive the following JSON response with status "200":
    """
    {"Enabled":true, "MFAOptions":[], "PreferredMfaSetting":null, "UserAttributes":[{"Name":"email", "Value":"user@ons.gov.uk"}, {"Name":"given_name", "Value":"Jane"}, {"Name":"family_name", "Value":"Doe"}], "UserCreateDate":"2021-01-01T12:00:00Z", "UserLastModifiedDate":"2021-01-01T12:00:00Z", "UserMFASettingList":[], "UserStatus":"CONFIRMED", "Username":"aaaa-bbbb-cccc-dddd"}
    """