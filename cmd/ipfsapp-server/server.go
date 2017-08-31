package main

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"github.com/mackzhong/ipfsapp/labrpc"
	"github.com/mackzhong/ipfsapp/sraft"
	"log"
	"math/big"
	"sync"
	"time"
)

const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Kind  string //"Put" or "Append" "Get"
	Key   string
	Value string
	ID    int64
	ReqID int
}

type RaftKV struct {
	mu      sync.Mutex
	me      int
	rf      *sraft.Raft
	applyCh chan sraft.ApplyMsg

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	db     map[string]string
	ack    map[int64]int
	result map[int]chan Op
}

func (kv *RaftKV) AppendEntryToLog(entry Op) bool {
	index, _, isLeader := kv.rf.Start(entry)
	if !isLeader {
		return false
	}

	kv.mu.Lock()
	ch, ok := kv.result[index]
	if !ok {
		ch = make(chan Op, 1)
		kv.result[index] = ch
	}
	kv.mu.Unlock()
	select {
	case op := <-ch:
		return op == entry
	case <-time.After(1000 * time.Millisecond):
		//log.Printf("timeout\n")
		return false
	}
}

func (kv *RaftKV) CheckDup(id int64, reqid int) bool {
	//kv.mu.Lock()
	//defer kv.mu.Unlock()
	v, ok := kv.ack[id]
	if ok {
		return v >= reqid
	}
	return false
}

func (kv *RaftKV) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	entry := Op{Kind: "Get", Key: args.Key, ID: args.Id, ReqID: args.ReqID}

	ok := kv.AppendEntryToLog(entry)
	if !ok {
		reply.WrongLeader = true
	} else {
		reply.WrongLeader = false

		reply.Err = OK
		kv.mu.Lock()
		reply.Value = kv.db[args.Key]
		kv.ack[args.Id] = args.ReqID
		//log.Printf("%d get:%v value:%s\n",kv.me,entry,reply.Value)
		kv.mu.Unlock()
	}
}

func (kv *RaftKV) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	entry := Op{Kind: args.Op, Key: args.Key, Value: args.Value, ID: args.Id, ReqID: args.ReqID}
	ok := kv.AppendEntryToLog(entry)
	if !ok {
		reply.WrongLeader = true
	} else {
		reply.WrongLeader = false
		reply.Err = OK
	}
}

func (kv *RaftKV) Apply(args Op) {
	switch args.Kind {
	case "Put":
		kv.db[args.Key] = args.Value
	case "Append":
		kv.db[args.Key] += args.Value
	}
	kv.ack[args.ID] = args.ReqID
}

//
// the tester calls Kill() when a RaftKV instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (kv *RaftKV) Kill() {
	kv.rf.Kill()
	// Your code here, if desired.
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots with persister.SaveSnapshot(),
// and Raft should save its state (including log) with persister.SaveRaftState().
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *sraft.Persister, maxraftstate int) *RaftKV {
	// call gob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	gob.Register(Op{})

	kv := new(RaftKV)
	kv.me = me
	kv.maxraftstate = maxraftstate

	// Your initialization code here.

	kv.applyCh = make(chan sraft.ApplyMsg)
	kv.rf = sraft.Make(servers, me, persister, kv.applyCh)

	kv.db = make(map[string]string)
	kv.ack = make(map[int64]int)
	kv.result = make(map[int]chan Op)

	go func() {
		for {
			msg := <-kv.applyCh
			if msg.UseSnapshot {
				var LastIncludedIndex int
				var LastIncludedTerm int

				r := bytes.NewBuffer(msg.Snapshot)
				d := gob.NewDecoder(r)

				kv.mu.Lock()
				d.Decode(&LastIncludedIndex)
				d.Decode(&LastIncludedTerm)
				kv.db = make(map[string]string)
				kv.ack = make(map[int64]int)
				d.Decode(&kv.db)
				d.Decode(&kv.ack)
				kv.mu.Unlock()
			} else {
				op := msg.Command.(Op)
				kv.mu.Lock()
				if !kv.CheckDup(op.ID, op.ReqID) {
					kv.Apply(op)
				}

				ch, ok := kv.result[msg.Index]
				if ok {
					select {
					case <-kv.result[msg.Index]:
					default:
					}
					ch <- op
				} else {
					kv.result[msg.Index] = make(chan Op, 1)
				}

				//need snapshot
				if maxraftstate != -1 && kv.rf.GetPerisistSize() > maxraftstate {
					w := new(bytes.Buffer)
					e := gob.NewEncoder(w)
					e.Encode(kv.db)
					e.Encode(kv.ack)
					data := w.Bytes()
					go kv.rf.StartSnapshot(data, msg.Index)
				}
				kv.mu.Unlock()
			}
		}
	}()

	return kv
}

type Clerk struct {
	servers []*labrpc.ClientEnd
	// You will have to modify this struct.
	id    int64
	reqid int
	mu    sync.Mutex
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	// You'll have to add code here.
	ck.id = nrand()
	ck.reqid = 0
	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {

	// You will have to modify this function.
	var args GetArgs
	args.Key = key
	args.Id = ck.id
	ck.mu.Lock()
	args.ReqID = ck.reqid
	ck.reqid++
	ck.mu.Unlock()
	for {
		for _, v := range ck.servers {
			var reply GetReply
			ok := v.Call("RaftKV.Get", &args, &reply)
			if ok && reply.WrongLeader == false {
				//if reply.Err == ErrNoKey {
				//	reply.Value = ""
				//	}
				return reply.Value
			}
		}
	}
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	var args PutAppendArgs
	args.Key = key
	args.Value = value
	args.Op = op
	args.Id = ck.id
	ck.mu.Lock()
	args.ReqID = ck.reqid
	ck.reqid++
	ck.mu.Unlock()
	for {
		for _, v := range ck.servers {
			var reply PutAppendReply
			ok := v.Call("RaftKV.PutAppend", &args, &reply)
			if ok && reply.WrongLeader == false {
				return
			}
		}
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
