FROM scratch
ARG IMAGE_SOURCE=""
LABEL org.opencontainers.image.source ${IMAGE_SOURCE}
WORKDIR /app
COPY dist/cc-auth-service cc-auth-service
ENTRYPOINT [ "/cc-auth-service" ]