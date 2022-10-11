#!/bin/bash
start=$(date +%s)
user_pool_id=$1
profile=${2:-default}

echo $(aws cognito-idp list-groups --user-pool-id "${user_pool_id}"  --query 'Groups[*].GroupName' --profile "${profile}" | jq '. | length')
echo $(( $(date +%s) - start ))