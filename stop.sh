#!/bin/bash

pids=$(pgrep -f "wechat_pay_server")

if [ -z "$pids" ]; then
    echo "No processes found with the name 'wechat_pay_server'."
else
    echo "Killing processes with IDs: $pids"

    pkill -f "wechat_pay_server"

    if [ $? -eq 0 ]; then
        echo "Processes with names containing 'wechat_pay_server' have been killed successfully."
    else
        echo "Failed to kill processes with names containing 'wechat_pay_server'."
    fi
fi

exit 0