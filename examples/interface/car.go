package main

import "github.com/samber/do"

type Car interface {
	// Start 汽车有启动的功能
	Start()
}

// carImplem 实现了 Car 接口
type carImplem struct {
	// Engine 引擎， 汽车依赖引擎
	Engine Engine
	// Wheels 轮子， 汽车依赖轮子（一共四个）
	Wheels []*Wheel
}

// Start 汽车启动
func (c *carImplem) Start() {
	println("vroooom")
}

// NewCar 创建汽车实例， 所有的汽车实例对象都从注册器 do.Injector 中获取， 所以输入参数是 do.Injector
func NewCar(i *do.Injector) (Car, error) {
	wheels := []*Wheel{
		// MustInvokeNamed: 通过轮子名称获取轮子实例， 这里多次调用NewCar 每次获取的都是同一个实例
		do.MustInvokeNamed[*Wheel](i, "wheel-1"),
		do.MustInvokeNamed[*Wheel](i, "wheel-2"),
		do.MustInvokeNamed[*Wheel](i, "wheel-3"),
		do.MustInvokeNamed[*Wheel](i, "wheel-4"),
	}

	// MustInvoke: 通过 do.Injector 获取 Engine 实例
	engine := do.MustInvoke[Engine](i)

	car := carImplem{
		Engine: engine,
		Wheels: wheels,
	}

	return &car, nil
}
