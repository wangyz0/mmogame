package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManger(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	//
	g := aoiMgr.GetSurroundGridsByGid(24)
	for _, v := range g {
		fmt.Println(v.GID)
	}
}

