#!/bin/bash
user_pool_id=$1
profile=${2:-default}

re='^[0-9]+$'
oldIFS="$IFS"
IFS=$'\n'
totalCount=0
start=$(date +%s)

groups=$(aws cognito-idp list-groups --user-pool-id "${user_pool_id}"  --query 'Groups[*].GroupName' --profile "${profile}" | jq '.[]')
groups=($groups)
for group in "${groups[@]}"
do
    groupDetails=($(aws cognito-idp get-group --user-pool-id "${user_pool_id}" --group-name "${group:1:${#group}-2}"  --query 'Group.[Description]' --profile "${profile}" | jq '.[]'  ) )
    mygroupmembers=( $(aws cognito-idp list-users-in-group --user-pool-id "${user_pool_id}" --group-name "${group:1:${#group}-2}" --query 'Users[*].Username' --profile "${profile}" | jq '. | length') )
    let totalCount=totalCount+"${mygroupmembers}"

done

echo "Total group members = "${totalCount}
echo $(( $(date +%s) - start ))
IFS="$oldIFS"