# The testing task
The full description of the task could be found in the [Task](doc/New Orca project Test Task.pdf)

# Requiremets for the system
Installed:
- docker (with ability to run from the user)
- make

# Build
```sh
$ make build
```
<details><summary>Example</summary>

```sh
$ make build
docker build \
        --build-arg GO_VER=1.17.8 \
        --build-arg ALPINE_VER=3.15 \
        --build-arg WORKDIR=/go/src/github.com/vasily.chertkov/orca-task \
        -t vasily.chertkov/orca-task:1.0.0 -f /home/vchertkov/dev/orca/task/docker/Dockerfile .
Sending build context to Docker daemon  148.1MB
Step 1/13 : ARG GO_VER
Step 2/13 : ARG ALPINE_VER
Step 3/13 : FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder
 ---> 16b6bb19f174
Step 4/13 : LABEL stage=server-intermediate
 ---> Using cache
 ---> cd20a236d9d7
Step 5/13 : ARG WORKDIR
 ---> Using cache
 ---> 79477d488c92
Step 6/13 : WORKDIR ${WORKDIR}
 ---> Using cache
 ---> 146c987cdb4f
Step 7/13 : RUN apk --no-cache --update add     git
 ---> Using cache
 ---> 9c953058b45d
Step 8/13 : COPY ./ ./
 ---> 8779de2510a6
Step 9/13 : RUN CGO_ENABLED=0 GOOS=linux go build -v -mod=vendor -o /tmp/orca-task ./cmd/server/*.go
 ---> Running in 854371d2f880
github.com/modern-go/reflect2
golang.org/x/sys/unix
github.com/modern-go/concurrent
net
github.com/sirupsen/logrus
github.com/json-iterator/go
vendor/golang.org/x/net/http/httpproxy
net/textproto
crypto/x509
mime/multipart
vendor/golang.org/x/net/http/httpguts
crypto/tls
net/http/httptrace
net/http
command-line-arguments
Removing intermediate container 854371d2f880
 ---> 19fc26510bb5
Step 10/13 : FROM alpine:${ALPINE_VER} as base
 ---> c059bfaa849c
Step 11/13 : COPY --from=builder /tmp/orca-task /usr/local/bin/orca-task
 ---> 004f2eaa84a0
Step 12/13 : EXPOSE 8080
 ---> Running in 189ec111ac3d
Removing intermediate container 189ec111ac3d
 ---> fbecbec77421
Step 13/13 : ENTRYPOINT ["/usr/local/bin/orca-task"]
 ---> Running in 989a7f989cb9
Removing intermediate container 989a7f989cb9
 ---> 787c843b9e82
Successfully built 787c843b9e82
Successfully tagged vasily.chertkov/orca-task:1.0.0
```
</details>

# Run
The service requires the json file with the data to be passed.
It can be done by setting the `INPUT_PATH` env var.
```sh
$ INPUT_PATH=cmd/server/fixtures/input-1000000.json make run
```
<details><summary>Example</summary>

```sh
$ INPUT_PATH=cmd/server/fixtures/input-1000000.json make run
docker build \
        --build-arg GO_VER=1.17.8 \
        --build-arg ALPINE_VER=3.15 \
        --build-arg WORKDIR=/go/src/github.com/vasily.chertkov/orca-task \
        -t vasily.chertkov/orca-task:1.0.0 -f /home/vchertkov/dev/orca/task/docker/Dockerfile .
Sending build context to Docker daemon  148.1MB
Step 1/13 : ARG GO_VER
Step 2/13 : ARG ALPINE_VER
Step 3/13 : FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder
 ---> 16b6bb19f174
Step 4/13 : LABEL stage=server-intermediate
 ---> Using cache
 ---> cd20a236d9d7
Step 5/13 : ARG WORKDIR
 ---> Using cache
 ---> 79477d488c92
Step 6/13 : WORKDIR ${WORKDIR}
 ---> Using cache
 ---> 146c987cdb4f
Step 7/13 : RUN apk --no-cache --update add     git
 ---> Using cache
 ---> 9c953058b45d
Step 8/13 : COPY ./ ./
 ---> Using cache
 ---> 8779de2510a6
Step 9/13 : RUN CGO_ENABLED=0 GOOS=linux go build -v -mod=vendor -o /tmp/orca-task ./cmd/server/*.go
 ---> Using cache
 ---> 19fc26510bb5
Step 10/13 : FROM alpine:${ALPINE_VER} as base
 ---> c059bfaa849c
Step 11/13 : COPY --from=builder /tmp/orca-task /usr/local/bin/orca-task
 ---> Using cache
 ---> 004f2eaa84a0
Step 12/13 : EXPOSE 8080
 ---> Using cache
 ---> fbecbec77421
Step 13/13 : ENTRYPOINT ["/usr/local/bin/orca-task"]
 ---> Using cache
 ---> 787c843b9e82
Successfully built 787c843b9e82
Successfully tagged vasily.chertkov/orca-task:1.0.0
docker run \
        --rm \
        -u `id -u`:`id -g` \
        -v /home/vchertkov/dev/orca/task/cmd/server/fixtures/input-1000000.json:/input.json \
        --name orca-task \
        -p 80:8080 \
        vasily.chertkov/orca-task:1.0.0
Alloc = 316 MiB TotalAlloc = 350 MiB    Sys = 345 MiB   NumGC = 2
time="2022-03-13T20:13:57Z" level=info msg="Loading data time: 718.554104ms"
time="2022-03-13T20:13:57Z" level=info msg="Processing FW Rules time: 312.4µs"
Alloc = 371 MiB TotalAlloc = 405 MiB    Sys = 401 MiB   NumGC = 2
time="2022-03-13T20:13:58Z" level=info msg="Processing VMs time: 675.079165ms"
Alloc = 656 MiB TotalAlloc = 754 MiB    Sys = 699 MiB   NumGC = 3
Alloc = 806 MiB TotalAlloc = 1637 MiB   Sys = 944 MiB   NumGC = 5
time="2022-03-13T20:14:00Z" level=info msg="Processing Tag Sets time: 2.418800015s"
time="2022-03-13T20:14:00Z" level=info msg="Preprocessing time: 3.817788574s"
time="2022-03-13T20:14:00Z" level=info msg="Listening on :8080..."
```
</details>

# API
The server runs on `8080` port in docker dontainer which is mapped to the `80` port on the host.

It's accessible by `http://localhost/api/v1/` URL.

## Obtaining the info about the potential attackers:
- http://localhost/api/v1/attach?vm_id={vm_id}
```sh
$ curl "http://localhost/api/v1/attack?vm_id=vm-8c849a1e"
```

## Getting statistics
- http://localhost/api/v1/stats
```sh
curl "http://localhost/api/v1/stats"
```

# Tests
The tests can be run with the following command:
```sh
$ make test
```

# Generation of the larger input
To generate input with 1000000 VMs and 1000 FW rules and save the result to `cmd/server/fixtures/input-1000000.json`
```sh
$ python cmd/server/fixtures/generate.py -vmc 1000000 -fwc 1000 cmd/server/fixtures/input-1000000.json
```
This generated file can be used later as a input for service:
```sh
$ INPUT_PATH=cmd/server/fixtures/input-1000000.json make run
```