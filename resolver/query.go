package resolver


//Query
func (_ *Resolver) Riff() *RiffResolver {
	return &RiffResolver{}
}

func (_ *Resolver) Nodes() *[]*NodeResolver {
	return &[]*NodeResolver{{}}
}
