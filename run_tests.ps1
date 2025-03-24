# Configuration
$SIM_NUMBERS = 1000
$TIMEOUT = 0
$LOG_FILE = "./tests.log"

# Delete tests log file
if (Test-Path $LOG_FILE)
{
    Remove-Item $LOG_FILE -Force
}

# Arguments
if ($args.Count -ge 1)
{
    # Check if first argument is a number
    if ($args[0] -match '^\d+$')
    {
        $SIM_NUMBERS = [int]$args[0]
    }
    else
    {
        $ERROR_MESSAGE = "Invalid argument: must be a positive integer, defaulting simulation numbers to $SIM_NUMBERS"
        $ERROR_MESSAGE | Tee-Object -FilePath $LOG_FILE -Append
    }
}

# Clear cache
"Cleaning golang cache..." | Tee-Object -FilePath $LOG_FILE -Append
go clean -cache | Tee-Object -FilePath $LOG_FILE -Append

# Run tests
"Testing $SIM_NUMBERS simulations per test..." | Tee-Object -FilePath $LOG_FILE -Append
go test -cover ./... -sim="$SIM_NUMBERS" -timeout="$TIMEOUT" -v >> $LOG_FILE
"Tests logged to '$LOG_FILE'" | Tee-Object -FilePath $LOG_FILE -Append