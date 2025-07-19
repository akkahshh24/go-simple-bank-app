package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// Header used by grpc-gateway to pass the user agent
	// This is sent by a HTTP client that uses grpc-gateway to communicate with the gRPC server.
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"

	// Header used by the client to pass the user agent
	// This is sent by a gRPC client directly communicating with the gRPC server.
	userAgentHeader = "user-agent"

	// Header used to pass the client IP address
	xForwardedForHeader = "x-forwarded-for"
)

// Metadata holds the metadata extracted from the gRPC context.
type Metadata struct {
	UserAgent string
	ClientIP  string
}

// extractMetadata extracts metadata from the gRPC context.
// It retrieves the user agent and client IP from the metadata headers.
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	// Extract metadata from the context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// Check for user agent in grpc-gateway header or user-agent header
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// Check for client IP in x-forwarded-for header
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	// If the client IP is not set, try to get it from the peer information
	// This is used when the gRPC server is directly accessed by a gRPC client.
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
