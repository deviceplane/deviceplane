RETRIES=20
if [[ -z "$DEBUG" ]]; then
    until echo 'select 1' | "$@" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
        echo "Retrying db connection, $((RETRIES--)) retries left..."
        sleep 1
    done
else
    until echo 'select 1' | "$@" || [ $RETRIES -eq 0 ]; do
        echo "Retrying db connection, $((RETRIES--)) retries left..."
        sleep 1
    done
fi
