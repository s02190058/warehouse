package migrations

type Option func(o *options)

func WithSteps(n int) Option {
	return func(o *options) {
		o.steps = n
	}
}
