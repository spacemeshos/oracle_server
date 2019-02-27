# oracle_server
a no-variance mock oracle server 

## compile protobufs

`make proto`


## build 

`make build`

## change ports

for now edit `main.go`

## WorldID 

In order to use to oracle for different tests at the same time, every request is required to pass
a unique `WorldID` basically an `int` that is used as an exectuion id. this id isolates all oracle functions like register/deregister and eligiblity to that id.

Tests that collapsed without unregistering might leave hanging worlds, hence its important to use a unique id everytime.
