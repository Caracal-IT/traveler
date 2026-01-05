DEST="/dist/$1/$2/$3/"
CATEGORY="$1"
CERT_NAME="/dist/$3"
CERT_NAME_CA="/dist/$CATEGORY/$2/ca/$2-ca"

CONFIG_TEMPLATE="/scripts/certs/cert.cnf" #Ensure this is securely mapped
CONFIG="/tmp/new_cert.cnf" # Write temp/working config to a writable location inside the container

cp "${CONFIG_TEMPLATE}" "${CONFIG}"

# Ensure the CN is properly set in the modified config
sed -i '/CN *=/d' "$CONFIG" # Remove existing CN lines
echo -e "\nCN = $3" >> "$CONFIG" # Append new CN line

# Exit on error
set -e

echo "-------------------------------------------------------------"
echo " Generate certificate signing request (CSR) and private key"
echo "-------------------------------------------------------------"

openssl req                           \
  -new                                \
  -nodes                              \
  -days 365                           \
  -keyout "$CERT_NAME.key"            \
  -out "$CERT_NAME.csr"               \
  -config "$CONFIG"                   \
  -passin pass:"$BROKER_PASSWORD"     \

echo "-------------------------------------------------------------"
echo " Signing the certificate with the CA"
echo "-------------------------------------------------------------"

openssl x509                         \
  -req                               \
  -in "$CERT_NAME.csr"               \
  -CA "$CERT_NAME_CA.pem"            \
  -CAkey "$CERT_NAME_CA.key"         \
  -CAcreateserial                    \
  -out "$CERT_NAME.crt"              \
  -days 365                          \
  -passin pass:"$BROKER_PASSWORD"

echo "-------------------------------------------------------------"
echo " Converting to PFX format"
echo "-------------------------------------------------------------"

openssl pkcs12                        \
  -export                             \
  -out "$CERT_NAME.pfx"               \
  -inkey "$CERT_NAME.key"             \
  -in "$CERT_NAME.crt"                \
  -passout pass:"$BROKER_PASSWORD"

echo "-------------------------------------------------------------"
echo " Extracting RSA private key"
echo "-------------------------------------------------------------"

openssl pkcs12                        \
  -in "$CERT_NAME.pfx"                \
  -nocerts                            \
  -nodes                              \
  -out "$CERT_NAME.rsa"               \
  -passin pass:"$BROKER_PASSWORD"

echo "-------------------------------------------------------------"
echo " Creating PEM format"
echo "-------------------------------------------------------------"

openssl pkcs12                        \
  -in "$CERT_NAME.pfx"                \
  -out "$CERT_NAME.pem"               \
  -nodes                              \
  -passin pass:"$BROKER_PASSWORD"


echo "-------------------------------------------------------------"
echo " Set Permissions"
echo "-------------------------------------------------------------"

chmod 644 "$CERT_NAME".pem

echo "-------------------------------------------------------------"
echo " Copy certificate and key to client dist folder"
echo "-------------------------------------------------------------"

mkdir -p "$DEST"

mv "${CERT_NAME}".crt "$DEST"
mv "${CERT_NAME}".key "$DEST"
mv "${CERT_NAME}".pem "$DEST"
mv "${CERT_NAME}".pfx "$DEST"

# Cleanup only temporary files; be tolerant if they don't exist
rm -f "${CONFIG}"
rm -f "${CERT_NAME}".*