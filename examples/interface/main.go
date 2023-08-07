package main

import (
	"github.com/samber/do"
)

func main() {
	injector := do.New()

	// ProvideNamedValue: 注册轮子实例， 通过该函数可以为某个类型提供一个命名的实例，这样在后续的调用中，可以通过该名称来获取该实例
	do.ProvideNamedValue(injector, "wheel-1", NewWheel())
	do.ProvideNamedValue(injector, "wheel-2", NewWheel())
	do.ProvideNamedValue(injector, "wheel-3", NewWheel())
	do.ProvideNamedValue(injector, "wheel-4", NewWheel())

	// provide car: 注意这里虽然Car依赖Engine， 但是可以在注册Engine之前注册Car， 因为do.Injector会自动解析依赖
	do.Provide(injector, NewCar)

	// provide engine
	do.Provide(injector, NewEngine)

	// MustInvoke: 通过 do.Injector 获取 Car 实例
	car := do.MustInvoke[Car](injector)
	car.Start()
}
