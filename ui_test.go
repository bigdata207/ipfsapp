package ipfsapp

import (
	"testing"

	"github.com/andlabs/ui"
)

func TestUI(t *testing.T) {
	err := ui.Main(LoginBox)
	if err != nil {
		panic(err)
	}
}
