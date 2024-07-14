package httpi

type (
	Option interface {
		apply(*options)
	}
	options struct {
		roundTripperFunc      RoundTripperFunc
		requestValidationFunc RequestValidationFunc
	}

	roundTripperFuncOption      RoundTripperFunc
	requestValidationFuncOption RequestValidationFunc
)

// WithRoundTripperFunc sets the RoundTripperFunc to be used by the Interceptor.
func WithRoundTripperFunc(f RoundTripperFunc) Option {
	return roundTripperFuncOption(f)
}

// WithRequestValidationFunc sets the RequestValidationFunc to be used by the Interceptor.
func WithRequestValidationFunc(f RequestValidationFunc) Option {
	return requestValidationFuncOption(f)
}

func (rtfo roundTripperFuncOption) apply(opts *options) {
	opts.roundTripperFunc = RoundTripperFunc(rtfo)
}

func (rvfo requestValidationFuncOption) apply(opts *options) {
	opts.requestValidationFunc = RequestValidationFunc(rvfo)
}
