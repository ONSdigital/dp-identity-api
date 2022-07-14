# Import users & groups

Import users and groups to Cognito from S3.

## Getting started

 Run ```go run scripts/import_users/import_users.go```

## Dependencies

No further dependencies other than following configuration.

## Configuration

| Environment variable         | Description
| ---------------------------- | -----------
| GROUPS_FILENAME              | Groups S3 backup filename
| GROUPUSERS_FILENAME          | User groups S3 backup filename
| FILENAME                     | User S3 backup filename
| S3_BUCKET                    | S3 bucket name
| S3_BASE_DIR                  | S3 backup DIR
| S3_REGION                    | S3 region name
| USER_POOL_ID                 | Cognito user pool id