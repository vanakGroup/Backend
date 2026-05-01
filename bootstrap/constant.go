package bootstrap

type Constants struct {
	Context Context
	Role    Role
}

func NewConstant() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
		},
		Role: Role{
			DefaultType: "customer",
		},
	}
}

type Context struct {
	Translator string
}

type Role struct {
	DefaultType string
}
