package Paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Acceptor struct {
	Lis          net.Listener
	ServerId     int         //服务器id
	MaxNumber    int         //接受者的承诺值
	AcceptNumber int         //已接受提案的最大提案编号
	AcceptValue  interface{} //已接受提案值
	Learners     []int
}

func (a *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {

	if args.Number > a.MaxNumber {
		a.MaxNumber = args.Number
		reply.Number = a.AcceptNumber
		reply.Value = a.AcceptValue
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {

	if args.Number >= a.MaxNumber {
		a.MaxNumber = args.Number
		a.AcceptValue = args.Value
		a.AcceptNumber = args.Number
		reply.Ok = true

		for _, lid := range a.Learners {
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = a.ServerId
				args.To = learner
				resp := new(MsgReply)
				ok := call(addr, "Learner.Learn", args, resp)
				if !ok {
					return
				}
			}(lid)
		}

	} else {
		reply.Ok = false
	}
	return nil
}

func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		ServerId: id,
		Learners: learners,
	}
	acceptor.server()
	return acceptor
}
func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(":%d", a.ServerId)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error : ", e)
	}
	a.Lis = l
	go func() {
		for {
			conn, err := a.Lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}
func (a *Acceptor) close() {
	a.Lis.Close()
}
