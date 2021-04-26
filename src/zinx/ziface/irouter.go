package ziface


/*
 路由抽象接口
 路由里的数据都是 IRequest
*/

type IRouter interface {

	// 在处理 conn 业务之前的钩子方法 hook
	PreHandle(request IRequest)

	// 在处理 conn  业务的主方法 hook
	Handler(request IRequest)

	// 在处理 conn 业务之后的钩子方法 hook
	PostHandler(request IRequest)

}









