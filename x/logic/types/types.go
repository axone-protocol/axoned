package types

import "google.golang.org/grpc"

// GrpcQueryServiceDesc represents the query server's RPC service specification.
// This gives access to the service name and method names needed for stargate
// queries.
func GrpcQueryServiceDesc() grpc.ServiceDesc {
	return _QueryService_serviceDesc
}
