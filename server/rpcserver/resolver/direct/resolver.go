package direct

/**
* created by mengqi on 2023/11/21
* 直连resolver
 */

import "google.golang.org/grpc/resolver"

type directResolver struct{}

func newDirectResolver() *directResolver {
	return &directResolver{}
}

func (r *directResolver) Close() {}

func (r *directResolver) ResolveNow(options resolver.ResolveNowOptions) {}
