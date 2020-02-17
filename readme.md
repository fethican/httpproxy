## Simple HTTP Proxy Server

SERVER_PORT=8080

PROXY_PROTO=https
PROXY_TO=s3.amazonaws.com

BUCKET_NAME=""

Exposes /health and /metrics endpoins. Forwards everything else to PROXY_TO.
