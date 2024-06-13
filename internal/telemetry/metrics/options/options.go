package options

type (
	serverMetricsOptions struct {
		Namespace string
		Subsystem string
	}

	Options interface {
		apply(*serverMetricsOptions)
	}

	serverMetricsOptionApplier func(*serverMetricsOptions)
)

var _ Options = (serverMetricsOptionApplier)(nil)

const (
	DefaultNamespace = "sys"
	DefaultSubsystem = "grpc_server"
)

func (applier serverMetricsOptionApplier) apply(opts *serverMetricsOptions) {
	applier(opts)
}

func WithNamespace(ns string) Options {
	var applier serverMetricsOptionApplier = func(opts *serverMetricsOptions) {
		opts.Namespace = ns
	}

	return applier
}

func WithSubsystem(ss string) Options {
	var applier serverMetricsOptionApplier = func(opts *serverMetricsOptions) {
		opts.Subsystem = ss
	}

	return applier
}
