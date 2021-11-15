#!/bin/bash

mode="execute"
while getopts lu:p:k:r:c:g:q: flag
do
    case "${flag}" in
        l) mode="login";;
        u) keycloak_user=${OPTARG};;
        p) keycloak_password=${OPTARG};;
        k) keycloak_url=${OPTARG};;
        r) keycloak_realm=${OPTARG};;
        c) keycloak_client=${OPTARG};;
        g) graphql_url=${OPTARG};;
        q) graphql_query=${OPTARG};;
    esac
done

echo "mode: $mode";

if [ "$mode" == "login" ]; then

    # reset stored token
    export OMNIKEEPER_ACCESS_TOKEN=""
    export OMNIKEEPER_GRAPHQL_URL=""

    if [[ -z "${keycloak_user}" ]]; then
        echo "Keycloak username (-u) not set"
        exit 1
    fi
    if [[ -z "${keycloak_password}" ]]; then
        echo "Keycloak password (-p) not set"
        exit 1
    fi
    if [[ -z "${keycloak_url}" ]]; then
        echo "Keycloak URL (-k) not set"
        exit 1
    fi
    if [[ -z "${keycloak_realm}" ]]; then
        echo "Keycloak realm (-r) not set"
        exit 1
    fi
    if [[ -z "${keycloak_client}" ]]; then
        echo "Keycloak client (-c) not set"
        exit 1
    fi
    if [[ -z "${graphql_url}" ]]; then
        echo "GraphQL URL (-g) not set"
        exit 1
    fi

    ACCESS_TOKEN=$(curl -s $keycloak_url/realms/$keycloak_realm/protocol/openid-connect/token \
        -d client_id=$keycloak_client \
        -d grant_type=password \
        -d username=$keycloak_user \
        -d password=$keycloak_password\
        | jq -r '.access_token')

    # store the access token AND the graphql URL in env variables for future use
    export OMNIKEEPER_ACCESS_TOKEN="$ACCESS_TOKEN"
    export OMNIKEEPER_GRAPHQL_URL="$graphql_url"
else 
    if [[ -z "${OMNIKEEPER_GRAPHQL_URL}" ]]; then
        echo "Environment variable OMNIKEEPER_GRAPHQL_URL not set. Did you login?"
        exit 1
    fi

    if [[ -z "${OMNIKEEPER_ACCESS_TOKEN}" ]]; then
        echo "Environment variable OMNIKEEPER_ACCESS_TOKEN not set. Did you login?"
        exit 1
    fi

    echo "$OMNIKEEPER_GRAPHQL_URL"

    curl $OMNIKEEPER_GRAPHQL_URL -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $OMNIKEEPER_ACCESS_TOKEN" \
        -d "{\"query\":\"$graphql_query\"}" -v
fi
