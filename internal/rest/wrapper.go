package rest

import "github.com/gin-gonic/gin"

// EndpointWrapper is a wrapper for all routes in the API
// It will be used to register global middlewares, like authentication
// and authorization or interceptors
type EndpointWrapper struct {
	globalInterceptors     []gin.HandlerFunc
	endpointInterceptorMap map[string][]gin.HandlerFunc
}

// NewGlobalWrapper creates a new EndpointWrapper
func newEndpointWrapper() *EndpointWrapper {
	return &EndpointWrapper{
		globalInterceptors:     []gin.HandlerFunc{},
		endpointInterceptorMap: map[string][]gin.HandlerFunc{},
	}
}

// InitEndpointWrapper initializes the global wrapper
func InitEndpointWrapper(globalInterceptors []gin.HandlerFunc, endpointInterceptorMap map[string][]gin.HandlerFunc) *EndpointWrapper {
	globalWrapper := newEndpointWrapper()
	globalWrapper.globalInterceptors = globalInterceptors
	globalWrapper.endpointInterceptorMap = endpointInterceptorMap

	return globalWrapper
}

// WrapEngine wraps a gin engine with the global wrapper
func (wrapper *EndpointWrapper) WrapEngine(engine *gin.Engine) {
	engine.Use(wrapper.globalInterceptors...)
	for endpoint, interceptors := range wrapper.endpointInterceptorMap {
		group := engine.Group(endpoint)
		group.Use(interceptors...)
	}
}

// entity -> request body, headers
// restTemplate.send(POST, entity, url, response)
