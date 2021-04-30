package core

import (
	"fmt"
	"sync"
)

/*
	一个 AOI 地图中的格子类型
*/
type Grid struct {
	// 格子ID
	IGID int
	//格子的左边边界坐标
	IMinX int
	//格子的右边边界坐标
	IMaxX int
	//格子的上边边界坐标
	IMinY int
	//格子的下边边界坐标
	IMaxY int
	//当前格子内玩家或者物体成员的ID集合
	mapPlayerIDs map[int]bool
	//当前集合的锁
	pIDLock sync.RWMutex
}

//初始化当前的格子的方法
func NewGrid(gID,minX,maxX,minY,maxY int) *Grid {
	return &Grid{
		IGID :gID,
		IMinX: minX,
		IMaxX: maxX,
		IMinY: minY,
		IMaxY: maxY,
		mapPlayerIDs : make(map[int]bool),
	}
}

//给格子添加一个玩家
func (pG *Grid) AddPlayer(iPlayerID int)  {
	pG.pIDLock.Lock()
	defer pG.pIDLock.Unlock()

	pG.mapPlayerIDs[iPlayerID] = true
}


//从格子中删除一个玩家
func (pG *Grid) RemovePlayer(iPlayerID int)  {
	pG.pIDLock.Lock()
	defer pG.pIDLock.Unlock()

	delete(  pG.mapPlayerIDs,iPlayerID)
}

//得到当前格子中所有的玩家ID
func (pG *Grid) GetAllPlayerIDs() (aPlayerIDs []int)  {
	pG.pIDLock.RLock()
	defer pG.pIDLock.RUnlock()

	for k,_:= range pG.mapPlayerIDs{
		aPlayerIDs = append(aPlayerIDs,k)
	}
	return
}

//调试使用打印出格子的基本信息 重写  String() 方法
// fmt.Println(Grid) 默认调用 grid. String()
func (pG *Grid) String() string  {

	//return  fmt.Sprintf("Grid %+v",pG)
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",pG.IGID,pG.IMinX,pG.IMaxX,pG.IMinY,pG.IMaxY,pG.mapPlayerIDs)
}




