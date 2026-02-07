# Corrected Dockerfile: Ubuntu with systemd and dbus
# Location: docker/ubuntu/Dockerfile

FROM ubuntu:22.04

# Non-interactive front-end and marker for systemd-in-docker
ENV container=docker \
    DEBIAN_FRONTEND=noninteractive

# Install systemd and dbus non-interactively
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    systemd \
    dbus \
    dbus-user-session \
    dbus-x11 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# systemd needs access to cgroups
VOLUME ["/sys/fs/cgroup"]
STOPSIGNAL SIGRTMIN+3

# Start systemd as PID 1
CMD ["/lib/systemd/systemd"]
