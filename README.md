# dp-identity-api

An API used to manage the authorisation of users accessing data publishing services.

## Getting started

Run `make debug` to get started

To run the app with reflex and have it restart when you make changes, run:

`make debug-watch`

### Dummy data

If test data is required in the local Cognito user pool:

* Run `make populate-local`

To remove create test data from Cognito user pool:

* Run `make remove-test-data`

### Dependencies

* No further dependencies other than those defined in `go.mod`

To run make validate-specification you require Node v20.x and to install @redocly/cli:

```sh
   npm install -g @redocly/cli
```

### Configuration

| Environment variable         | Default   | Description                                                                                                        
|------------------------------|-----------|--------------------------------------------------------------------------------------------------------------------
| BIND_ADDR                    | :25600    | The host and port to bind to                                                                                       
| GRACEFUL_SHUTDOWN_TIMEOUT    | 20s       | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)                                                            
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) 
| AWS_REGION                   | eu-west-2 | The default AWS region for the identity api service                                                                
| AWS_COGNTIO_USER_POOL_ID     | -         | The ID of the user pool to be used                                                                                 
| AWS_COGNITO_CLIENT_ID        | -         | Cognito client ID                                                                                                  
| AWS_COGNITO_CLIENT_SECRET    | -         | Cognito client secret                                                                                              
| AWS_AUTH_FLOW                | -         | A parameter to define the request to the InitiateAuth endpoint in cognito                                          
| AWS_AUTH_FLOW                | -         | A parameter to define the request to the InitiateAuth endpoint in cognito                                          
| ENABLE_PLUS_EMAIL_BLOCKING   | true      | A feature flag to allow/disallow emails addresses with plus sign during user creation                              
| HTTP_WRITE_TIMEOUT           | [^dpnet]  | How long the dispatcher waits for us to write to it (`time.Duration` format)                                       

[^dpnet]: dp-net default

To get the values for the other AWS Cognito secrets:

* AWS_COGNITO_USER_POOL_ID get from AWS > Cognito > User Pools > user pool ID
* AWS_COGNITO_CLIENT_ID get from AWS > Cognito > User Pools > App Integration > App clients > dp-identity-api > client
  id
* AWS_COGNITO_CLIENT_SECRET get from AWS > Cognito > User Pools > App Integration > App clients > dp-identity-api >
  client secret

### Configuration needed to import user and group from s3

```sh
export GroupsFilename=""
export GroupUsersFilename=""
export UserFileName=""
export S3Bucket=""
export S3BaseDir=""
export S3Region=""
export AWSCognitoUserPoolID=""
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2024, [Office for National Statistics](https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
