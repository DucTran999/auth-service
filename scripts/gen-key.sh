#!/bin/bash

# Output directory
KEY_DIR="./keys"
PRIVATE_KEY="$KEY_DIR/private.pem"
PUBLIC_KEY="$KEY_DIR/public.pem"

# Create the directory if it doesn't exist
mkdir -p "$KEY_DIR"

# Generate RSA private key (2048 bits is secure, 4096 is stronger)
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY" -pkeyopt rsa_keygen_bits:2048

# Extract public key from private key
openssl rsa -pubout -in "$PRIVATE_KEY" -out "$PUBLIC_KEY"

echo "✅ RSA key pair generated:"
echo "🔒 Private Key: $PRIVATE_KEY"
echo "🔓 Public Key:  $PUBLIC_KEY"
