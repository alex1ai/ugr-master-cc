#!/bin/bash

az ad sp create-for-rbac > tmp.json

export AZURE_TENANT_ID=$( cat tmp.json | jq '.tenant' | sed 's/"//g' )
export AZURE_CLIENT_ID=$( cat tmp.json | jq '.appId' | sed 's/"//g' )
export AZURE_CLIENT_SECRET=$( cat tmp.json | jq '.password' | sed 's/"//g' )

export AZURE_SUBSCRIPTION_ID=$(az account list --query "[?isDefault].id" -o tsv)

rm tmp.json
