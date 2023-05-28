package starter

type StarterContext map[string]interface{}

type Starter interface {
	// 1. 系统启动，初始化一些基础资源，比如数据库，消息队列等
	Init(ctx StarterContext)
	// 2. 系统基础资源的安装，比如初始化数据库表，或者其他基础资源的初始化
	Setup(ctx StarterContext)
	// 3. 启动基础资源，比如启动数据库连接池，消息队列消费者
	Start(ctx StarterContext)
	// 启动器是否可阻塞，比如启动RPC服务，启动器需要一直阻塞才能提供RPC服务
	StartBlocking() bool
	// 4. 系统停止，一般是一些基础资源的回收和停止，比如数据库连接池，消息队列连接池等
	Stop(ctx StarterContext)
}

type BaseStarter struct{}

func (b *BaseStarter) Init(ctx StarterContext)  {}
func (b *BaseStarter) Setup(ctx StarterContext) {}
func (b *BaseStarter) Start(ctx StarterContext) {}
func (b *BaseStarter) StartBlocking() bool      { return false }
func (b *BaseStarter) Stop(ctx StarterContext)  {}

type StarterRegister struct {
	starters []Starter
}

func (sr *StarterRegister) Register(s Starter) {
	sr.starters = append(sr.starters, s)
}

func (sr *StarterRegister) AllStarters() []Starter {
	return sr.starters
}

var starterRegister *StarterRegister = new(StarterRegister)

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
