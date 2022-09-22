dp-identity-api Cognito Clear down
================

### Purpose

These processes are to be used to make an existing cognito userpool back to a clean state.  
If using an awsb environment you need to...  
`aws sso login --profile <environment aws profile name>`  
To view and amend execution permissions  
`ls –l [file_name]`  
`chmod +x filename`  

Please run in the given order...

### Step 1 remove all user team/group membership
to run 
`./delete_users_by_group.sh <userpool> <environment aws profile name>`

### Step 2 remove all teams (except teams that are roles)
to run
`./delete_groups.sh <userpool> <environment aws profile name>`

### Step 3 remove all users 
to run
`./delete_users.sh <userpool> <environment aws profile name>`


### Complete run...  
`./delete_cognito_objects.sh
<userpool> <environment aws profile name>`

please note that these will continue if there is an issue in any  of the above steps 
