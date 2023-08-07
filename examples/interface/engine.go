package main

import "github.com/samber/do"

type Engine interface{}

type engineImplem struct {
}

// NewEngine 创建一个引擎， 该函数的返回定义为 Engine 接口， 但是实际返回的是 engineImplem 实例
func NewEngine(i *do.Injector) (Engine, error) {
	return &engineImplem{}, nil
}
