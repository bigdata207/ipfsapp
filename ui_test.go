package ipfsapp

import (
	"testing"

	"github.com/andlabs/ui"
)

func TestUI(t *testing.T) {
	err := ui.Main(initUI)
	if err != nil {
		panic(err)
	}
}
