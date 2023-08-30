package Paxos

import "fmt"

type Proposer struct {
	ServerId  int   //服务器id
	Round     int   //当前提议者已知最大轮次
	Number    int   //提案编号  （轮次，服务器id）
	Acceptors []int //接受者列表
}

func (p *Proposer) propose(v interface{}) interface{} {

	p.Round++
	p.Number = p.proposalNumber()

	//第一阶段

	prepareCount := 0
	maxNumber := 0

	for _, aid := range p.Acceptors {
		args := MsgArgs{
			Number: p.Number,
			From:   p.ServerId,
			To:     aid,
		}
		reply := new(MsgReply)
		err := call(fmt.Sprintf("127.0.0.1:%d", aid), "Acceptor.Prepare", args, reply)
		//fmt.Println("回复：", reply)
		if !err {
			continue
		}
		if reply.Ok {
			prepareCount++
			if reply.Number > maxNumber {
				maxNumber = reply.Number
				v = reply.Value
			}
		}
		if prepareCount == p.majority() {
			break
		}
	}

	//第二阶段

	acceptCount := 0
	if prepareCount >= p.majority() {

		for _, aid := range p.Acceptors {
			args := MsgArgs{
				Number: p.Number,
				Value:  v,
				From:   p.ServerId,
				To:     aid,
			}
			reply := new(MsgReply)
			err := call(fmt.Sprintf("127.0.0.1:%d", aid), "Acceptor.Accept", args, reply)
			if !err {
				continue
			}
			if reply.Ok {
				acceptCount++
			}
		}
	}

	if acceptCount >= p.majority() {
		return v
	}
	return nil
}

func (p *Proposer) majority() int {
	return len(p.Acceptors)/2 + 1
}
func (p *Proposer) proposalNumber() int {
	return p.Round<<16 | p.ServerId
}
