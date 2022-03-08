Feature: JWKS

#   Get JSON web key set with JWT Key ID and associated RSA Public Signing Key
    Scenario: GET /v1/jwt-keys and checking the response status 200
        Given I am an admin user
        When I GET "/v1/jwt-keys"
        Then I should receive the following JSON response with status "200":
            """
            {
                "2a8vXmIK67ZZ3hFZ=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApU1DqxJalEmlznkrM+U4aSBMw9u5axcIqNeUq8+ZHo98uKy8Xy5zCOXfWx6KafOPJhbOZInFaSh9UMaluSSw11l/PR4KrGBFzJODQ+RMq6bHW6FlwwHSkMTSfQ0hwzO7y91BiZFmJnaUECf52H3QBApGT4TT060ri5zt1ygpliRwjLLlHW1XX0epzZH3ogrikn4i65e8w6uUcsGBhQvQQqiHvEpcgCQAB",
				"GRBevIroJzPBvaGa=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtvDfudfY9n+8sFJmHGFfgbKqKf8iiEcbvRXNMEi9qd2NGAekhdNJKdeW3sMSwR+sb4Ly6IypowCE2eueYk/GatzYyyolWny/Krdp0EWPT/PnK8Iq1FTIuHxFb08B8iLnH/2nKqgOjVvwEU4eSBh0YHKti2v77a+a4bnx6aOC2YkF2AyIRmbXAHaq4Js9u33X8gGMXZcVsxcSpG8Py/NJ3s+PLKebQFQAB"
            }
            """
