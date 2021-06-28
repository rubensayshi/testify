package assert

import (
	"reflect"
	"sync"
)

const equalMethodName = "Equal"

var equallerCacheMu sync.RWMutex
var equallerCache map[reflect.Type]bool

func init() {
	equallerCache = make(map[reflect.Type]bool, 0)
}

func checkIsEqualByEqualler(expected, actual interface{}) bool {
	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	if !isEqualler(actualType) {
		return false
	}

	equalMethod, _ := actualType.MethodByName(equalMethodName)
	expectedValue := reflect.ValueOf(expected)

	// if the equal method takes our expected value as arg then we can just call it
	if equalMethod.Type.In(1) == expectedValue.Type() {
		res := reflect.ValueOf(actual).MethodByName(equalMethodName).Call([]reflect.Value{expectedValue})
		return res[0].Bool()
	}

	// if the equal expects a dereferenced argument then we can dereference expectedValue and call it
	if expectedValue.Type().Kind() == reflect.Ptr && equalMethod.Type.In(1) == expectedValue.Type().Elem() {
		expectedValue = expectedValue.Elem()
		res := reflect.ValueOf(actual).MethodByName(equalMethodName).Call([]reflect.Value{expectedValue})
		return res[0].Bool()
	}

	return false
}

func isEqualler(t reflect.Type) bool {
	isEqualler, cached := isEquallerCached(t)
	if !cached {
		isEqualler = determineIsEqualler(t)
		setIsEquallerCached(t, isEqualler)
	}

	return isEqualler
}

func determineIsEqualler(t reflect.Type) bool {
	equalMethod, hasEqualMethod := t.MethodByName(equalMethodName)
	if hasEqualMethod {
		// should have only 1 return value which should be a bool
		if equalMethod.Type.NumOut() != 1 || equalMethod.Type.Out(0).Kind() != reflect.Bool {
			return false
		}
		// and should have exactly 2 arguments (pointer method so first is self)
		if equalMethod.Type.NumIn() != 2 {
			return false
		}

		// if the receiver and the argument match then we can use it
		if equalMethod.Type.In(0) == equalMethod.Type.In(1) {
			return true
		}

		// or if it's a pointer receiver and the argument isn't a pointer
		// then we check if the dereferenced receiver matches the argument
		if equalMethod.Type.In(0).Kind() == reflect.Ptr && equalMethod.Type.In(1).Kind() != reflect.Ptr &&
			equalMethod.Type.In(0).Elem() == equalMethod.Type.In(1) {
			return true
		}

		return false
	}

	return false
}

func isEquallerCached(t reflect.Type) (bool, bool) {
	equallerCacheMu.RLock()
	defer equallerCacheMu.RUnlock()

	isEqualler, cached := equallerCache[t]

	return isEqualler, cached
}

func setIsEquallerCached(t reflect.Type, isEqualler bool) {
	equallerCacheMu.Lock()
	defer equallerCacheMu.Unlock()

	equallerCache[t] = isEqualler
}
