package Paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Learner struct {
	Lis         net.Listener
	Id          int             //学习者id
	AcceptedMsg map[int]MsgArgs //记录接受者已经接受的提案  【接受者id：请求消息】
}

func (l *Learner) Learn(args *MsgArgs, reply *MsgReply) error {

	a := l.AcceptedMsg[args.From]

	if a.Number < args.Number {
		l.AcceptedMsg[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (l *Learner) Chosen() interface{} {
	acceptCounts := make(map[int]int)
	acceptMsg := make(map[int]MsgArgs)

	for _, accepted := range l.AcceptedMsg {
		if accepted.Number != 0 {
			acceptCounts[accepted.Number]++
			acceptMsg[accepted.Number] = accepted
		}
	}

	for n, count := range acceptCounts {
		if count >= l.majority() {
			return acceptMsg[n].Value
		}
	}
	return nil

}

func (l *Learner) majority() int {
	return len(l.AcceptedMsg)/2 + 1
}

func newLearner(id int, acceptorIds []int) *Learner {
	Learner := &Learner{
		Id:          id,
		AcceptedMsg: make(map[int]MsgArgs),
	}
	for _, aid := range acceptorIds {
		Learner.AcceptedMsg[aid] = MsgArgs{
			Number: 0,
			Value:  nil,
		}
	}
	Learner.server(id)
	return Learner
}
func (l *Learner) server(id int) {
	rpcs := rpc.NewServer()
	rpcs.Register(l)
	addr := fmt.Sprintf(":%d", id)
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error : ", e)
	}
	l.Lis = lis
	go func() {
		for {
			conn, err := l.Lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}
func (l *Learner) close() {
	l.Lis.Close()
}
