#!/bin/sh
set -e
cd testdata
for task in $(dog | awk '{print $1}'); do
	echo -e "\n## Test: $task"
	dog $task || true
done
