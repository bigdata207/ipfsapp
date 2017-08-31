package sraft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"github.com/mackzhong/ipfsapp/labrpc"
	"math/rand"
	"sync"
	"time"
)

//定义常量
const (
	LEADER    = iota // 0
	CANDIDATE        // 1
	FOLLOWER         // 2

	HBINTERVAL = 50 * time.Millisecond // 50ms
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
//接收信息结构体
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

//日志结构体
type LogEntry struct {
	LogIndex   int
	LogTerm    int
	LogCommand interface{}
}

//
// 一个raft节点实现.
//
type Raft struct {
	mu        sync.Mutex
	peers     []*labrpc.ClientEnd
	persister *Persister
	me        int // index into peers[]

	// Your data here.
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	//channel
	//信道变量，用于各节点通信
	state         int
	voteCount     int
	chanCommit    chan bool
	chanHeartbeat chan bool
	chanGrantVote chan bool
	chanLeader    chan bool
	chanApply     chan ApplyMsg

	//persistent state on all server
	//所有服务节点的持久状态
	currentTerm int
	votedFor    int
	log         []LogEntry

	//volatile state on all servers
	//所有节点上的可变状态
	commitIndex int
	lastApplied int

	//volatile state on leader
	//leader上的可变状态
	nextIndex  []int
	matchIndex []int
}

// 返回当前节点任期以及是否是leader
func (rf *Raft) GetState() (int, bool) {
	return rf.currentTerm, rf.state == LEADER
}

//返回最新的日志编号
func (rf *Raft) getLastIndex() int {
	return rf.log[len(rf.log)-1].LogIndex
}

//返回最新日志是在第几个任期保存的
func (rf *Raft) getLastTerm() int {
	return rf.log[len(rf.log)-1].LogTerm
}

//返回节点是否是Leader
func (rf *Raft) IsLeader() bool {
	return rf.state == LEADER
}

//将节点状态持久化存储到硬盘中从而可以在节点崩溃或者重启后恢复
func (rf *Raft) persist() {
	// Your code here.
	// Example:
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.log)
	data := w.Bytes()
	rf.persister.SaveRaftState(data)
}

//读取快照恢复节点状态
func (rf *Raft) readSnapshot(data []byte) {

	rf.readPersist(rf.persister.ReadRaftState())

	if len(data) == 0 {
		return
	}

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)

	var LastIncludedIndex int
	var LastIncludedTerm int

	d.Decode(&LastIncludedIndex)
	d.Decode(&LastIncludedTerm)

	rf.commitIndex = LastIncludedIndex
	rf.lastApplied = LastIncludedIndex

	rf.log = truncateLog(LastIncludedIndex, LastIncludedTerm, rf.log)

	msg := ApplyMsg{UseSnapshot: true, Snapshot: data}

	go func() {
		rf.chanApply <- msg
	}()
}

//还原原来的持久状态
func (rf *Raft) readPersist(data []byte) {
	// Your code here.
	// Example:
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.log)
}

//RequestVoteArgs 请求选票 的RPC结构体
type RequestVoteArgs struct {
	// Your data here.
	Term         int
	CandidateID  int
	LastLogTerm  int
	LastLogIndex int
}

//RequestVoteReply 回复选票请求的RPC结构体
type RequestVoteReply struct {
	// Your data here.
	Term        int
	VoteGranted bool
}

//AppendEntriesArgs 请求添加日志的RPC结构体
type AppendEntriesArgs struct {
	// Your data here.
	Term         int
	LeaderID     int
	PrevLogTerm  int
	PrevLogIndex int
	Entries      []LogEntry
	LeaderCommit int
}

//AppendEntriesReply 回复日志添加请求的RPC结构体
type AppendEntriesReply struct {
	// Your data here.
	Term      int
	Success   bool
	NextIndex int
}

//RequestVote 选票请求处理，args就是当前节点参数的一个集合
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer rf.persist()
	reply.VoteGranted = false
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	//如果中的任期大于打前节点的任期，则将当前节点设为follower,当前任期设为请求中的任期，默认选举人编号为-1
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = FOLLOWER
		rf.votedFor = -1
	}
	reply.Term = rf.currentTerm

	term := rf.getLastTerm()
	index := rf.getLastIndex()
	uptoDate := false

	//If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote
	//Raft determines which of two logs is more up-to-date by comparing the index and term of the last entries in the logs.
	// If the logs have last entries with different terms,then the log with the later term is more up-to-date.
	// If the logs end with the same term, then whichever log is longer is more up-to-date
	if args.LastLogTerm > term {
		uptoDate = true
	}

	if args.LastLogTerm == term && args.LastLogIndex >= index {
		// at least up to date
		uptoDate = true
	}

	if (rf.votedFor == -1 || rf.votedFor == args.CandidateID) && uptoDate {
		rf.chanGrantVote <- true
		rf.state = FOLLOWER
		reply.VoteGranted = true
		rf.votedFor = args.CandidateID
	}
}

//AppendEntries 处理日志添加请求
func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer rf.persist()

	reply.Success = false
	//如果日志添加请求里的任期小于当前节点任期，那么直接返回当前节点任期以及当前节点最新日志编号的下一位
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
		return
	}
	rf.chanHeartbeat <- true
	//If RPC request or response contains term T > currentTerm: set currentTerm = T, convert to follower
	//如果请求中的任期大于打前节点的任期，则将当前节点设为follower,当前任期设为请求中的任期，默认选举人编号为-1
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = FOLLOWER
		rf.votedFor = -1
	}
	reply.Term = args.Term
	//如果请求里记录的上一条日志的编号大于当前节点的最新日志编号，说明当前节点有日志缺失，直接返回当前节点第一条缺失日志的编号以便同步
	if args.PrevLogIndex > rf.getLastIndex() {
		reply.NextIndex = rf.getLastIndex() + 1
		return
	}

	baseIndex := rf.log[0].LogIndex

	// If a follower’s log is inconsistent with the leader’s, the AppendEntries consis- tency check will fail in the next AppendEntries RPC.
	// Af- ter a rejection, the leader decrements nextIndex and retries the AppendEntries RPC
	//Eventually nextIndex will reach a point where the leader and follower logs match
	//which removes any conflicting entries in the follower’s log and appends entries from the leader’s log (if any).
	if args.PrevLogIndex > baseIndex {
		term := rf.log[args.PrevLogIndex-baseIndex].LogTerm
		if args.PrevLogTerm != term {
			for i := args.PrevLogIndex - 1; i >= baseIndex; i-- {
				if rf.log[i-baseIndex].LogTerm != term {
					reply.NextIndex = i + 1
					break
				}
			}
			return
		}
	}
	if args.PrevLogIndex < baseIndex {

	} else {
		//Append any new entries not already in the log
		rf.log = rf.log[:args.PrevLogIndex+1-baseIndex]
		rf.log = append(rf.log, args.Entries...)
		reply.Success = true
		reply.NextIndex = rf.getLastIndex() + 1
	}
	//If leaderCommit > commitIndex, set commitIndex =min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		last := rf.getLastIndex()
		if args.LeaderCommit > last {
			rf.commitIndex = last
		} else {
			rf.commitIndex = args.LeaderCommit
		}
		rf.chanCommit <- true
	}
	return
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// returns true if labrpc says the RPC was delivered.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
//调用rpc服务发送选票请求
func (rf *Raft) sendRequestVote(server int, args RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		term := rf.currentTerm
		//如果节点状态已经不是candidate，直接返回请求是否发到
		if rf.state != CANDIDATE {
			return ok
		}
		if args.Term != term {
			return ok
		}
		//如果reply的任期比当前节点任期高,结点任期更新，状态改为follower
		if reply.Term > term {
			rf.currentTerm = reply.Term
			rf.state = FOLLOWER
			rf.votedFor = -1
			rf.persist()
		}
		//如果该节点同意，那么选票加1
		if reply.VoteGranted {
			rf.voteCount++
			//如果当前节点状态为candinate且选票大于一半的节点，那么当前节点变为leader
			if rf.state == CANDIDATE && rf.voteCount > len(rf.peers)/2 {
				rf.state = LEADER
				rf.chanLeader <- true
			}
		}
	}
	return ok
}

//调用rpc服务发送日志添加请求
func (rf *Raft) sendAppendEntries(server int, args AppendEntriesArgs, reply *AppendEntriesReply) bool {
	//ok代表请求是否发送到服务器，false表示无法连接到服务器,ture表示发送成功
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		//如果节点状态已经不是leader，直接返回ok
		if rf.state != LEADER {
			return ok
		}
		//如果日志中的任期不等于当前节点的任期，直接返回ok
		if args.Term != rf.currentTerm {
			return ok
		}
		//如果reply的任期大于当前节点，说明当前节点是老Leader，需要更新任期并转换为follower
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = FOLLOWER
			rf.votedFor = -1
			rf.persist()
			return ok
		}
		//如果获得回复则更新索引日志数组
		if reply.Success {
			//如果添加的日志数不为0则更新
			if len(args.Entries) > 0 {
				rf.nextIndex[server] = args.Entries[len(args.Entries)-1].LogIndex + 1
				//reply.NextIndex
				//rf.nextIndex[server] = reply.NextIndex
				rf.matchIndex[server] = rf.nextIndex[server] - 1
			}
		} else {
			rf.nextIndex[server] = reply.NextIndex
		}
	}
	return ok
}

//InstallSnapshotArgs 快照相关部分，用于在节点奔溃或者重启后的状态恢复
type InstallSnapshotArgs struct {
	Term              int //任期
	LeaderID          int //leader编号
	LastIncludedIndex int //两个节点重合日志的最大编号
	LastIncludedTerm  int //两个节点重合的编号最大的日志的任期，这两个参数可以通过truncatelog函数获取两个节点日志中其中一个少掉的那部分日志然后从另一个节点同步回去
	Data              []byte
}

type InstallSnapshotReply struct {
	Term int
}

func (rf *Raft) GetPerisistSize() int {
	return rf.persister.RaftStateSize()
}

func (rf *Raft) StartSnapshot(snapshot []byte, index int) {

	rf.mu.Lock()
	defer rf.mu.Unlock()

	baseIndex := rf.log[0].LogIndex
	lastIndex := rf.getLastIndex()

	if index <= baseIndex || index > lastIndex {
		// in case having installed a snapshot from leader before snapshotting
		// second condition is a hack
		return
	}

	var newLogEntries []LogEntry

	newLogEntries = append(newLogEntries, LogEntry{LogIndex: index, LogTerm: rf.log[index-baseIndex].LogTerm})

	for i := index + 1; i <= lastIndex; i++ {
		newLogEntries = append(newLogEntries, rf.log[i-baseIndex])
	}

	rf.log = newLogEntries

	rf.persist()

	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(newLogEntries[0].LogIndex)
	e.Encode(newLogEntries[0].LogTerm)

	data := w.Bytes()
	data = append(data, snapshot...)
	rf.persister.SaveSnapshot(data)
}

//获取两个节点日志中其中一个少掉的那部分日志并返回
func truncateLog(lastIncludedIndex int, lastIncludedTerm int, log []LogEntry) []LogEntry {

	var newLogEntries []LogEntry
	//先把最新日志存到新数组
	newLogEntries = append(newLogEntries, LogEntry{LogIndex: lastIncludedIndex, LogTerm: lastIncludedTerm})
	//找到第一条编号和任期和参数中的一样的日志，从该日志的索引往后的所有日志都截取出来合并到新数组中返回
	for index := len(log) - 1; index >= 0; index-- {
		if log[index].LogIndex == lastIncludedIndex && log[index].LogTerm == lastIncludedTerm {
			newLogEntries = append(newLogEntries, log[index+1:]...)
			break
		}
	}
	return newLogEntries
}

//InstallSnapshot 节点根据给定快照参数从快照中恢复快照中的状态
func (rf *Raft) InstallSnapshot(args InstallSnapshotArgs, reply *InstallSnapshotReply) {
	// Your code here.
	rf.mu.Lock()
	defer rf.mu.Unlock()
	//如果参数即快照中的任期小于当前节点，直接将当前节点任期更新并return
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}
	//更新节点参数
	rf.chanHeartbeat <- true
	rf.state = FOLLOWER
	//rf.currentTerm = rf.currentTerm
	rf.currentTerm = args.Term
	//保存参数数据
	rf.persister.SaveSnapshot(args.Data)

	//获取不同的日志部分
	rf.log = truncateLog(args.LastIncludedIndex, args.LastIncludedTerm, rf.log)

	msg := ApplyMsg{UseSnapshot: true, Snapshot: args.Data}

	rf.lastApplied = args.LastIncludedIndex
	rf.commitIndex = args.LastIncludedIndex

	rf.persist()

	rf.chanApply <- msg
}

//调用rpc服务将自己的节点状态打包发送给其他节点
func (rf *Raft) sendInstallSnapshot(server int, args InstallSnapshotArgs, reply *InstallSnapshotReply) bool {
	ok := rf.peers[server].Call("Raft.InstallSnapshot", args, reply)
	if ok {
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = FOLLOWER
			rf.votedFor = -1
			return ok
		}

		rf.nextIndex[server] = args.LastIncludedIndex + 1
		rf.matchIndex[server] = args.LastIncludedIndex
	}
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	index := -1
	term := rf.currentTerm
	isLeader := rf.state == LEADER
	if isLeader {
		index = rf.getLastIndex() + 1
		//fmt.Printf("raft:%d start\n",rf.me)
		rf.log = append(rf.log, LogEntry{LogTerm: term, LogCommand: command, LogIndex: index}) // append new entry from client
		rf.persist()
	}
	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//候选人拉取选票
func (rf *Raft) broadcastRequestVote() {
	var args RequestVoteArgs
	rf.mu.Lock()
	args.Term = rf.currentTerm
	args.CandidateID = rf.me
	args.LastLogTerm = rf.getLastTerm()
	args.LastLogIndex = rf.getLastIndex()
	rf.mu.Unlock()

	for i := range rf.peers {
		//当i不等于自己的编号并且自己还是candidate的时候不断发送选票请求
		if i != rf.me && rf.state == CANDIDATE {
			go func(i int) {
				var reply RequestVoteReply
				rf.sendRequestVote(i, args, &reply)
			}(i)
		}
	}
}

/**
 * Log replication
 */
//广播日志复制
func (rf *Raft) broadcastAppendEntries() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	//广播日志的编号
	N := rf.commitIndex
	//最新一条日志的编号
	last := rf.getLastIndex()
	//第一条日志的编号
	baseIndex := rf.log[0].LogIndex
	//If there exists an N such that N > commitIndex, a majority of matchIndex[i] ≥ N, and log[N].term == currentTerm: set commitIndex = N
	//找到最新的大于当前节点提交日志编号的N能够使得有一半以上节点最新日志编号大于等于commitIndex且第i-baseIndex条日志的任期等于当前任期
	for i := rf.commitIndex + 1; i <= last; i++ {
		num := 1
		for j := range rf.peers {
			if j != rf.me && rf.matchIndex[j] >= i && rf.log[i-baseIndex].LogTerm == rf.currentTerm {
				num++
			}
		}
		if 2*num > len(rf.peers) {
			N = i
		}
	}
	//触发Make中的第二个子进程，传递消息回复给applyCh
	if N != rf.commitIndex {
		rf.commitIndex = N
		rf.chanCommit <- true
	}

	for i := range rf.peers {
		//当编号不是自己并且自己的状态依旧是leader时不断遍历向所有节点发送日志添加请求
		if i != rf.me && rf.state == LEADER {

			//copy(args.Entries, rf.log[args.PrevLogIndex + 1:])

			if rf.nextIndex[i] > baseIndex {
				var args AppendEntriesArgs
				args.Term = rf.currentTerm
				args.LeaderID = rf.me
				args.PrevLogIndex = rf.nextIndex[i] - 1
				//	fmt.Printf("baseIndex:%d PrevLogIndex:%d\n",baseIndex,args.PrevLogIndex )
				args.PrevLogTerm = rf.log[args.PrevLogIndex-baseIndex].LogTerm
				//args.Entries = make([]LogEntry, len(rf.log[args.PrevLogIndex + 1:]))
				args.Entries = make([]LogEntry, len(rf.log[args.PrevLogIndex+1-baseIndex:]))
				copy(args.Entries, rf.log[args.PrevLogIndex+1-baseIndex:])
				args.LeaderCommit = rf.commitIndex
				go func(i int, args AppendEntriesArgs) {
					var reply AppendEntriesReply
					rf.sendAppendEntries(i, args, &reply)
				}(i, args)
			} else {
				var args InstallSnapshotArgs
				args.Term = rf.currentTerm
				args.LeaderID = rf.me
				args.LastIncludedIndex = rf.log[0].LogIndex
				args.LastIncludedTerm = rf.log[0].LogTerm
				args.Data = rf.persister.snapshot
				go func(server int, args InstallSnapshotArgs) {
					reply := &InstallSnapshotReply{}
					rf.sendInstallSnapshot(server, args, reply)
				}(i, args)
			}
		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here.
	rf.state = FOLLOWER
	rf.votedFor = -1
	rf.log = append(rf.log, LogEntry{LogTerm: 0})
	rf.currentTerm = 0
	rf.chanCommit = make(chan bool, 100)
	rf.chanHeartbeat = make(chan bool, 100)
	rf.chanGrantVote = make(chan bool, 100)
	rf.chanLeader = make(chan bool, 100)
	rf.chanApply = applyCh

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	rf.readSnapshot(persister.ReadSnapshot())
	//每个raft节点都有一个负责切换节点状态的子进程
	go func() {
		for {
			switch rf.state {
			case FOLLOWER:
				select {
				//转化为candidate的三种情况
				case <-rf.chanHeartbeat:
				case <-rf.chanGrantVote:
				case <-time.After(time.Duration(rand.Int63()%333+550) * time.Millisecond):
					rf.state = CANDIDATE
				}
			//leader广播日志添加请求,	然后等待50ms回复
			case LEADER:
				//fmt.Printf("Leader:%v %v\n",rf.me,"boatcastAppendEntries	")
				rf.broadcastAppendEntries()
				time.Sleep(HBINTERVAL)
			//candidate 会在当前任期上加1,然后广播选票请求争取当上leader走向人上巅峰
			case CANDIDATE:
				rf.mu.Lock()
				//To begin an election, a follower increments its current term and transitions to candidate state
				rf.currentTerm++
				//It then votes for itself and issues RequestVote RPCs in parallel to each of the other servers in the cluster.
				rf.votedFor = rf.me
				rf.voteCount = 1
				rf.persist()
				rf.mu.Unlock()
				//(a) it wins the election, (b) another server establishes itself as leader, or (c) a period of time goes by with no winner
				go rf.broadcastRequestVote()
				select {
				//随机等待一个介于510~810ms之间的时间间隔
				case <-time.After(time.Duration(rand.Int63()%300+510) * time.Millisecond):
				case <-rf.chanHeartbeat:
					rf.state = FOLLOWER
				//	fmt.Printf("CANDIDATE %v reveive chanHeartbeat\n",rf.me)
				//选票达到一半并且原来还是candidate时转化为leadert
				case <-rf.chanLeader:
					rf.mu.Lock()
					rf.state = LEADER
					//fmt.Printf("%v is Leader\n",rf.me)
					rf.nextIndex = make([]int, len(rf.peers))
					rf.matchIndex = make([]int, len(rf.peers))
					for i := range rf.peers {
						//The leader maintains a nextIndex for each follower, which is the index of the next log entry the leader will send to that follower.
						// When a leader first comes to power, it initializes all nextIndex values to the index just after the last one in its log
						rf.nextIndex[i] = rf.getLastIndex() + 1
						rf.matchIndex[i] = 0
					}
					rf.mu.Unlock()
					//rf.boatcastAppendEntries()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-rf.chanCommit:
				rf.mu.Lock()
				commitIndex := rf.commitIndex
				baseIndex := rf.log[0].LogIndex
				for i := rf.lastApplied + 1; i <= commitIndex; i++ {
					msg := ApplyMsg{Index: i, Command: rf.log[i-baseIndex].LogCommand}
					applyCh <- msg
					//fmt.Printf("me:%d %v\n",rf.me,msg)
					rf.lastApplied = i
				}
				rf.mu.Unlock()
			}
		}
	}()
	return rf
}
