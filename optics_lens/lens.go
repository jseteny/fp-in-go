package optics_lens

import "fmt"

type Lens struct {
	Get
	Set
}
type Get func(obj interface{}) interface{}
type Set func(obj, value interface{}) interface{}

func Compose(outer, inner Lens) Lens {
	return Lens{
		Get: func(a interface{}) interface{} {
			return inner.Get(outer.Get(a))
		},
		Set: func(obj, value interface{}) interface{} {
			return outer.Set(obj, inner.Set(outer.Get(obj), value))
		},
	}
}

func Example_lens() {
	
	type Street struct {
		name   string
		number int
	}

	type Address struct {
		country string
		city    string
		street  Street
	}


	street := Street{"Mura utca", 8}
	address := Address{"Magyarország", "Szeged", street}
	_streetNumberLens := Lens{
		func(s interface{}) interface{} {
			return s.(Street).number
		},
		func(street, number interface{}) interface{} {
			return Street{street.(Street).name, number.(int)}
		},
	}
	streetNumberLens := struct { // Type Safe Lens
		l   Lens
		Get func(s Street) int
		Set func(s Street, number int) Street
	}{
		l: _streetNumberLens,
		Get: func(s Street) int {
			return _streetNumberLens.Get(s).(int)
		},
		Set: func(s Street, number int) Street {
			return _streetNumberLens.Set(s, number).(Street)
		},
	}
	_addressStreetLens := Lens{
		func(address interface{}) interface{} {
			return address.(Address).street
		},
		func(address, street interface{}) interface{} {
			return Address{address.(Address).country, address.(Address).city, street.(Street)}
		},
	}
	fmt.Println(streetNumberLens.Get(street))
	fmt.Println(_addressStreetLens.Get(address))

	_addressStreetNumberLens := Compose(_addressStreetLens, _streetNumberLens)
	fmt.Println(_addressStreetNumberLens.Set(address, 77))
}
