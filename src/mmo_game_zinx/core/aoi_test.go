package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化 AOIManager
	objAoiMgr := NewAOIManager(0,250,5,0,250,5)


	// 打印 AOIManager
	fmt.Println(objAoiMgr)
}


func TestAOIManagerSuroundGridsByGid(t *testing.T) {
	// 初始化 AOIManager
	objAoiMgr := NewAOIManager(0,250,5,0,250,5)

	for gid,_ := range objAoiMgr.mapGrids{
		grids := objAoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid: ",gid," grids len = ",len(grids))
		gIDs := make([]int,0,len(grids))
		for _,grid := range grids{
			gIDs = append(gIDs,grid.IGID)
		}
		fmt.Println("surounding grid ids are ",gIDs)
	}

	// 打印 AOIManager
	fmt.Println(objAoiMgr)
}







