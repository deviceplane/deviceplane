PROD="prod"
STAGE="staging"
LOCAL="local"

if [ "$ENV" = "$PROD" ]; then
    if [[ -z "$DUMP" ]]; then
        echo ssh ubuntu@54.200.126.157 mysql -h deviceplane.coed7waagekn.us-west-2.rds.amazonaws.com -u deviceplane --password=$DB_PASS -P 3306 -D deviceplane
    else
        echo ssh ubuntu@54.200.126.157 mysqldump -h deviceplane.coed7waagekn.us-west-2.rds.amazonaws.com -u deviceplane --password=$DB_PASS -P 3306 --databases deviceplane
    fi
    exit 0
fi
if [ "$ENV" = "$STAGE" ]; then
    echo "no stage yet"
    exit 1
fi
if [ "$ENV" = "$LOCAL" ]; then
    echo mysql -h 127.0.0.1 -u user --password=pass -P 3306 -D deviceplane
    exit 0
fi

echo "ENV ($ENV) is not valid"
exit 1