package grpc

import "google.golang.org/grpc"

func RegisterOption(grpcOpts []grpc.ServerOption) SetOption {
	return func(o *grpcServerOption) {
		if grpcOpts == nil {
			panic("grpc.RegisterService: no opts in param")
		}

		o.grpcOpts = grpcOpts
	}
}

func RegisterService(srvSet SetService) SetOption {
	return func(o *grpcServerOption) {
		if srvSet == nil {
			panic("grpc.RegisterService: no opts in param")
		}

		o.srvSet = srvSet
	}
}
