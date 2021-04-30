package core

import "sync"

/*
*	当前游戏世界总管理模块
 */

type WorldManager struct {
	// AOIManager 当前世界地图的AOI 管理模块
	PAoiMgr *AOIManager

	// 当前全部在线的 players 集合
	MapPlayers map[int32]*Player

	// 保护 players 的集合的锁
	pLock sync.RWMutex
}

// 提供一个对外的世界管理模块的句柄（对外唯一）
var PWorldMgrObj *WorldManager

//初始化方法  全局使用，故而使用 init 初始化
func init() {
	PWorldMgrObj = &WorldManager{
		PAoiMgr:    NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		MapPlayers: make(map[int32]*Player),
	}
}

//添加一个玩家
func (pWm *WorldManager) AddPlayer(pPlayer *Player)  {
	pWm.pLock.Lock()
	pWm.MapPlayers[pPlayer.Pid] = pPlayer
	pWm.pLock.Unlock()

	// 将玩家Player 添加到 AOIManager 中
	pWm.PAoiMgr.AddToGridByPos(int(pPlayer.Pid),pPlayer.X,pPlayer.Z)
}


//删除一个玩家
func (pWm *WorldManager) RemovePlayerByPid(pid int32) {

	// 先 从AOIManager 删除
	player := pWm.MapPlayers[pid]
	pWm.PAoiMgr.RemoveFromGridByPos(int(pid),player.X,player.Z)

	// 再将从世界管理中删除
	pWm.pLock.Lock()
	delete(pWm.MapPlayers,pid)
	pWm.pLock.Unlock()
}

//通过玩家ID查询Player对象
func (pWm *WorldManager) GetPlayerByPid(pid int32) *Player {

	pWm.pLock.Lock()
	defer  pWm.pLock.Unlock()
	return pWm.MapPlayers[pid]
}

//获取全部的在线玩家
func (pWm *WorldManager) GetAllPlayers() []*Player {

	pWm.pLock.Lock()
	defer  pWm.pLock.Unlock()
	players := make([]*Player,0)

	for _,p := range pWm.MapPlayers {
		players = append(players,p)
	}
	return players
}
