package Paxos

import (
	"fmt"
	"net/rpc"
)

type MsgArgs struct {
	Number int
	Value  interface{}
	From   int
	To     int
}

type MsgReply struct {
	Ok     bool
	Number int
	Value  interface{}
}

func call(srv string, name string, args interface{}, reply interface{}) bool {

	c, err := rpc.Dial("tcp", srv)
	if err != nil {
		return false
	}

	defer c.Close()
	err = c.Call(name, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}
