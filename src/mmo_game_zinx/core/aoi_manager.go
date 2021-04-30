package core

import (
	"fmt"
)

// 定义一些 AOI 边界值
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 10
)

/*
	AOI 区域管理模块
*/

type AOIManager struct {
	//区域的左边界坐标
	IMinX int
	//区域的右边界坐标
	IMaxX int
	//X方向格子的数量
	ICntsX int
	//区域的上边界坐标
	IMinY int
	//区域的下边界坐标
	IMaxY int
	//Y方向格子的数量
	ICntsY int
	//当前区域中有哪些格子 [map-(key = 格子的ID),(Value= 格子对象)]
	mapGrids map[int]*Grid
}

// 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	objAoiMgr := &AOIManager{
		IMinX:    minX,
		IMaxX:    maxX,
		ICntsX:   cntsX,
		IMinY:    minY,
		IMaxY:    maxY,
		ICntsY:   cntsY,
		mapGrids: make(map[int]*Grid),
	}

	// 给AOI 初始化区域的格子所有格子进行编号 和 初始化动作
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子的ID，根据 x,y 编号
			// 格子的编号 ： id = idy * cntsX  + idx
			gid := y*cntsX + x
			// 初始化 gid 格子
			objAoiMgr.mapGrids[gid] = NewGrid(gid,
				objAoiMgr.IMinX+x*objAoiMgr.GetGridWidthX(),
				objAoiMgr.IMinX+(x+1)*objAoiMgr.GetGridWidthX(),
				objAoiMgr.IMinY+y*objAoiMgr.GetGridLengthY(),
				objAoiMgr.IMinY+(y+1)*objAoiMgr.GetGridLengthY())
		}
	}
	return objAoiMgr
}

// 得到每个格子在 X 轴方向的宽度
func (pM *AOIManager) GetGridWidthX() int {
	return (pM.IMaxX - pM.IMinX) / pM.ICntsX
}

// 得到每个格子在 Y 轴方向的长度
func (pM *AOIManager) GetGridLengthY() int {
	return (pM.IMaxY - pM.IMinY) / pM.ICntsY
}

// 打印格子信息
func (pM *AOIManager) String() string {
	// 打印 AOIManager 信息
	s := fmt.Sprintf("AOIManager:\n minX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d \n Grids in AOIManager\n",
		pM.IMinX, pM.IMaxX, pM.ICntsX, pM.IMinY, pM.IMaxY, pM.ICntsY)

	// 打印 全部格子信息
	for _, grid := range pM.mapGrids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// 根据格子GID 得到周边九宫格格子的ID集合
func (pM *AOIManager) GetSurroundGridsByGid(iGID int) (aGrids []*Grid) {
	// 判断gid (iGID)是否在 AOIManager 中
	if _, ok := pM.mapGrids[iGID]; !ok {
		return
	}

	// 初始化 grids 返回值切片, 将 当前 gid 本身加入 九宫格切片
	aGrids = append(aGrids, pM.mapGrids[iGID])

	// 判断 gid 左边和右边是否有格子
	// 通过 gid 得到当前格子 x 轴的编号
	idx := iGID % pM.ICntsX
	// 判断 idx 编号 坐标是否左边还有格子，有就放入 gidx 集合
	if idx > 0 {
		aGrids = append(aGrids, pM.mapGrids[iGID-1])
	}
	// 判断 idx 编号 坐标是否右边还有格子，有就放入 gidx 集合
	if idx < pM.ICntsX-1 {
		aGrids = append(aGrids, pM.mapGrids[iGID+1])
	}
	// 将X轴当前的格子都取出，进行遍历，再分别得到每个格子上下是否还有格子
	gidsX := make([]int, 0, len(aGrids))
	for _, v := range aGrids {
		gidsX = append(gidsX, v.IGID)
	}

	// 遍历 gidx 集合中每个格子的 gid
	for _, v := range gidsX {
		// 得到当前格子的 ID的 y 轴的编号
		idy := v / pM.ICntsY
		// gid 上边和下边 是否还有格子
		if idy > 0 {
			// 上边有格子
			aGrids = append(aGrids, pM.mapGrids[v-pM.ICntsX])
		}
		if idy < pM.ICntsY-1 {
			// 下边有格子
			aGrids = append(aGrids, pM.mapGrids[v+pM.ICntsX])
		}
	}
	return
}

// 通过 横纵轴坐标 x,y 得到 GID格子编号
func (pM *AOIManager) GetGidByPos(x, y float32) int {

	idx := (int(x) - pM.IMinX) / pM.GetGridWidthX()
	idy := (int(x) - pM.IMinY) / pM.GetGridLengthY()

	return idy*pM.ICntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部的playerIDs
func (pM *AOIManager) GetPidsByPos(x, y float32) (aPlayerIDs []int) {

	// 得到当前玩家的GID格子
	iGid := pM.GetGidByPos(x, y)

	// 通过GID得到周边九宫格信息
	aGrids := pM.GetSurroundGridsByGid(iGid)

	// 将九宫格的信息 全部的 Player 放在 aPlayerIDs 中
	for _, Grid := range aGrids {
		aPlayerIDs = append(aPlayerIDs, Grid.GetAllPlayerIDs()...)
	}
	return
}

//添加一个 playerID 到一个格子中
func (pM *AOIManager) AddPidToGrid(pID, gID int) {
	pM.mapGrids[gID].AddPlayer(pID)
}

//移除一个格子中的 playerID
func (pM *AOIManager) RemovePidFromGrid(pID, gID int) {
	pM.mapGrids[gID].RemovePlayer(pID)
}

//通过GID 获取全部的 playerID
func (pM *AOIManager) GetPidsByGid(gID int) (aPlayerIDs []int) {
	aPlayerIDs = pM.mapGrids[gID].GetAllPlayerIDs()
	return
}

//通过坐标 将 playerID 添加一个格子中
func (pM *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := pM.GetGidByPos(x, y)
	objGrid := pM.mapGrids[gID]
	objGrid.AddPlayer(pID)
}

//通过坐标把一个 player 从一个格子中 删除
func (pM *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := pM.GetGidByPos(x, y)
	objGrid := pM.mapGrids[gID]
	objGrid.RemovePlayer(pID)
}
