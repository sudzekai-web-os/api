package objects

type Flag struct {
	name         string
	bindProperty string
	usage        string

	valueGetter func() any
	visited     bool
}

func NewFlag(
	name string,
	bindProperty string,
	value any,
	usage string,
) Flag {
	return Flag{
		name:         name,
		bindProperty: bindProperty,
		usage:        usage,
		valueGetter: func() any {
			return value
		},
	}
}

func (f *Flag) GetName() string {
	return f.name
}

func (f *Flag) GetBindProperty() string {
	return f.bindProperty
}

func (f *Flag) GetValue() any {
	if f.valueGetter == nil {
		return nil
	}

	return f.valueGetter()
}

func (f *Flag) IsVisited() bool {
	return f.visited
}

func (f *Flag) SetIsVisited(value bool) {
	f.visited = value
}

func (f *Flag) SetValueGetter(getter func() any) {
	f.valueGetter = getter
}

func (f *Flag) GetUsage() string {
	return f.usage
}
