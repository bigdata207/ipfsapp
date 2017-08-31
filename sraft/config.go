package sraft

//
// support for Raft tester.
//
// we will use the original config.go to test your code for grading.
// so, while you can modify this code to help you debug, please
// test with the original before submitting.
//

import "github.com/mackzhong/ipfsapp/labrpc"
import "log"
import "sync"
import "testing"
import "runtime"

//import crand "crypto/rand"
import crand "github.com/mackzhong/ipfsapp/fastrand"
import "encoding/base64"
import "sync/atomic"
import "time"
import (
	"fmt"
	"math/rand"
)

//产生随机字符串做终端名
func randstring(n int) string {
	b := make([]byte, 2*n)
	crand.Read(b)
	s := base64.URLEncoding.EncodeToString(b)
	return s[0:n]
}

//管理中心......
type config struct {
	mu        sync.Mutex
	t         *testing.T
	net       *labrpc.Network
	n         int
	done      int32 // tell internal threads to die
	rafts     []*Raft
	applyErr  []string // from apply channel readers
	connected []bool   // whether each server is on the net
	saved     []*Persister
	endnames  [][]string    // the port file names each sends to
	logs      []map[int]int // copy of each server's committed entries
}

//初始化服务
func make_config(t *testing.T, n int, unreliable bool) *config {
	runtime.GOMAXPROCS(4)
	cfg := &config{}
	cfg.t = t
	cfg.net = labrpc.MakeNetwork()
	cfg.n = n
	cfg.applyErr = make([]string, cfg.n)
	cfg.rafts = make([]*Raft, cfg.n)
	cfg.connected = make([]bool, cfg.n)
	cfg.saved = make([]*Persister, cfg.n)
	cfg.endnames = make([][]string, cfg.n)
	cfg.logs = make([]map[int]int, cfg.n)

	cfg.setunreliable(unreliable)

	cfg.net.LongDelays(true)

	//启动所有节点
	// create a full set of Rafts.
	for i := 0; i < cfg.n; i++ {
		cfg.logs[i] = map[int]int{}
		cfg.start1(i)
	}
	//将所有节点加到服务网络中
	// connect everyone
	for i := 0; i < cfg.n; i++ {
		cfg.connect(i)
	}

	return cfg
}

// shut down a Raft server but save its persistent state.
//关闭一个Raft节点但是保留他的状态

func (cfg *config) crash1(i int) {
	//与服务网络断开
	cfg.disconnect(i)
	//从服务网络中删除
	cfg.net.DeleteServer(i) // disable client connections to the server.

	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// a fresh persister, in case old instance
	// continues to update the Persister.
	// but copy old persister's content so that we always
	// pass Make() the last persisted state.
	if cfg.saved[i] != nil {
		//获取原来状态的一个备份
		cfg.saved[i] = cfg.saved[i].Copy()
	}

	rf := cfg.rafts[i]
	if rf != nil {
		cfg.mu.Unlock()
		rf.Kill()
		cfg.mu.Lock()
		cfg.rafts[i] = nil
	}

	if cfg.saved[i] != nil {
		raftlog := cfg.saved[i].ReadRaftState()
		cfg.saved[i] = &Persister{}
		cfg.saved[i].SaveRaftState(raftlog)
	}
}

//
// start or re-start a Raft.
// if one already exists, "kill" it first.
// allocate new outgoing port file names, and a new
// state persister, to isolate previous instance of
// this server. since we cannot really kill it.
//
//启动一个Raft节点
func (cfg *config) start1(i int) {
	//先关闭该节点
	cfg.crash1(i)

	// a fresh set of outgoing ClientEnd names.
	// so that old crashed instance's ClientEnds can't send.
	cfg.endnames[i] = make([]string, cfg.n)
	for j := 0; j < cfg.n; j++ {
		cfg.endnames[i][j] = randstring(20)
	}

	// a fresh set of ClientEnds.
	//添加到服务网络中
	ends := make([]*labrpc.ClientEnd, cfg.n)
	for j := 0; j < cfg.n; j++ {
		ends[j] = cfg.net.MakeEnd(cfg.endnames[i][j])
		cfg.net.Connect(cfg.endnames[i][j], j)
	}

	cfg.mu.Lock()

	// a fresh persister, so old instance doesn't overwrite
	// new instance's persisted state.
	// but copy old persister's content so that we always
	// pass Make() the last persisted state.
	if cfg.saved[i] != nil {
		cfg.saved[i] = cfg.saved[i].Copy()
	} else {
		cfg.saved[i] = MakePersister()
	}

	cfg.mu.Unlock()

	// listen to messages from Raft indicating newly committed messages.
	applyCh := make(chan ApplyMsg)
	go func() {
		for m := range applyCh {
			err_msg := ""
			if m.UseSnapshot {
				// ignore the snapshot
			} else if v, ok := (m.Command).(int); ok {
				cfg.mu.Lock()
				for j := 0; j < len(cfg.logs); j++ {
					if old, oldok := cfg.logs[j][m.Index]; oldok && old != v {
						// some server has already committed a different value for this entry!
						err_msg = fmt.Sprintf("commit index=%v server=%v %v != server=%v %v",
							m.Index, i, m.Command, j, old)
					}
				}
				_, prevok := cfg.logs[i][m.Index-1]
				cfg.logs[i][m.Index] = v
				cfg.mu.Unlock()

				if m.Index > 1 && prevok == false {
					err_msg = fmt.Sprintf("server %v apply out of order %v", i, m.Index)
				}
			} else {
				err_msg = fmt.Sprintf("committed command %v is not an int", m.Command)
			}

			if err_msg != "" {
				log.Fatalf("apply error: %v\n", err_msg)
				cfg.applyErr[i] = err_msg
				// keep reading after error so that Raft doesn't block
				// holding locks...
			}
		}
	}()

	rf := Make(ends, i, cfg.saved[i], applyCh)

	cfg.mu.Lock()
	cfg.rafts[i] = rf
	cfg.mu.Unlock()

	svc := labrpc.MakeService(rf)
	srv := labrpc.MakeServer()
	srv.AddService(svc)
	cfg.net.AddServer(i, srv)
}

//停止所有节点
func (cfg *config) cleanup() {
	for i := 0; i < len(cfg.rafts); i++ {
		if cfg.rafts[i] != nil {
			cfg.rafts[i].Kill()
		}
	}
	atomic.StoreInt32(&cfg.done, 1)
}

// attach server i to the net.
//将i节点加入网络，原理和disconnect相同
func (cfg *config) connect(i int) {
	// fmt.Printf("connect(%d)\n", i)

	cfg.connected[i] = true

	// outgoing ClientEnds
	for j := 0; j < cfg.n; j++ {
		if cfg.connected[j] {
			endname := cfg.endnames[i][j]
			cfg.net.Enable(endname, true)
		}
	}

	// incoming ClientEnds
	for j := 0; j < cfg.n; j++ {
		if cfg.connected[j] {
			endname := cfg.endnames[j][i]
			cfg.net.Enable(endname, true)
		}
	}
}

// detach server i from the net.
//将编号i的节点从网络中断开
func (cfg *config) disconnect(i int) {
	// fmt.Printf("disconnect(%d)\n", i)
	//将连接状态设为false
	cfg.connected[i] = false

	// outgoing ClientEnds
	//关闭所有i节点对其他节点的输出
	for j := 0; j < cfg.n; j++ {
		if cfg.endnames[i] != nil {
			endname := cfg.endnames[i][j]
			cfg.net.Enable(endname, false)
		}
	}

	// incoming ClientEnds
	//关闭所有其他节点对i节点的输入
	for j := 0; j < cfg.n; j++ {
		if cfg.endnames[j] != nil {
			endname := cfg.endnames[j][i]
			cfg.net.Enable(endname, false)
		}
	}
}

//获取服务节点数
func (cfg *config) rpcCount(server int) int {
	return cfg.net.GetCount(server)
}

func (cfg *config) setunreliable(unrel bool) {
	cfg.net.Reliable(!unrel)
}

func (cfg *config) setlongreordering(longrel bool) {
	cfg.net.LongReordering(longrel)
}

// check that there's exactly one leader.
// try a few times in case re-elections are needed.
//检查是否只有不超过1个Leader
func (cfg *config) checkOneLeader() int {
	for iters := 0; iters < 10; iters++ {
		leaders := make(map[int][]int)
		//遍历所有节点找出Leader的节点(可能不止一个)
		for i := 0; i < cfg.n; i++ {
			if cfg.connected[i] {
				if t, leader := cfg.rafts[i].GetState(); leader {
					leaders[t] = append(leaders[t], i)
				}
			}
		}
		//保存最新的leader的任期，即最大的任期
		lastTermWithLeader := -1
		fmt.Println("hhh")
		fmt.Println(len(leaders))
		for t, leadersatt := range leaders {
			if len(leadersatt) > 1 {
				cfg.t.Fatalf("term %d has %d (>1) leaders", t, len(leadersatt))
			}
			if t > lastTermWithLeader {
				//更新lastTermWithLeader
				lastTermWithLeader = t
			}
		}
		//
		if len(leaders) != 0 {
			fmt.Println(leaders[lastTermWithLeader][0])
			fmt.Println(leaders)
			//返回具有最大任期的Leader的编号
			return leaders[lastTermWithLeader][0]
		}
		//没找到的话等待500ms选举新的leader
		time.Sleep(500 * time.Millisecond)
	}
	cfg.t.Fatalf("expected one leader, got none")
	//十次循环(5s)内没找到leader，返回-1
	return -1
}

// check that everyone agrees on the term.
func (cfg *config) checkTerms() int {
	term := -1
	//遍历所有节点
	for i := 0; i < cfg.n; i++ {
		//对服务中的节点进行判断
		if cfg.connected[i] {
			xterm, _ := cfg.rafts[i].GetState()
			if term == -1 {
				term = xterm
			} else if term != xterm {
				//如果当前节点与之前节点任期不同，则说明网络中存在异常节点(可能为condinate)
				cfg.t.Fatalf("servers disagree on term")
			}
		}
	}
	//返回当前任期
	return term
}

// check that there's no leader
func (cfg *config) checkNoLeader() {
	for i := 0; i < cfg.n; i++ {
		if cfg.connected[i] {
			_, isLeader := cfg.rafts[i].GetState()
			if isLeader {
				//cfg.t.Fatalf("expected no leader, but %v claims to be leader", i)
				fmt.Printf("expected no leader, but %v claims to be leader\n", i)
			}
		}
	}
}

// how many servers think a log entry is committed?
//统计有多少人已经存有第index条记录并且相同
func (cfg *config) nCommitted(index int) (int, interface{}) {
	//初始化已保存第index条日志的节点数以及保存的日志内容
	count := 0
	cmd := -1
	//遍历所有网络中的节点(正常服务节点)
	for i := 0; i < len(cfg.rafts); i++ {
		if cfg.applyErr[i] != "" {
			cfg.t.Fatal(cfg.applyErr[i])
		}
		//加锁，提取出第i节点的编号index的日志
		cfg.mu.Lock()
		cmd1, ok := cfg.logs[i][index]
		cfg.mu.Unlock()
		//DPrintf("i=",i,"index=",index,"cmd=",cmd1)
		//如果存在编号index的日志
		if ok {
			//如果已存在存有编号index日志的节点但是日志内容与当前节点保存的不一致，则说明日志不一致
			if count > 0 && cmd != cmd1 {
				cfg.t.Fatalf("committed values do not match: index %v, %v, %v\n",
					index, cmd, cmd1)
			}
			//否则已保存日志的节点数加1
			count += 1
			//更新cmd为当前节点保存编号index的日志
			cmd = cmd1
		}
	}
	return count, cmd
}

// wait for at least n servers to commit.
// but don't wait forever.
func (cfg *config) wait(index int, n int, startTerm int) interface{} {
	to := 10 * time.Millisecond
	for iters := 0; iters < 30; iters++ {
		nd, _ := cfg.nCommitted(index)
		//如果保存节点数已达到期望的n,则结束循环进行下一步
		if nd >= n {
			break
		}
		time.Sleep(to)
		//to = to * 2^iters 但是最多是1s
		if to < time.Second {
			to *= 2
		}
		if startTerm > -1 {
			for _, r := range cfg.rafts {
				//如果某个节点的任期大于编号index的日志广播时的任期，说明有了新的leader，
				// 并且该日志保存节点数少于期望的n,日志作废
				if t, _ := r.GetState(); t > startTerm {
					// someone has moved on
					// can no longer guarantee that we'll "win"
					return -1
				}
			}
		}
	}
	nd, cmd := cfg.nCommitted(index)
	//如果在规定时间内（即设定的30次循环内)记录节点数没有达到期望的n,则记录失败
	if nd < n {
		cfg.t.Fatalf("only %d decided for index %d; wanted %d\n",
			nd, index, n)
	}
	//返回节点中编号index的日志
	return cmd
}

// do a complete agreement.
// it might choose the wrong leader initially,
// and have to re-submit after giving up.
// entirely gives up after about 10 seconds.
// indirectly checks that the servers agree on the
// same value, since nCommitted() checks this,
// as do the threads that read from applyCh.
// returns index.
//保存一条日志并设置需要的最少已保存该日志的节点数，返回日志编号
func (cfg *config) one(cmd int, expectedServers int) int {

	for i := 0; i < len(cfg.rafts); i++ {
		//DPrintf(cfg.rafts)
	}

	t0 := time.Now()
	r := rand.New(rand.NewSource(t0.UnixNano()))
	fmt.Println("New loop")
	//随机一个初始结点查找Leader
	starts := r.Intn(5)
	//找到Leader广播一条日记记录要求大家记录
	for time.Since(t0).Seconds() < 10 {
		// try all the servers, maybe one is the leader.
		index := -1
		for si := 0; si < cfg.n; si++ {
			//取余保证每个节点都会被检测到
			starts = (starts + 1) % cfg.n
			fmt.Println(starts)
			var rf *Raft
			//锁定网络状态
			cfg.mu.Lock()
			//判断该节点是否正常工作，即在服务中
			if cfg.connected[starts] {
				rf = cfg.rafts[starts]
			}
			//解锁网络状态
			cfg.mu.Unlock()
			if rf != nil {
				//index1表示当前广播日志编号，ok表示当前结点是否是Leader
				index1, _, ok := rf.Start(cmd)
				if ok {
					//Leader已广播，index设置为当前广播日志编号
					index = index1
					break
				}
			}
		}
		//如果Leader已经广播，那么说明已经获取到日志编号，统计2s内已保存日志的节点数
		if index != -1 {
			// somebody claimed to be the leader and to have
			// submitted our command; wait a while for agreement.
			t1 := time.Now()
			//如果已经广播，那么设置2s的等待保存时间
			for time.Since(t1).Seconds() < 2 {
				nd, cmd1 := cfg.nCommitted(index)
				fmt.Printf("%d : %v\n", nd, cmd1)
				//如果节点数>0并且大于所需确定的节点数
				if nd > 0 && nd >= expectedServers {
					// 验证是否已经提交并且保存的日志为所广播的日志
					if cmd2, ok := cmd1.(int); ok && cmd2 == cmd {
						return index
					}
				}
				//每20ms进行一次检测
				time.Sleep(20 * time.Millisecond)
			}
		} else {
			//如果没找到Leader，设置等待50ms
			time.Sleep(50 * time.Millisecond)
		}
	}
	printLogs(cfg)
	cfg.t.Fatalf("one(%v) failed to reach agreement", cmd)
	return -1
}

func printLogs(cfg *config) {
	for si := 0; si < cfg.n; si++ {
		//var rf *Raft = cfg.rafts[si]
		//fmt.Printf("%s",rf.role)
		//DPrintf(strconv.Itoa(rf.me),"data=",rf.logs)
	}
}
