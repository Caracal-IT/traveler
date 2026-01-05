#!/usr/bin/env sh

# Create client certificates for caracal-client for several profiles, skipping any that already exist
client="caracal-client"
SCRIPT_IN_CONTAINER="/scripts/client/create-client.sh"
LOCAL_SCRIPT="$(dirname "$0")/create-client.sh"
PROFILES="default-client-dev default-client-iot default-client-mqtt"

call_create() {
  profile="$1"
  if [ -x "$SCRIPT_IN_CONTAINER" ]; then
    "$SCRIPT_IN_CONTAINER" "$client" "$profile"
  elif [ -f "$SCRIPT_IN_CONTAINER" ]; then
    sh "$SCRIPT_IN_CONTAINER" "$client" "$profile"
  elif [ -x "$LOCAL_SCRIPT" ]; then
    "$LOCAL_SCRIPT" "$client" "$profile"
  elif [ -f "$LOCAL_SCRIPT" ]; then
    sh "$LOCAL_SCRIPT" "$client" "$profile"
  else
    echo "Error: no create-client.sh found at $SCRIPT_IN_CONTAINER or $LOCAL_SCRIPT" >&2
    return 1
  fi
}

for profile in $PROFILES; do
  target_dir="/dist/clients/$client/$profile"
  if [ -d "$target_dir" ]; then
    echo "Skipping creation: $target_dir already exists"
    continue
  fi

  echo "Creating client certs for $client / $profile"
  call_create "$profile" || { echo "Failed to create certificates for $profile" >&2; exit 1; }
done
