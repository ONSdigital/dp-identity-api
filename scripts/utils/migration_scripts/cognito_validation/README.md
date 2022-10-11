dp-identity-api Cognito validation
================

### Purpose

These processes are to be used validate the number of users, teams/groups and user team/group membership  in a cognito userpool.  
If using an awsb environment you need to aws `sso login --profile <environment aws profile name>`  
To view and amend execution permissions  
`ls –l [file_name]`  
`chmod +x filename`  

### Number of Users
to run 
`./number_of_users.sh <userpool> <environment aws profile name>`

### Number of groups
to run
`./number_of_groups.sh <userpool> <environment aws profile name>`

### Number of groups for users
to run
`./number_of_users_teams.sh <userpool> <environment aws profile name>`