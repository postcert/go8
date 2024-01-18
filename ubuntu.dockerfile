# Use the latest Ubuntu base image
FROM ubuntu:latest

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=America/Los_Angeles
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone

# Update the system and install Go, Delve, and Expect
RUN apt-get update \
  && apt-get install -y wget git expect \
  && rm -rf /var/lib/apt/lists/*

# Define Go version and installation path
ENV GO_VERSION 1.21.6
ENV GOPATH /root/go
ENV PATH $GOPATH/bin:$PATH

RUN wget --progress=dot:mega https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz \
  && tar -C /usr/local -xzvf /tmp/go.tar.gz --strip-components 1 \
  && rm /tmp/go.tar.gz

# Make Go binary executable
RUN chmod +x /usr/local/bin/go

# Install Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy all Go files to /app directory in the container
COPY *.go /app/
COPY go.* /app/

# Copy the scripts into the container
COPY echo_ver_and_run_test.sh /app/
COPY dlv_test.exp /app/

# Set the working directory to /app
WORKDIR /app

# Make the scripts executable
RUN chmod +x /app/echo_ver_and_run_test.sh && chmod +x /app/dlv_test.exp

# Run the Expect script when the container starts
CMD ["/app/echo_ver_and_run_test.sh"]
