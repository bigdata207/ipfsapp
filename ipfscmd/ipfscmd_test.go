package ipfscmd

import (
	"flag"
	"fmt"
	"github.com/ipfs/go-ipfs/commands"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsCientErr(t *testing.T) {
	t.Log("Catch both pointers and values")
	if !isClientError(commands.Error{Code: commands.ErrClient}) {
		t.Errorf("misidentified value")
	}
	if !isClientError(&commands.Error{Code: commands.ErrClient}) {
		t.Errorf("misidentified pointer")
	}
}

func TestRunMain(t *testing.T) {
	args := flag.Args()
	Args := append([]string{os.Args[0]}, args...)
	ret := mainRet(Args...)

	p := os.Getenv("IPFS_COVER_RET_FILE")
	if len(p) != 0 {
		ioutil.WriteFile(p, []byte(fmt.Sprintf("%d\n", ret)), 0777)
	}

	// close outputs so go testing doesn't print anything
	null, _ := os.Open(os.DevNull)
	os.Stderr = null
	os.Stdout = null
}
