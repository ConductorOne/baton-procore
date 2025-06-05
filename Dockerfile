FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-procore"]
COPY baton-procore /