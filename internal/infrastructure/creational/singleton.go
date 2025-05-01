package creational

import (
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"sync"
)

type Supplier func() (interface{}, error)

type Singleton[T any] struct {
	instance T
	once     sync.Once
	supplier Supplier
	err      error
}

func NewSingleton[T any](supplier Supplier) *Singleton[T] {
	return &Singleton[T]{
		supplier: supplier,
	}
}

func (s *Singleton[T]) Get() (T, error) {
	s.once.Do(func() {
		instance, err := s.supplier()
		if err != nil {
			s.err = err
			return
		}
		_instances, ok := instance.(T)
		if !ok {
			s.err = errors.New(errors.ErrnoCastError)
			return
		}
		s.instance = _instances
	})
	return s.instance, s.err
}

func (s *Singleton[T]) Reset() {
	s.once = sync.Once{}
	s.instance = *new(T)
	s.err = nil
}

type SingletonFactory struct {
	singletons map[string]interface{}
	lock       sync.RWMutex
}

func NewSingletonFactory() *SingletonFactory {
	return &SingletonFactory{
		singletons: make(map[string]interface{}),
	}
}

func (f *SingletonFactory) Register(typeName string, supplier Supplier) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if _, exists := f.singletons[typeName]; !exists {
		f.singletons[typeName] = NewSingleton[interface{}](supplier)
	}
}

func (f *SingletonFactory) Get(typeName string) (interface{}, error) {
	f.lock.RLock()
	singleton, exists := f.singletons[typeName]
	f.lock.RUnlock()
	if !exists {
		return nil, errors.New(errors.ErrnoResourceNotFoundException)
	}
	s, ok := singleton.(*Singleton[interface{}])
	if !ok {
		return nil, errors.New(errors.ErrnoCastError)
	}
	return s.Get()
}

func (f *SingletonFactory) Clear(typeName string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if s, exists := f.singletons[typeName]; exists {
		if _s, ok := s.(*Singleton[interface{}]); ok {
			_s.Reset()
		}
		delete(f.singletons, typeName)
	}
}
