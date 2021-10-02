#!/bin/bash

# ======================================
#  birdol-API automated testing tool
#  [Dependencies]: nodejs, newman, jq
# ======================================

# Generate RSA Key
echo "Generating RSA Key..."
PUBKEY_A=`./rsa-sign/signing-cs gen rsa-sign/keyfile/key-a`
PUBKEY_B=`./rsa-sign/signing-cs gen rsa-sign/keyfile/key-b`

# Generate UUID
echo "Generating UUID..."
UUID_A=`./rsa-sign/signing-cs uuid`
UUID_B=`./rsa-sign/signing-cs uuid`

echo "Starting newman test process..."

# Run Test 1: CreateAccount - DeviceA
newman run collection.json -e environment.json --folder CreateAccount --env-var "PUBLIC_KEY_A=$PUBKEY_A" --env-var "DEVICE_ID_A=$UUID_A" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 1"
    exit 
fi

# Run Test 2: LoginWithToken - DeviceA
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:" rsa-sign/keyfile/key-a.xml`

newman run collection.json -e environment_tmp.json --folder LoginWithToken --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
echo "Fail: Test 2"
    exit
fi

# Run Test 3: SetDataLink - DeviceA
PASS=`cat environment.json | jq -r '.values | from_entries | .LINK_PASSWORD'`
BODY=`cat collection.json | jq -r '.item[].item[0].item[2].request.body.raw'`
BODY=`echo "$BODY" | sed s/{{LINK_PASSWORD}}/$PASS/`

TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:$BODY" rsa-sign/keyfile/key-a.xml`

newman run collection.json -e environment_tmp.json --folder SetDataLink --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 3"
    exit
fi

# Run Test 4: Refresh Token - DeviceA
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:" rsa-sign/keyfile/key-a.xml`

newman run collection.json -e environment_tmp.json --folder RefreshToken --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 4"
    exit
fi

# Run Test 5: AccountLink - DeviceB
newman run collection.json -e environment_tmp.json --folder AccountLink --env-var "PUBLIC_KEY_B=$PUBKEY_B" --env-var "DEVICE_ID_B=$UUID_B" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 5"
    exit
fi

# Run Test 6: LoginWithToken - DeviceB
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:" rsa-sign/keyfile/key-b.xml`

newman run collection.json -e environment_tmp.json --folder LoginWithTokenB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 6"
    exit
fi

# Run Test 7: SetDataLink - DeviceB
PASS=`cat environment.json | jq -r '.values | from_entries | .LINK_PASSWORD'`
BODY=`cat collection.json | jq -r '.item[].item[1].item[2].request.body.raw'`
BODY=`echo "$BODY" | sed s/{{LINK_PASSWORD}}/$PASS/`

TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:$BODY" rsa-sign/keyfile/key-b.xml`

newman run collection.json -e environment_tmp.json --folder SetDataLinkB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 7"
    exit
fi

# Run Test 8: Refresh Token - DeviceB
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`./rsa-sign/signing-cs sign "v2:$TS:" rsa-sign/keyfile/key-b.xml`

newman run collection.json -e environment_tmp.json --folder RefreshTokenB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 8"
    exit
fi

echo "All test passed."
echo "Cleaning..."
rm environment_tmp.json ./rsa-sign/keyfile/*
