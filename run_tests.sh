#!/usr/bin/env bash

# Configuration
SIM_NUMBERS=1000
TIMEOUT="0" # disabled timeout
LOG_FILE='./tests.log'

# Delete tests log file
[ -f ./tests.log ] && rm ./tests.log

# Arguments
if [[ "$#" -ge 1 ]]; then
  # Check if first argument is a number
  if [[ "$1" =~ ^[0-9]+$ ]]; then
    SIM_NUMBERS=$1
  else
    ERROR="Invalid argument: must be a positive integer, defaulting simulation numbers to $SIM_NUMBERS"
    echo "$ERROR" | tee -a $LOG_FILE
  fi
fi

# Clear cache
echo "Cleaning golang cache..." | tee -a $LOG_FILE
go clean -cache | tee -a $LOG_FILE

# Run tests
echo "Testing $SIM_NUMBERS simulations per test..." | tee -a $LOG_FILE
go test -cover ./... -sim="$SIM_NUMBERS" -timeout="$TIMEOUT" -v >> $LOG_FILE
echo "Tests logged to '$LOG_FILE'" | tee -a $LOG_FILE