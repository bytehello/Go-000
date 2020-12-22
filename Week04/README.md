## 作业题目
> 1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架

参考 https://github.com/go-kratos/kratos 自动生成

## 笔记
/cmd/myapp 只放很少的代码

/internal 只希望自己使用 internal/app/myapp 用来标识意义；去看标准库有很多internal，go编译器实现的

/internal/pkg 整个项目中，可以被整个git项目中的其他使用，

gitlab 不加 internal 是可以被别人import的

/pkg 外部可以导入的