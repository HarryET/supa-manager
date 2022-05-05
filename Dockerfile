FROM rust:1.60.0 as builder

RUN USER=root cargo new --bin supa-manager
WORKDIR ./supa-manager
COPY ./Cargo.toml ./Cargo.toml
RUN cargo build --release
RUN rm src/*.rs

ADD . ./

RUN rm ./target/release/deps/supa-manager*
RUN cargo build --release

FROM alpine

ENV ROCKET_PROFILE="release" \
    ROCKET_ADDRESS=0.0.0.0 \
    ROCKET_PORT=80

COPY --from=builder /supa-manager/target/release/supa-manager /usr/bin/supa-manager

EXPOSE 80

WORKDIR /usr/bin

CMD ["./rust-docker-web"]