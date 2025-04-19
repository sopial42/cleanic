FROM gcr.io/distroless/base-debian11

ARG TARGETARCH

ARG APPNAME

WORKDIR /

COPY build/${APPNAME}.${TARGETARCH} /server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]
