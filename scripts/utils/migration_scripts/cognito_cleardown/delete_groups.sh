#!/bin/bash
user_pool_id=$1
profile=${2:-default}

roles=("role-admin" "role-publisher")
totalCount=0
start=$(date +%s)

groups=($(aws cognito-idp list-groups --user-pool-id "${user_pool_id}"  --query 'Groups[*].GroupName' --profile "${profile}" | jq -r '.[]') )
for group in "${groups[@]}"
do
    if [[ ! "${roles[*]}" =~ ${group} ]]; then
        echo "${group}"
        aws cognito-idp delete-group --user-pool-id "${user_pool_id}" --group-name ${group} --profile "${profile}"
    fi
done

echo $(( $(date +%s) - start ))