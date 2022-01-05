FROM starudream/golang AS builder

WORKDIR /build

COPY . .

RUN make build && make upx

FROM starudream/alpine-glibc:latest

WORKDIR /

COPY --from=builder /build/bin/app /app

CMD /app
