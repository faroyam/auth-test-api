mkdir keys 2>/dev/null
rm keys/* 2>/dev/null
openssl genrsa -out keys/app.rsa 1024
openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub

openssl ecparam -genkey -name secp384r1 -out keys/server.key
openssl req -new -x509 -sha256 -key keys/server.key -out keys/server.crt -days 3650