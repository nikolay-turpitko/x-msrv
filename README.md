# x-msrv
x-msrv is a sample bluerprint microserivece on Go programming language with
systemd and Docker Compose.

This project initially created as a job interview test solution.
Though, it's may be useful snippet of code for future use.

## Original Requirements

1. Microservice's configuration contains
    - message queue IP address and topic name to listen to;
    - N - max number of messages to process during one invocation;
    - M - period of invocations in seconds.
2. Every M seconds microserivece connects to topic and receives N messages at most.
3. Messages are JSON objects, they should be validated, parsed and stored into DB.
4. Queue impl - http://nsq.io.
5. DB impl - http://www.aerospike.com/.
6. Service runs as a daemon and correctly handles `SIGTERM`.
7. Any additional implementation decisions not listed here make as you feel suitable.
8. Document your implementation.

## Additional requirements

Let's fix some moving parts, 

1. Let's use Docker Compose to make environment reproducible.
2. Because of some implementation difficulties (Go scheduler and threading
   implementation), it's hard to correctly daemonize program using pure Go.
   There are abundance of wrappers and launchers for that task.
   Most of modern Linux distro use `systemd`, which is perfect for this job.
   For older distro it's possible to use `runit` or similar.
   Good option to run within Docker container (using `runit`):
   https://github.com/phusion/baseimage-docker.
   Let's use `systemd`, as more flexible, robust and tried out alternative.
3. Let's deploy binary file and accompanying configuration files using standard
   deployment formats (deb, rpm, ...). This way microserivece may be installed
   on minimal standard distro with any convenient tools, as well as in different
   container environments. __Note__: this step a bit complicates build process
   and can be viewed by someone as unnecessary hurdle. Well, it, probably, is.
   But it gives a bit more flexibility.
4. Let's use `glide` for vendoring dependencies to produce reproducible builds.
5. Let's build in the Docker container to produce reproducible builds.

## Prerequisites

I created and tested this solution on the following host environment:

- Ubuntu 16.04.2 LTS (4.4.0-64-generic x86_64)
- Docker version 1.12.6, build 78d1802
- docker-compose version 1.8.1, build 878cff1
- go version go1.8 linux/amd64
- glide version v0.12.3
- GNU Make 4.1
- GNU bash, version 4.3.46(1)-release (x86_64-pc-linux-gnu)


## High-level solution description

Microservice implemented as a `systemd` timer.

This way `systemd` manages service creation and termination, regular execution,
restart, priveleges, etc.

All required settings can be changed via ./`systemd`/x-msrv.service and
./`systemd`/x-msrv.timer files. Particularly, interval between execution can be
changed in ./`systemd`/x-msrv.timer file. Both files can be amended by administrator
after installation.

`systemd` can be configured to send desired termination signal. It sends `SIGTERM`
by default. It can be also configured to implement desired kill strategy.

If required, `systemd` can be used within Docker container. Though for many
simple applications I would not use it (Docker container, I mean).
I would start with bare CentOS or Debian with required dependencies and standard
tools - rpm / apt, systemd, rsyslog.

In small apps microseriveces may even have no need to know about each other, so
simple configuration files would be enough.  More complex orchestration may be
introduced later with app growth.  It will be possible to continue use systemd
within container and not very hard (from the service's implementation point of
view) to replace it with some alternative, if necessary. So, I want to
emphasize, that I used Docker container here mostly to create reproducible test
environment, not for deployment into production.

Use `make` to build and execute application.

Build script in `Makefile` uses Docker Compose to create Docker container, which
performs actual Go build and RPM packaging. After that Docker Compose used again
to prepare execution environment and start microservice within Docker container.

File `./docker-compose-build.yml` used to prepare Docker container for build
and `./docker-compose.yml` - to execute app. Notice, that environment variable
`PKG_TYPE` (rpm/deb) should be provided to Docker Compose.

Directory `./rpm-build` contains dockerfile and Makefile for build and RPM
packaging. They used to prepare container with required versions of go compiler
and glide (dependency management tool). Versions should be provided as build
arguments via `./docker-compose-build.yml`. Sources and output directories and
spec file passed to container via volumes, mapped to host directories. Actual
build commands are in the `./rpm-spec/x-msrv.spec` file, which is a file used
by RPM package manager. Spec file contains instructions to build Go sources and
to package binary files, documentation and configuration files into RPM archive.

Similarly, directory `./deb-build` contains dockerfile and required scripts to
build deb package in Docker container.

Directory `systemd` contains our service's unit and timer description files for
`systemd`.

Directory `./*-deploy` contains dockerfile to create microservice's Docker
container.  For that it uses `./rpm-deploy/rpms` or `./deb-deploy/debs` subdir,
which is mapped as a volume to previously described building and packaging
Docker container, which puts resulting rpm file into this directory.

This blueprint implementation installs simple man page for microservice (just
to illustrate how to do it). I included man installation into Docker image, so
that this feature could be tested. Use following command to check it from host:

    sudo PKG_TYPE=rpm docker-compose exec app man x-msrv

BTW, to check service status and logs use:

    sudo PKG_TYPE=rpm docker-compose exec app systemctl status x-msrv
    sudo PKG_TYPE=rpm docker-compose exec app journalctl -efx

TODO: config file
