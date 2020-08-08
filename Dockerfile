FROM golang:1.14 AS builder

COPY . /workspace/
WORKDIR /workspace
RUN make bin/server

FROM scratch

COPY --from=builder /workspace/bin/server ./server

ENTRYPOINT ["/server"]

