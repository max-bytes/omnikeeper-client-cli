#!/bin/bash

KEYCLOAK_URL=https://auth-dev.mhx.at/auth
KEYCLOAK_REALM=acme
KEYCLOAK_CLIENT=landscape-omnikeeper
KEYCLOAK_USER=omnikeeper-client-library-test
KEYCLOAK_PASSWORD=omnikeeper-client-library-test
GRAPHQL_URL=https://10.0.0.43:45455/graphql

ACCESS_TOKEN=$(curl -s $KEYCLOAK_URL/realms/$KEYCLOAK_REALM/protocol/openid-connect/token \
    -d client_id=$KEYCLOAK_CLIENT \
    -d grant_type=password \
    -d username=$KEYCLOAK_USER \
    -d password=$KEYCLOAK_PASSWORD\
    | jq -r '.access_token')

# store the access token AND the graphql URL in env variables for future use
export OMNIKEEPER_ACCESS_TOKEN="$ACCESS_TOKEN"
export OMNIKEEPER_GRAPHQL_URL="$GRAPHQL_URL"

echo "Logged in at $KEYCLOAK_URL, access token stored"