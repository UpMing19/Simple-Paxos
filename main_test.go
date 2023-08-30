package Paxos

import (
	"testing"
)

func start(acceptorIds []int, learnerIds []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	learners := make([]*Learner, 0)

	for _, aid := range acceptorIds {
		a := newAcceptor(aid, learnerIds)
		acceptors = append(acceptors, a)
	}
	for _, lid := range learnerIds {
		l := newLearner(lid, acceptorIds)
		learners = append(learners, l)
	}

	return acceptors, learners
}
func cleanup(acceptors []*Acceptor, learners []*Learner) {
	for _, a := range acceptors {
		a.close()
	}
	for _, l := range learners {
		l.close()
	}
}
func TestSingleProposer(t *testing.T) {
	// 1001 - 1002 - 1003是接受者id
	acceptorIds := []int{1001, 1002, 1003}
	//2001是学习者id
	learnerIds := []int{2001}

	acceptors, learns := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learns)

	p := &Proposer{
		ServerId:  1,
		Acceptors: acceptorIds,
	}
	value := p.propose("hello")

	if value != "hello" {
		t.Errorf("value = %s, excepted %s ", value, "hello")
	}

	learnerValue := learns[0].Chosen()

	if learnerValue != value {
		t.Errorf("learnerValue = %s, excepted %s ", learnerValue, "hello")
	}

}

func TestDoubleProposer(t *testing.T) {
	// 1001 - 1002 - 1003是接受者id
	acceptorIds := []int{1001, 1002, 1003}
	//2001是学习者id
	learnerIds := []int{2001}

	acceptors, learns := start(acceptorIds, learnerIds)

	defer cleanup(acceptors, learns)

	p := &Proposer{
		ServerId:  1,
		Acceptors: acceptorIds,
	}

	p2 := &Proposer{
		ServerId:  2,
		Acceptors: acceptorIds,
	}
	value1 := p.propose("hello")
	value2 := p2.propose("world")

	if value2 != value1 {
		t.Errorf("value1 : %s  value2 : %s ", value1, value2)
	}
	learnerValue := learns[0].Chosen()
	if learnerValue != value1 {
		t.Errorf("value1 : %s  value2 : %s  LearnValuse:%s", value1, value2, learnerValue)
	}

}
