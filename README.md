# Go8

## Delve Issue

Run `./build_and_run_test_container.sh {DOCKERFILE}` to build the container and execute the expect script.

## Output

Arch test container

```bash
❯ sudo ./build_and_run_test_container.sh arch.dockerfile
DEPRECATED: The legacy builder is deprecated and will be removed in a future release.
            Install the buildx component to build images with BuildKit:
            https://docs.docker.com/go/buildx/

Sending build context to Docker daemon  5.641MB
Step 1/8 : FROM archlinux:latest
 ---> 69f38d8f6347
Step 2/8 : RUN pacman -Syu --noconfirm   && pacman -S go delve expect --noconfirm
 ---> Using cache
 ---> 2eeaf4d35198
Step 3/8 : COPY *.go /app/
 ---> Using cache
 ---> a0f258aa4efd
Step 4/8 : COPY go.* /app/
 ---> Using cache
 ---> d21fc3ef5b7b
Step 5/8 : COPY dlv_test.exp /app/
 ---> Using cache
 ---> 9e363a90a1b4
Step 6/8 : WORKDIR /app
 ---> Using cache
 ---> 8698d7d19dde
Step 7/8 : RUN chmod +x /app/dlv_test.exp
 ---> Using cache
 ---> da285ad443cd
Step 8/8 : CMD ["/app/dlv_test.exp"]
 ---> Using cache
 ---> 8fce3a22dbec
Successfully built 8fce3a22dbec
Successfully tagged temp-dlv-image:latest
spawn dlv test ./
go: downloading github.com/stretchr/testify v1.8.4
go: downloading github.com/davecgh/go-spew v1.1.1
go: downloading gopkg.in/yaml.v3 v3.0.1
go: downloading github.com/pmezard/go-difflib v1.0.0
Type 'help' for list of commands.
(dlv) break ./emulate_cycle.go:194
Breakpoint 1 set at 0x63691d for github.com/postcert/go8.(*Chip8).addRegistersWithCarry() ./emulate_cycle.go:194
(dlv) continue
> github.com/postcert/go8.(*Chip8).addRegistersWithCarry() ./emulate_cycle.go:194 (hits goroutine(35):1 total:1) (PC: 0x63691d)
   189:         chip.V[registerX>>8] ^= chip.V[registerY>>4]
   190: }
   191:
   192: func (chip *Chip8) addRegistersWithCarry(opcode uint16, value uint16) {
   193:         // TODO: 255 + 1 = 0, not 256 becuase of type. Failing overflow test below
=> 194:         temp := uint16(chip.V[opcode>>8]) + uint16(chip.V[opcode>>4])
   195:         if temp > 255 {
   196:                 chip.V[0xF] = 1
   197:         } else {
   198:                 chip.V[0xF] = 0
   199:         }
(dlv) stack -full
Sending output to pager...
0  0x000000000063691d in github.com/postcert/go8.(*Chip8).addRegistersWithCarry
   at ./emulate_cycle.go:194
       chip = (unreadable empty OP stack)
       opcode = (unreadable empty OP stack)
       value = (unreadable empty OP stack)

1  0x0000000000639b94 in github.com/postcert/go8.TestAddRegistersWithCarry_Success
   at ./emulate_cycle_test.go:254
       t = (*testing.T)(0xc0001c8d00)
       chip = ("*github.com/postcert/go8.Chip8")(0xc0001d2600)

2  0x0000000000522453 in testing.tRunner
   at /usr/lib/go/src/testing/testing.go:1595
       t = (*testing.T)(0xc0001c8d00)
       fn = github.com/postcert/go8.TestAddRegistersWithCarry_Success

3  0x0000000000524093 in testing.(*T).Run.func1
   at /usr/lib/go/src/testing/testing.go:1648

4  0x0000000000474d01 in runtime.goexit
   at /usr/lib/go/src/runtime/asm_amd64.s:1650

(dlv) %

```

Ubuntu test container

```bash
❯ sudo ./build_and_run_test_container.sh ubuntu.dockerfile
DEPRECATED: The legacy builder is deprecated and will be removed in a future release.
            Install the buildx component to build images with BuildKit:
            https://docs.docker.com/go/buildx/

Sending build context to Docker daemon  5.644MB
Step 1/17 : FROM ubuntu:latest
 ---> e34e831650c1
Step 2/17 : ARG DEBIAN_FRONTEND=noninteractive
 ---> Using cache
 ---> 1c8a728b4b2c
Step 3/17 : ENV TZ=America/Los_Angeles
 ---> Using cache
 ---> c90973460d49
Step 4/17 : RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime   && echo $TZ > /etc/timezone
 ---> Using cache
 ---> ed2637d4ee55
Step 5/17 : RUN apt-get update   && apt-get install -y wget git expect   && rm -rf /var/lib/apt/lists/*
 ---> Using cache
 ---> 6af791e25a89
Step 6/17 : ENV GO_VERSION 1.21.6
 ---> Using cache
 ---> 8d8303392aca
Step 7/17 : ENV GOPATH /root/go
 ---> Using cache
 ---> e4a29dec9327
Step 8/17 : ENV PATH $GOPATH/bin:$PATH
 ---> Using cache
 ---> c348045c745b
Step 9/17 : RUN wget --progress=dot:mega https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz   && tar -C /usr/local -xzvf /tmp/go.tar.gz --strip-components 1   && rm /tmp/go.tar.gz
 ---> Using cache
 ---> 89198f2b8b93
Step 10/17 : RUN chmod +x /usr/local/bin/go
 ---> Using cache
 ---> dc45f272c818
Step 11/17 : RUN go install github.com/go-delve/delve/cmd/dlv@latest
 ---> Using cache
 ---> 6cae4a8c1a92
Step 12/17 : COPY *.go /app/
 ---> Using cache
 ---> a1cab535f46f
Step 13/17 : COPY go.* /app/
 ---> Using cache
 ---> 3c2fbc98bebe
Step 14/17 : COPY dlv_test.exp /app/
 ---> Using cache
 ---> 188fe9278791
Step 15/17 : WORKDIR /app
 ---> Using cache
 ---> d54ca6b50cdd
Step 16/17 : RUN chmod +x /app/dlv_test.exp
 ---> Using cache
 ---> 981ab12c18fe
Step 17/17 : CMD ["/app/dlv_test.exp"]
 ---> Using cache
 ---> 425be3d7ad58
Successfully built 425be3d7ad58
Successfully tagged temp-dlv-image:latest
spawn dlv test ./
go: downloading github.com/stretchr/testify v1.8.4
go: downloading github.com/davecgh/go-spew v1.1.1
go: downloading github.com/pmezard/go-difflib v1.0.0
Type 'help' for list of commands.
(dlv) break ./emulate_cycle.go:194
Breakpoint 1 set at 0x63691d for github.com/postcert/go8.(*Chip8).addRegistersWithCarry() ./emulate_cycle.go:194
(dlv) continue
> github.com/postcert/go8.(*Chip8).addRegistersWithCarry() ./emulate_cycle.go:194 (hits goroutine(35):1 total:1) (PC: 0x63691d)
   189:         chip.V[registerX>>8] ^= chip.V[registerY>>4]
   190: }
   191:
   192: func (chip *Chip8) addRegistersWithCarry(opcode uint16, value uint16) {
   193:         // TODO: 255 + 1 = 0, not 256 becuase of type. Failing overflow test below
=> 194:         temp := uint16(chip.V[opcode>>8]) + uint16(chip.V[opcode>>4])
   195:         if temp > 255 {
   196:                 chip.V[0xF] = 1
   197:         } else {
   198:                 chip.V[0xF] = 0
   199:         }
(dlv) stack -full
Sending output to pager...
0  0x000000000063691d in github.com/postcert/go8.(*Chip8).addRegistersWithCarry
   at ./emulate_cycle.go:194
       chip = (unreadable empty OP stack)
       opcode = (unreadable empty OP stack)
       value = (unreadable empty OP stack)

1  0x0000000000639b94 in github.com/postcert/go8.TestAddRegistersWithCarry_Success
   at ./emulate_cycle_test.go:254
       t = (*testing.T)(0xc0001bf1e0)
       chip = ("*github.com/postcert/go8.Chip8")(0xc0001d2600)

2  0x0000000000522453 in testing.tRunner
   at /usr/local/src/testing/testing.go:1595
       t = (*testing.T)(0xc0001bf1e0)
       fn = github.com/postcert/go8.TestAddRegistersWithCarry_Success

3  0x0000000000524093 in testing.(*T).Run.func1
   at /usr/local/src/testing/testing.go:1648

4  0x0000000000474d01 in runtime.goexit
   at /usr/local/src/runtime/asm_amd64.s:1650

(dlv) %
```
