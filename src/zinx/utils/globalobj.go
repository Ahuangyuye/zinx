package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/src/zinx/ziface"
)

/*
	存储一切有关 zinx 框架的全局参数，共其他模块使用
	一些参数可以通过 zinx.json 由用户配置
*/

type GLobalObj struct {
	/*
	server
	*/
	TcpServer ziface.IServer // 当前zinx 全局的 server 对象
	StrHostIP string					// 当前服务器主机监听的IP
	ITcpPort int		// 当前服务器主机监听的port
	StrName string 			// 当前服务器主的名称
	/*
		zinx
	*/
	StrVersion string 	// 当前zinx 的版本号
	IMaxConn int 	// 当前zinx 主机所允许最大连接数
	IMaxPackageSize uint32 	// 当前zinx 数据包的最大值


}

/*
	定义一个全局的对外 GlobalObj
*/
var GlobalObject *GLobalObj



// 从 zinx.json 去 加载用于自定义的参数
func (pG *GLobalObj)Reload()  {

	data,err :=ioutil.ReadFile("conf/zinx.json")

	//strFile, _ := exec.LookPath(os.Args[0])
	//strPath, _ := filepath.Abs(strFile)
	//index := strings.LastIndex(strPath, string(os.PathSeparator))
	//strPath = strPath[:index]
	//datapath := path.Join(strPath, "conf/zinx.json")
	//fmt.Println(" Reload conf path",datapath)
	//data,err :=ioutil.ReadFile(datapath)

	if err != nil{
		panic(err)
		//fmt.Println(" Reload conf/zinx.json err=",err)
	}

	// 将json 文件数据解析到struct 中
	err = json.Unmarshal(data,&GlobalObject)
	if err != nil{
		panic(err)
		//fmt.Println(" Reload conf/zinx.json err=",err)
	}
}


/*
	提供一个 init 方法，初始化当前的 pGlobalObject
*/
func init()  {
	// 如果配置文件没有加载，默认值
	GlobalObject = &GLobalObj {
		StrName: "ZinxServerApp",
		StrVersion: "V0.6",
		ITcpPort: 8999,
		IMaxConn: 1000,
		IMaxPackageSize: 4096,
	}

	// 应该尝试从 conf/zinx.json 去加载用户自定义的参数
	GlobalObject.Reload()
}



