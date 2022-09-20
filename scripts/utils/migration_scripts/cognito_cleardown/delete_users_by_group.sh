#!/bin/bash
user_pool_id=$1
profile=${2:-default}

re='^[0-9]+$'
oldIFS="$IFS"
IFS=$'\n'
totalCount=0
start=$(date +%s)

groups=( $(aws cognito-idp list-groups --user-pool-id "${user_pool_id}"  --query 'Groups[*].GroupName' --profile "${profile}" | jq '.[]')  )
echo "-------"
for group in "${groups[@]}"
do
    echo "group" "${group}"
    echo "modified group" "${group:1:${#group}-2}"
    mygroupmembers=( $(aws cognito-idp list-users-in-group --user-pool-id "${user_pool_id}" --group-name "${group:1:${#group}-2}" --query 'Users[*].Username' --profile "${profile}" | jq '.[]'  ))
    echo  "${mygroupmembers}"
    if [ ${#mygroupmembers[@]} >  0 ] ; then
        for gm in "${mygroupmembers[@]}"
        do
            echo $group $gm
            aws cognito-idp admin-remove-user-from-group --user-pool-id "${user_pool_id}" --group-name "${group:1:${#group}-2}" --username "${gm:1:${#gm}-2}" --profile "${profile}"
        done
    fi
done
echo "-------"
echo $(( $(date +%s) - start ))
IFS="$oldIFS"
