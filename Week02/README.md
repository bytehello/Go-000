go error 学习笔记

原则：

错误只处理一次(handle once)

用法：
1. 应用代码中，用pkg/error errors.New errors.Errorf 保存堆栈数据
2. 调用包内其他函数，直接简单返回，避免因为wrap 多倍的堆栈信息
3. 调用其他第三方库，调用基础库标准，自己的库，应该是wrap保存根因
4. 直接返回操作，不用导出打日志
5. 工作的goroutine顶部使用 %+v打印

总结：

1. 作为kit/标准库（被很多人使用）、可重用性高的包只能返回根错误，eg sql/database就是如此，肯定不会wrap给你的，
2. 不打算处理的包/不降级，wrap起来，往上抛 ，举例：调用dao层报错，dao层作为与db库交互的最后一层， 会 wrap 错误上抛