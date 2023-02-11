# docker build -t leehom/v2raya-guard:latest .

# -- Stage 1 -- #
# Compile the app.
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN apk add g++ make git
RUN make build

# -- Stage 2 -- #
# Create the final environment with the compiled binary.
FROM alpine:3.17.2
MAINTAINER leehom Chen <clh021@gmail.com>
LABEL maintainer="leehom Chen <clh021@gmail.com>"
COPY --from=builder /app/bin/v2raya-guard /usr/local/bin/
CMD ["v2raya-guard"]
