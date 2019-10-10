TIME=20
until echo 'select 1' | "$@" > /dev/null 2>&1 || [ $TIME -eq 0 ]; do
    echo "Waiting for postgres server, $((TIME--)) seconds left..."
    sleep 1
done