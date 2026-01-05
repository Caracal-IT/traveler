#!/usr/bin/env sh

# Create client certificates for caracal-client for several profiles, skipping any that already exist
CLIENT="caracal-server"
SCRIPT="/scripts/certs/create-cert.sh"
PROFILES="default-server-dev"

call_create() {
  profile="$1"
  sh "$SCRIPT" "server" "$CLIENT" "$profile"
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
