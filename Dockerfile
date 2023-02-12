# docker build -t leehom/v2raya-guard:latest .
# docker run -d --name v2raya-guard -v /etc/localtime:/etc/localtime:ro -e username="root" -e password="l" -e serverbaseurl="http://192.168.1.2:2017/api/" -e cronExp="0 3,12,21 * * *"  leehom/v2raya-guard

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
