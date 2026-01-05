DEST="$2"
CATEGORY="$1"
CONFIG_TEMPLATE="/scripts/certs/cert.cnf" #Ensure this is securely mapped
CONFIG="/tmp/new_cert.cnf" # Write temp/working config to a writable location inside the container
CERT_NAME_CA="/dist/$DEST-ca"

cp "${CONFIG_TEMPLATE}" "${CONFIG}"

# Ensure the CN is properly set in the modified config
sed -i '/CN *=/d' "$CONFIG" # Remove existing CN lines
echo -e "\nCN = $DEST" >> "$CONFIG" # Append new CN line

# Exit on error
set -e

echo "-------------------------------------------------------------"
echo " Generate CA private key and self-signed certificate"
echo "-------------------------------------------------------------"

openssl req                           \
  -new                                \
  -x509                               \
  -days 365                           \
  -extensions v3_ca                   \
  -keyout "$CERT_NAME_CA.key"         \
  -out "$CERT_NAME_CA.pem"            \
  -config "$CONFIG"                   \
  -passin pass:"$BROKER_PASSWORD"     \
  -passout pass:"$BROKER_PASSWORD"    \

echo "-------------------------------------------------------------"
echo " Convert CA certificate to CRT format"
echo "-------------------------------------------------------------"

openssl x509 -outform der -in "$CERT_NAME_CA.pem" -out "$CERT_NAME_CA.crt"

# Cleanup only temporary files; be tolerant if they don't exist
rm -f ${CONFIG}

echo "-------------------------------------------------------------"
echo " Copy CA certificate and key to client dist folder"
echo "-------------------------------------------------------------"

mkdir -p "/dist/$CATEGORY/$DEST/ca"
mv "${CERT_NAME_CA}".* "/dist/$CATEGORY/$DEST/ca"
