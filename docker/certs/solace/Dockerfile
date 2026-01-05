FROM alpine:3.21.2

# Install OpenSSL
RUN apk update && \
    apk add --no-cache openssl

# Create a non-root user
RUN addgroup -r solaceuser && adduser -r -g solaceuser solaceuser

# switch to the non-root user
USER solaceuser