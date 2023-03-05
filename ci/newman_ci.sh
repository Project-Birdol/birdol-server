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

echo "######################### Birdol API testing tool #########################"
echo

TOOL_BIN="./verify-test-tool/bin/Release/net7.0/linux-x64/birdolcrypto"

# Generate RSA Key
echo "[Key Generation Benchmark]"
echo ">> Generating ECDSA Key..."
PUBKEY_A=$($TOOL_BIN keygen /tmp/key-a ecdsa | base64 -w 0)
PUBKEY_B=$($TOOL_BIN keygen /tmp/key-b ecdsa | base64 -w 0)

echo "[API Testing Process]"

# Generate UUID
echo ">> Generating UUID..."
UUID_A=$("$TOOL_BIN" uuid)
UUID_B=$("$TOOL_BIN" uuid)

# Run Test 1: CreateAccount - DeviceA
echo "> Test 1: CreateAccount - DeviceA"
newman run collection.json -e environment.json \
    --folder CreateAccount \
    --env-var "PUBLIC_KEY_A=$PUBKEY_A" --env-var "DEVICE_ID_A=$UUID_A" \
    --export-environment environment_tmp.json || { echo "Fail: Test 1"; clean; exit 1; }

echo ">> Passed"

# Run Test 2: LoginWithToken - DeviceA
echo "> Test 2: LoginWithToken - DeviceA"
TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$("$TOOL_BIN" sign "v2:$TS:" ecdsa /tmp/key-a.priv | base64 -w 0)
echo "$SIG" | base64 -d

newman run collection.json -e environment_tmp.json \
    --folder LoginWithToken \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 2"; clean; exit 1;}

echo ">> Passed"

# Run Test 3: SetDataLink - DeviceA
echo "> Test 3: SetDataLink - DeviceA"
PASS=$(jq -r '.values | from_entries | .LINK_PASSWORD' < environment.json)
BODY=$(jq -r '.item[].item[0].item[2].request.body.raw' < collection.json)
BODY=${BODY//"{{LINK_PASSWORD}}/$PASS"}

TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$("$TOOL_BIN" sign "v2:$TS:$BODY" ecdsa /tmp/key-a.priv | base64 -w 0)

newman run collection.json -e environment_tmp.json \
    --folder SetDataLink \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 3"; clean; exit 1; }

echo ">> Passed"

# Run Test 4: Refresh Token - DeviceA
echo "> Test 4: Refresh Token - DeviceA"
TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$("$TOOL_BIN" sign "v2:$TS:" ecdsa /tmp/key-a.priv | base64 -w 0)

newman run collection.json -e environment_tmp.json \
    --folder RefreshToken \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 4"; clean; exit 1; }

echo ">> Passed"

# Run Test 5: AccountLink - DeviceB
echo "> Test 5: AccountLink - DeviceB"
newman run collection.json -e environment_tmp.json \
    --folder AccountLink \
    --env-var "PUBLIC_KEY_B=$PUBKEY_B" --env-var "DEVICE_ID_B=$UUID_B" \
    --export-environment environment_tmp.json || { echo "Fail: Test 5"; clean; exit 1; }

echo ">> Passed"

# Run Test 6: LoginWithToken - DeviceB
echo "> Test 6: LoginWithToken - DeviceB"
TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$("$TOOL_BIN" sign ecdsa "v2:$TS:" /tmp/key-b.priv | base64 -w 0)

newman run collection.json -e environment_tmp.json \
    --folder LoginWithTokenB \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 6"; clean; exit 1; }

echo ">> Passed"

# Run Test 7: SetDataLink - DeviceB
echo "> Test 7: SetDataLink - DeviceB"
PASS=$(jq -r '.values | from_entries | .LINK_PASSWORD' < environment.json)
BODY=$(jq -r '.item[].item[1].item[2].request.body.raw' < collection.json)
BODY=${BODY//"{{LINK_PASSWORD}}/$PASS"}

TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$("$TOOL_BIN" sign "v2:$TS:$BODY" ecdsa /tmp/key-b.priv | base64 -w 0)

newman run collection.json -e environment_tmp.json \
    --folder SetDataLinkB \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 7"; clean; exit 1; }

echo ">> Passed"

# Run Test 8: Refresh Token - DeviceB
echo "> Test 8: Refresh Token - DeviceB"
TS=$(date '+%Y-%m-%d-%H-%M-%S')
SIG=$($TOOL_BIN sign "v2:$TS:" ecdsa /tmp/key-b.priv | base64 -w 0)

newman run collection.json -e environment_tmp.json \
    --folder RefreshTokenB \
    --env-var "SIGNATURE=$SIG" --env-var "TIMESTAMP=$TS" \
    --export-environment environment_tmp.json || { echo "Fail: Test 8"; clean; exit 1; }

echo ">> Passed"

echo "All test passed."
clean
