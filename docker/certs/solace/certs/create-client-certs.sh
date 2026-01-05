#!/usr/bin/env sh

# Create client certificates for caracal-client for several profiles, skipping any that already exist
CLIENT="caracal-client"
SCRIPT="/scripts/certs/create-cert.sh"
PROFILES="default-client-dev default-client-iot default-client-mqtt"

call_create() {
  profile="$1"
  sh "$SCRIPT" "client" "$CLIENT" "$profile"
}

for profile in $PROFILES; do
  target_dir="/dist/$CLIENT/$profile"

  if [ -d "$target_dir" ]; then
    echo "Skipping creation: $target_dir already exists"
    continue
  fi

  echo "Creating client certs for $CLIENT / $profile"
  call_create "$profile" || { echo "Failed to create certificates for $profile" >&2; exit 1; }
done
