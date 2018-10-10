package resolver

type NodeResolver struct {
}
func (r *NodeResolver) DataCenter() string {
	return ""
}
func (r *NodeResolver) Ip() string {
	return ""
}
func (r *NodeResolver) IsSelf() bool {
	return false
}
func (r *NodeResolver) Name() string {
	return ""
}
func (r *NodeResolver) Port() int32 {
	return 0
}
func (r *NodeResolver) Services() *[]*NestServiceResover {
	return &[]*NestServiceResover{{}}
}
func (r *NodeResolver) SnapShot() string {
	return ""
}
func (r *NodeResolver) State() string {
	return ""
}
func (r *NodeResolver) Version() int32 {
	return 0
}