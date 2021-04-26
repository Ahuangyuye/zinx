package znet

import "zinx/src/zinx/ziface"

// 实现 router 时， 先嵌入这个BaseRouter 基类，然后根据需要
// 对这个基类的方法进行重写就好了
type BaseRouter struct {
	
}

// 这里之所以 BaseRouter 的方法都为空
// 是因为有的roter 不希望有 PreHandle PostHandler 这两个业务
//  所以 Roter 全部继承 BaseRouter 好处就是 不需要实现  PreHandle PostHandler
// 在处理 conn 业务之前的钩子方法 hook
func (pBR *BaseRouter )PreHandle(request ziface.IRequest){

}

// 在处理 conn  业务的主方法 hook
func (pBR *BaseRouter ) Handler(request  ziface.IRequest){

}

// 在处理 conn 业务之后的钩子方法 hook
func (pBR *BaseRouter ) PostHandler(request  ziface.IRequest){

}











