package assert

import (
	"reflect"
	"testing"
)

type Foo struct {
	id int
}

func (f Foo) Equal(other Foo) bool {
	return other.id == f.id
}

type Bar struct {
	id int
}

func (f Bar) Equal(other Bar) bool {
	return other.id == f.id
}

func TestCheckIsEqualByEqualler(t *testing.T) {
	f1a := Foo{1}
	f1b := Foo{1}
	f2 := Foo{2}
	f1ptr := &Foo{1}

	if !checkIsEqualByEqualler(f1a, f1a) {
		t.Fail()
	}
	if !checkIsEqualByEqualler(f1a, f1b) {
		t.Fail()
	}
	if checkIsEqualByEqualler(f1a, f2) {
		t.Fail()
	}

	if !checkIsEqualByEqualler(f1ptr, f1ptr) {
		t.Fail()
	}
	if !checkIsEqualByEqualler(f1a, f1ptr) {
		t.Fail()
	}
	if !checkIsEqualByEqualler(f1ptr, f1a) {
		t.Fail()
	}
	if checkIsEqualByEqualler(f1ptr, f2) {
		t.Fail()
	}
	if checkIsEqualByEqualler(f2, f1ptr) {
		t.Fail()
	}
}

func TestCheckIsEqualByEquallerNoEq(t *testing.T) {
	foo := FooNoEq{1}

	if checkIsEqualByEqualler(foo, foo) {
		t.Fail()
	}
}

func TestCheckIsEqualByEquallerNotTheSameType(t *testing.T) {
	foo := Foo{1}
	bar := Bar{1}

	if checkIsEqualByEqualler(foo, bar) {
		t.Fail()
	}
}

func TestDetermineIsEqualler(t *testing.T) {
	fooT := reflect.TypeOf(Foo{})
	fooTIsEqualler := determineIsEqualler(fooT)
	if !fooTIsEqualler {
		t.Errorf("Foo should be isEqualler")
	}

	fooPtrT := reflect.TypeOf(&Foo{})
	fooPtrTIsEqualler := determineIsEqualler(fooPtrT)
	if !fooPtrTIsEqualler {
		t.Errorf("*Foo should be isEqualler")
	}
}

type FooNoEq struct {
	id int
}

func TestDetermineIsEquallerNoMethod(t *testing.T) {
	fooT := reflect.TypeOf(FooNoEq{})
	fooTIsEqualler := determineIsEqualler(fooT)
	if fooTIsEqualler {
		t.Errorf("FooNoEq doesn't have Equal method, shouldn't be isEqualler")
	}
}

type FooFunkyEqReturn struct {
	id int
}

func (f FooFunkyEqReturn) Equal(other FooFunkyEqReturn) (bool, bool) {
	return other.id == f.id, true
}

func TestDetermineIsEquallerFunkyEqReturn(t *testing.T) {
	fooT := reflect.TypeOf(FooFunkyEqReturn{})
	fooTIsEqualler := determineIsEqualler(fooT)
	if fooTIsEqualler {
		t.Errorf("FooFunkyEqReturn has a weird return value for Equal, shouldn't be isEqualler")
	}
}

type FooFunkyEqArg struct {
	id int
}

func (f FooFunkyEqArg) Equal(other Foo) (bool, bool) {
	return other.id == f.id, true
}

func TestDetermineIsEquallerFunkyEqArg(t *testing.T) {
	fooT := reflect.TypeOf(FooFunkyEqArg{})
	fooTIsEqualler := determineIsEqualler(fooT)
	if fooTIsEqualler {
		t.Errorf("FooFunkyEqArg has a weird argument value for Equal, shouldn't be isEqualler")
	}
}

func TestIsEquallerCache(t *testing.T) {
	fooT := reflect.TypeOf(Foo{})
	fooPtrT := reflect.TypeOf(&Foo{})

	// reset cache
	equallerCache = make(map[reflect.Type]bool, 0)

	if _, isCached := isEquallerCached(fooT); isCached {
		t.Errorf("Foo shouldn't be cached yet")
	}
	if _, isCached := isEquallerCached(fooPtrT); isCached {
		t.Errorf("*Foo shouldn't be cached yet")
	}

	setIsEquallerCached(fooT, true)

	if isEqualler, isCached := isEquallerCached(fooT); !isCached && !isEqualler {
		t.Errorf("Foo should be cached and true")
	}
	if _, isCached := isEquallerCached(fooPtrT); isCached {
		t.Errorf("*Foo shouldn't be cached yet")
	}

	setIsEquallerCached(fooPtrT, false)

	if isEqualler, isCached := isEquallerCached(fooT); !isCached && !isEqualler {
		t.Errorf("Foo should be cached and true")
	}
	if isEqualler, isCached := isEquallerCached(fooPtrT); isCached && isEqualler {
		t.Errorf("*Foo should be cached and false")
	}

}

// the tests for determineIsEqualler should cover most cases, here we just test we are using the cache
func TestIsEqualler(t *testing.T) {
	fooT := reflect.TypeOf(Foo{})
	fooPtrT := reflect.TypeOf(&Foo{})

	// reset cache
	equallerCache = make(map[reflect.Type]bool, 0)

	setIsEquallerCached(fooT, true)

	if !isEqualler(fooT) {
		t.Errorf("Foo should be cached and true")
	}

	setIsEquallerCached(fooPtrT, true)

	if !isEqualler(fooPtrT) {
		t.Errorf("*Foo should be cached and true")
	}

	setIsEquallerCached(fooPtrT, false)

	if !isEqualler(fooT) {
		t.Errorf("Foo should be cached and true")
	}
	if isEqualler(fooPtrT) {
		t.Errorf("*Foo should be cached and false")
	}
}
