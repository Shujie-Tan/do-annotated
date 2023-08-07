package do

import (
	"fmt"
)

type Service[T any] interface {
	getName() string
	getInstance(*Injector) (T, error)
	healthcheck() error
	shutdown() error
	clone() any
}

type healthcheckableService interface {
	healthcheck() error
}

type shutdownableService interface {
	shutdown() error
}

// generateServiceName 生成服务名称, 如果是结构体， 则返回结构体类型名称， 如果是接口， 则返回接口类型名称
func generateServiceName[T any]() string {
	var t T

	// struct
	name := fmt.Sprintf("%T", t)
	if name != "<nil>" {
		return name
	}

	// interface
	return fmt.Sprintf("%T", new(T))
}

type Healthcheckable interface {
	HealthCheck() error
}

type Shutdownable interface {
	Shutdown() error
}

type cloneableService interface {
	clone() any
}
