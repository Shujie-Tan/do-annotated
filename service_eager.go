package do

// ServiceEager 和 ServiceLazy 的区别在于， ServiceEager 在注册时就会创建实例， 而 ServiceLazy 在获取实例时才会创建实例, 所以ServiceEager不需要provider， 也不需要加锁
type ServiceEager[T any] struct {
	name     string
	instance T
}

func newServiceEager[T any](name string, instance T) Service[T] {
	return &ServiceEager[T]{
		name:     name,
		instance: instance,
	}
}

//nolint:unused
func (s *ServiceEager[T]) getName() string {
	return s.name
}

//nolint:unused
func (s *ServiceEager[T]) getInstance(i *Injector) (T, error) {
	return s.instance, nil
}

func (s *ServiceEager[T]) healthcheck() error {
	instance, ok := any(s.instance).(Healthcheckable)
	if ok {
		return instance.HealthCheck()
	}

	return nil
}

func (s *ServiceEager[T]) shutdown() error {
	instance, ok := any(s.instance).(Shutdownable)
	if ok {
		return instance.Shutdown()
	}

	return nil
}

func (s *ServiceEager[T]) clone() any {
	return s
}
