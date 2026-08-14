#!/bin/bash

# Create certs directory
mkdir -p certs
cd certs

echo "Generating CA..."
openssl req -x509 -nodes -newkey rsa:2048 -days 365 -keyout ca.key -out ca.crt -subj "/CN=Phalanx Test CA"

echo "Generating Server Cert..."
openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=server"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile <(printf "subjectAltName=DNS:service1.local,DNS:service2.local,DNS:localhost")

echo "Generating Client Cert..."
openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr -subj "/CN=client"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365

echo "Certificates generated in certs/ directory."
