#!/bin/bash
user_pool_id=$1
profile=${2:-default}


totalCount=0
start=$(date +%s)


users=($(aws cognito-idp list-users --user-pool-id "${user_pool_id}"  --query 'Users[*].Username' --profile "${profile}" | jq -r '.[]'))

for user in "${users[@]}"
do
    echo "${user}"
    aws cognito-idp admin-delete-user --user-pool-id "${user_pool_id}" --username ${user} --profile "${profile}"

done

echo $(( $(date +%s) - start ))