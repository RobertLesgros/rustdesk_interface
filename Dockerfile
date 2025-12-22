# Build stage (optional, for building from source)
FROM alpine:3.19 AS base

# Security: Use specific version, not latest
# Security: Install only necessary packages
RUN apk add --no-cache tzdata ca-certificates

# Production stage
FROM base AS production

ARG BUILDARCH
WORKDIR /app

# Security: Create non-root user for running the application
RUN addgroup -g 1000 rustdesk && \
    adduser -u 1000 -G rustdesk -h /app -D rustdesk

# Copy the pre-built binary
COPY ./${BUILDARCH}/release /app/

# Security: Create necessary directories with proper permissions
RUN mkdir -p /app/data /app/runtime /app/conf && \
    chown -R rustdesk:rustdesk /app

# Security: Set restrictive permissions on binary
RUN chmod 550 /app/apimain

# Set up volume for persistent data
VOLUME /app/data

# Expose API port
EXPOSE 21114

# Security: Run as non-root user
USER rustdesk

# Security: Set read-only filesystem hint (enforced by docker run --read-only)
# Healthcheck for container orchestration
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:21114/api/health || exit 1

CMD ["./apimain"]
