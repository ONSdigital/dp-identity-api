#!/bin/bash
user_pool_id=$1
profile=${2:-default}

chmod +x delete_users_by_group.sh
chmod +x delete_groups.sh
chmod +x delete_users.sh

./delete_users_by_group.sh "$user_pool_id" "$profile"
./delete_groups.sh "$user_pool_id" "$profile"
./delete_users.sh  "$user_pool_id" "$profile"

echo $(( $(date +%s) - start ))