# 目录

```bash
|- main.go # 入口
|- starter
  |- config.go # 资源配置文件
  |- starter.go # 启动器的统一入口
```

# 笔记

## 统一初始化入口

什么时候需要统一初始化入口函数呢？

当我们需要在 `main` 函数中做一些初始化的工作，比如初始化日志，初始化配置文件，初始化数据库等等

这些都是需要在 `main` 函数中做的初始化工作，这个时候就可以使用统一初始化入口的方式来实现

## 抽象出生命周期接口

使用统一初始化入口的方式，需要我们抽象出生命周期接口

比如下面我们定义了 `Init`、`Setup`、`Start`、`StartBlocking`、`Stop` 五个接口，用来完成五件事情：

- `Init`: 系统启动，初始化一些基础资源，比如数据库，消息队列等
- `Setup`: 安装一些基础资源，比如初始化数据库表，或者其他基础资源的初始化
- `Start`: 启动基础资源，比如启动数据库连接池，消息队列消费者
- `StartBlocking`: 启动器是否可阻塞，比如启动 `RPC` 服务，启动器需要一直阻塞才能提供 `RPC` 服务
- `Stop`: 系统停止，一般是一些基础资源的回收和停止，比如数据库连接池，消息队列连接池等

```go
type Starter interface {
	Init(ctx StarterContext)
	Setup(ctx StarterContext)
	Start(ctx StarterContext)
	StartBlocking() bool
	Stop(ctx StarterContext)
}
```

## 实现生命周期接口

下一步是实现这个接口，在 go 中 `interface` 是隐式实现的，只要实现了接口中的方法，就是实现了这个接口

```go
type BaseStarter struct{}

func (b *BaseStarter) Init(ctx StarterContext)  {}
func (b *BaseStarter) Setup(ctx StarterContext) {}
func (b *BaseStarter) Start(ctx StarterContext) {}
func (b *BaseStarter) StartBlocking() bool      { return false }
func (b *BaseStarter) Stop(ctx StarterContext)  {}
```

## 实现自动装配

所有的启动器都实现了 `Starter` 接口

这时就可以将所有的启动器都放在一个切片中，然后遍历切片，依次调用 `Init`、`Setup`、`Start`、`StartBlocking`、`Stop` 方法

```go
type StarterRegister struct {
	starters []Starter
}
func (sr *StarterRegister) Register(s Starter) {
	sr.starters = append(sr.starters, s)
}
func (sr *StarterRegister) AllStarters() []Starter {
	return sr.starters
}
```

提供一个注册和运行的方法，这样就可以实现自动装配了

```go
func Register(s Starter) {
	starterRegister.Register(s)
}
func SystemRun() {
	ctx := StarterContext{}
	for _, s := range starterRegister.AllStarters() {
		s.Init(ctx)
		s.Setup(ctx)
		s.Start(ctx)
		s.StartBlocking()
		s.Stop(ctx)
	}
}
```

## 使用

1. 为每个资源创建一个 `config.go` 文件，然后在 `init` 函数中注册

```go
func init() {
	fmt.Println("config")
	Register(&Config{})
}

type Config struct{}

func (c *Config) Init(ctx StarterContext) {
	fmt.Println("Init")
}

func (c *Config) Setup(ctx StarterContext) {
	fmt.Println("Setup")
}

func (c *Config) Start(ctx StarterContext) {
	fmt.Println("Start")
}

func (c *Config) StartBlocking() bool {
  fmt.Println("StartBlocking")
	return true
}

func (c *Config) Stop(ctx StarterContext) {
	fmt.Println("Stop")
}
```

2. 在 `main` 函数中调用 `SystemRun` 方法

通过这种方式，就可以实现自动装配了，有多少资源，就创建多少 `config.go` 文件，并实现 `Starter` 接口
