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

Sending build context to Docker daemon  5.659MB
Step 1/9 : FROM archlinux:latest
 ---> 69f38d8f6347
Step 2/9 : RUN pacman -Syu --noconfirm   && pacman -S go delve expect --noconfirm
 ---> Using cache
 ---> 2eeaf4d35198
Step 3/9 : COPY *.go /app/
 ---> Using cache
 ---> a0f258aa4efd
Step 4/9 : COPY go.* /app/
 ---> Using cache
 ---> d21fc3ef5b7b
Step 5/9 : COPY echo_ver_and_run_test.sh /app/
 ---> Using cache
 ---> b3fff430a740
Step 6/9 : COPY dlv_test.exp /app/
 ---> Using cache
 ---> 843615636709
Step 7/9 : WORKDIR /app
 ---> Using cache
 ---> ca9035967af3
Step 8/9 : RUN chmod +x /app/echo_ver_and_run_test.sh && chmod +x /app/dlv_test.exp
 ---> Using cache
 ---> 76b4f0d8930b
Step 9/9 : CMD ["/app/echo_ver_and_run_test.sh"]
 ---> Using cache
 ---> 80a210a78551
Successfully built 80a210a78551
Successfully tagged temp-dlv-image:latest
Go info:
go version go1.21.6 linux/amd64
/app/echo_ver_and_run_test.sh: line 5: which: command not found


Delve info:
Delve Debugger
Version: 1.22.0
Build: $Id: 61ecdbbe1b574f0dd7d7bad8b6a5d564cce981e9 $
/app/echo_ver_and_run_test.sh: line 10: which: command not found


spawn dlv test ./
go: downloading github.com/stretchr/testify v1.8.4
go: downloading gopkg.in/yaml.v3 v3.0.1
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
       t = (*testing.T)(0xc0001c0ea0)
       chip = ("*github.com/postcert/go8.Chip8")(0xc0001d0600)

2  0x0000000000522453 in testing.tRunner
   at /usr/lib/go/src/testing/testing.go:1595
       t = (*testing.T)(0xc0001c0ea0)
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

Sending build context to Docker daemon  5.659MB
Step 1/18 : FROM ubuntu:latest
 ---> e34e831650c1
Step 2/18 : ARG DEBIAN_FRONTEND=noninteractive
 ---> Using cache
 ---> 1c8a728b4b2c
Step 3/18 : ENV TZ=America/Los_Angeles
 ---> Using cache
 ---> c90973460d49
Step 4/18 : RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime   && echo $TZ > /etc/timezone
 ---> Using cache
 ---> ed2637d4ee55
Step 5/18 : RUN apt-get update   && apt-get install -y wget git expect   && rm -rf /var/lib/apt/lists/*
 ---> Using cache
 ---> 6af791e25a89
Step 6/18 : ENV GO_VERSION 1.21.6
 ---> Using cache
 ---> 8d8303392aca
Step 7/18 : ENV GOPATH /root/go
 ---> Using cache
 ---> e4a29dec9327
Step 8/18 : ENV PATH $GOPATH/bin:$PATH
 ---> Using cache
 ---> c348045c745b
Step 9/18 : RUN wget --progress=dot:mega https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz   && tar -C /usr/local -xzvf /tmp/go.tar.gz --strip-components 1   && rm /tmp/go.tar.gz
 ---> Using cache
 ---> 89198f2b8b93
Step 10/18 : RUN chmod +x /usr/local/bin/go
 ---> Using cache
 ---> dc45f272c818
Step 11/18 : RUN go install github.com/go-delve/delve/cmd/dlv@latest
 ---> Using cache
 ---> 6cae4a8c1a92
Step 12/18 : COPY *.go /app/
 ---> Using cache
 ---> a1cab535f46f
Step 13/18 : COPY go.* /app/
 ---> Using cache
 ---> 3c2fbc98bebe
Step 14/18 : COPY echo_ver_and_run_test.sh /app/
 ---> 305a0bf4dea9
Step 15/18 : COPY dlv_test.exp /app/
 ---> 020cedefc061
Step 16/18 : WORKDIR /app
 ---> Running in fadb32b99d20
Removing intermediate container fadb32b99d20
 ---> cb83ee948a8e
Step 17/18 : RUN chmod +x /app/echo_ver_and_run_test.sh && chmod +x /app/dlv_test.exp
 ---> Running in ffff39817345
Removing intermediate container ffff39817345
 ---> f56c798d073e
Step 18/18 : CMD ["/app/echo_ver_and_run_test.sh"]
 ---> Running in 45f1e9558aa2
Removing intermediate container 45f1e9558aa2
 ---> bae1b8c4bdef
Successfully built bae1b8c4bdef
Successfully tagged temp-dlv-image:latest
Go info:
go version go1.21.6 linux/amd64
/usr/local/bin/go

Delve info:
Delve Debugger
Version: 1.22.0
Build: $Id: 61ecdbbe1b574f0dd7d7bad8b6a5d564cce981e9 $
/root/go/bin/dlv

spawn dlv test ./
go: downloading github.com/stretchr/testify v1.8.4
go: downloading github.com/pmezard/go-difflib v1.0.0
go: downloading github.com/davecgh/go-spew v1.1.1
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
       t = (*testing.T)(0xc0001c1040)
       chip = ("*github.com/postcert/go8.Chip8")(0xc0001d0600)

2  0x0000000000522453 in testing.tRunner
   at /usr/local/src/testing/testing.go:1595
       t = (*testing.T)(0xc0001c1040)
       fn = github.com/postcert/go8.TestAddRegistersWithCarry_Success

3  0x0000000000524093 in testing.(*T).Run.func1
   at /usr/local/src/testing/testing.go:1648

4  0x0000000000474d01 in runtime.goexit
   at /usr/local/src/runtime/asm_amd64.s:1650

(dlv) %
```
