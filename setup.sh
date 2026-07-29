#!/bin/sh

y="y"
yes="yes"
n="n"
no="no"

echo "Creating loselose-game"
go build ./cmd/loselose-game/

echo "Creating injector"
go build ./cmd/loselose-injection

read -p "Would you like to modify /proc/sys/kernel/yama/ptrace_scope to increase stability (Not required)? [Y/N] " choice
if [[ ${choice,,} == $y  ||  ${choice,,} == $yes ]]; then
    echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
    echo "Modified /proc/sys/kernel/yama/ptrace_scope to 0"
elif [[ ${choice,,} == $n  ||  ${choice,,} == $no ]]; then
    echo "No modification done."
else
    echo "Invalid option given, proceeding with no modification."
fi

echo "Moving loselose-game to staging location"
mv ./loselose-game /tmp/loselose

echo "Setup completed, good luck! :)"