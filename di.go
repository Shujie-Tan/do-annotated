package do

import (
	"fmt"
)

// Provide: 注册一个服务， 该服务的实例是通过 provider 函数来创建的
// Provider[T] 是一个函数类型， 该函数接收一个 *Injector 指针， 返回一个 T 类型的实例和一个 error
func Provide[T any](i *Injector, provider Provider[T]) {
	name := generateServiceName[T]()

	_i := getInjectorOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	service := newServiceLazy(name, provider)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

// ProvideNamed 注册命名的服务, 而不是使用默认的服务名称， 这样可以注册多个同类型的服务, 例如注册多个数据库连接
func ProvideNamed[T any](i *Injector, name string, provider Provider[T]) {
	_i := getInjectorOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	service := newServiceLazy(name, provider)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

// ProvideValue 直接注册一个实例， 而不是指定一个 provider 函数来创建实例
func ProvideValue[T any](i *Injector, value T) {
	name := generateServiceName[T]()

	_i := getInjectorOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

// ProvideNamedValue 注册命名的服务实例， 同时指定服务名称和服务的实例
func ProvideNamedValue[T any](i *Injector, name string, value T) {
	_i := getInjectorOrDefault(i)
	if _i.exists(name) {
		panic(fmt.Errorf("DI: service `%s` has already been declared", name))
	}

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s injected", name)
}

func Override[T any](i *Injector, provider Provider[T]) {
	name := generateServiceName[T]()

	_i := getInjectorOrDefault(i)

	service := newServiceLazy(name, provider)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

func OverrideNamed[T any](i *Injector, name string, provider Provider[T]) {
	_i := getInjectorOrDefault(i)

	service := newServiceLazy(name, provider)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

func OverrideValue[T any](i *Injector, value T) {
	name := generateServiceName[T]()

	_i := getInjectorOrDefault(i)

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

func OverrideNamedValue[T any](i *Injector, name string, value T) {
	_i := getInjectorOrDefault(i)

	service := newServiceEager(name, value)
	_i.set(name, service)

	_i.logf("service %s overridden", name)
}

// Invoke: 获取一个服务的实例
func Invoke[T any](i *Injector) (T, error) {
	// 生成服务名称, 用于不需要命名的服务
	name := generateServiceName[T]()
	return InvokeNamed[T](i, name)
}

// MustInvoke 确保invoke成功， 如果失败则panic
func MustInvoke[T any](i *Injector) T {
	s, err := Invoke[T](i)
	must(err)
	return s
}

func InvokeNamed[T any](i *Injector, name string) (T, error) {
	return invokeImplem[T](i, name)
}

func MustInvokeNamed[T any](i *Injector, name string) T {
	s, err := InvokeNamed[T](i, name)
	must(err)
	return s
}

func invokeImplem[T any](i *Injector, name string) (T, error) {
	_i := getInjectorOrDefault(i)

	// 检查服务是否已注册
	serviceAny, ok := _i.get(name)
	if !ok {
		return empty[T](), _i.serviceNotFound(name)
	}

	// 尝试将服务转换为 Service[T] 类型, Service[T] 是一个接口类型, 该接口有两种实现: ServiceEager[T] 和 ServiceLazy[T]
	service, ok := serviceAny.(Service[T])
	if !ok {
		return empty[T](), _i.serviceNotFound(name)
	}

	// 调用 Service[T].getInstance 方法获取服务实例， 如果是 ServiceEager[T] 类型， 则直接返回实例， 如果是 ServiceLazy[T] 类型， 则调用 ServiceLazy[T].build 方法创建实例
	instance, err := service.getInstance(_i)
	if err != nil {
		return empty[T](), err
	}

	// 调用注册的回调函数
	_i.onServiceInvoke(name)

	_i.logf("service %s invoked", name)

	return instance, nil
}

func HealthCheck[T any](i *Injector) error {
	name := generateServiceName[T]()
	return getInjectorOrDefault(i).healthcheckImplem(name)
}

func HealthCheckNamed(i *Injector, name string) error {
	return getInjectorOrDefault(i).healthcheckImplem(name)
}

func Shutdown[T any](i *Injector) error {
	name := generateServiceName[T]()
	return getInjectorOrDefault(i).shutdownImplem(name)
}

func MustShutdown[T any](i *Injector) {
	name := generateServiceName[T]()
	must(getInjectorOrDefault(i).shutdownImplem(name))
}

func ShutdownNamed(i *Injector, name string) error {
	return getInjectorOrDefault(i).shutdownImplem(name)
}

func MustShutdownNamed(i *Injector, name string) {
	must(getInjectorOrDefault(i).shutdownImplem(name))
}
