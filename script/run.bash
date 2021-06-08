#!/bin/bash

set -e
set -x

# Run
bin/controller --addr ":80" \
               --mysql "$DP_DB_CONN" \
               --email-provider smtp \
               --email-from-address noreply@deviceplane.com \
               --smtp-server "$DP_SMTP_SERVER" \
               --smtp-port 465 \
               --smtp-username apikey \
               --smtp-password "$DP_SMTP_PW" \
               --auth0-audience "$DP_AUTH0_AUD" \
               --auth0-domain "$DP_AUTH0_DOM" \
               --db-max-open-conns 5 \
               --db-max-idle-conns 5 \
               --db-max-conn-lifetime 5m                
