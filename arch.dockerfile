# Use the latest Arch Linux base image
FROM archlinux:latest

# Update the system and install Go, Delve, and Expect
RUN pacman -Syu --noconfirm \
  && pacman -S go delve expect --noconfirm 

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

