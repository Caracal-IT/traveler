CERT_NAME_CA="/dist/$1-ca"
CONFIG_TEMPLATE="/scripts/client/cert.cnf" #Ensure this is securely mapped
# Write temp/working config to a writable location inside the container
CONFIG="/tmp/new_cert.cnf"

cp "${CONFIG_TEMPLATE}" "${CONFIG}"

# Ensure the CN is properly set in the modified config
sed -i '/CN *=/d' "$CONFIG" # Remove existing CN lines
echo -e "\nCN = $1" >> "$CONFIG" # Append new CN line

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

mkdir -p "/dist/clients/$1/"
mv "${CERT_NAME_CA}".* "/dist/clients/$1/"
