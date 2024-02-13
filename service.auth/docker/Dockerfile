FROM scratch
ARG IMAGE_SOURCE=""
LABEL org.opencontainers.image.source ${IMAGE_SOURCE}
LABEL org.opencontainers.image.description "Authentication service for cc-microsvcs."
WORKDIR /app
COPY dist/cc-auth-service cc-auth-service
ENTRYPOINT [ "/cc-auth-service" ]