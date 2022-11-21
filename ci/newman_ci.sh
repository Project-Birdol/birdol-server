#!/bin/bash

# =================================================
#  birdol-API automated testing tool
#  [Dependencies]: dotnet-sdk, nodejs, newman, jq
# =================================================

function clean() {
    echo "Cleaning..."
    rm environment_tmp.json
    return 0
}

# Generate RSA Key
echo "Generating RSA Key..."
PUBKEY_A=`dotnet run -c Release --project verify-test-tool gen /tmp/key-a`
PUBKEY_B=`dotnet run -c Release --project verify-test-tool gen /tmp/key-b`

# Generate UUID
echo "Generating UUID..."
UUID_A=`dotnet run -c Release --project verify-test-tool uuid`
UUID_B=`dotnet run -c Release --project verify-test-tool uuid`

echo "Starting newman test process..."

# Run Test 1: CreateAccount - DeviceA
newman run collection.json -e environment.json --folder CreateAccount --env-var "PUBLIC_KEY_A=$PUBKEY_A" --env-var "DEVICE_ID_A=$UUID_A" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 1"
    clean
    exit 1
fi

# Run Test 2: LoginWithToken - DeviceA
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet -c Release --project verify-test-tool sign "v2:$TS:" /tmp/key-a.xml`

newman run collection.json -e environment_tmp.json --folder LoginWithToken --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 2"
    clean
    exit 1
fi

# Run Test 3: SetDataLink - DeviceA
PASS=`cat environment.json | jq -r '.values | from_entries | .LINK_PASSWORD'`
BODY=`cat collection.json | jq -r '.item[].item[0].item[2].request.body.raw'`
BODY=`echo "$BODY" | sed s/{{LINK_PASSWORD}}/$PASS/`

TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet run -c Release --project verify-test-tool sign "v2:$TS:$BODY" /tmp/key-a.xml`

newman run collection.json -e environment_tmp.json --folder SetDataLink --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 3"
    clean
    exit 1
fi

# Run Test 4: Refresh Token - DeviceA
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet run -c Release --project verify-test-tool sign "v2:$TS:" /tmp/key-a.xml`

newman run collection.json -e environment_tmp.json --folder RefreshToken --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 4"
    clean
    exit 1
fi

# Run Test 5: AccountLink - DeviceB
newman run collection.json -e environment_tmp.json --folder AccountLink --env-var "PUBLIC_KEY_B=$PUBKEY_B" --env-var "DEVICE_ID_B=$UUID_B" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 5"
    clean
    exit 1
fi

# Run Test 6: LoginWithToken - DeviceB
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet run -c Release --project verify-test-tool sign "v2:$TS:" /tmp/key-b.xml`

newman run collection.json -e environment_tmp.json --folder LoginWithTokenB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 6"
    clean
    exit 1
fi

# Run Test 7: SetDataLink - DeviceB
PASS=`cat environment.json | jq -r '.values | from_entries | .LINK_PASSWORD'`
BODY=`cat collection.json | jq -r '.item[].item[1].item[2].request.body.raw'`
BODY=`echo "$BODY" | sed s/{{LINK_PASSWORD}}/$PASS/`

TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet run -c Release --project verify-test-tool sign "v2:$TS:$BODY" /tmp/key-b.xml`

newman run collection.json -e environment_tmp.json --folder SetDataLinkB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 7"
    clean
    exit 1
fi

# Run Test 8: Refresh Token - DeviceB
TS=`date '+%Y-%m-%d-%H-%M-%S'`
SIG=`dotnet run -c Release --project verify-test-tool sign "v2:$TS:" /tmp/key-b.xml`

newman run collection.json -e environment_tmp.json --folder RefreshTokenB --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" --export-environment environment_tmp.json

if [ $? -ne 0 ]; then
    echo "Fail: Test 8"
    clean
    exit 1
fi

echo "All test passed."
clean
