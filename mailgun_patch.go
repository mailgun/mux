package mux

import "reflect"

func (r *Router) NewUnattachedRoute() *Route {
	route := &Route{routeConf: copyRouteConf(r.routeConf), namedRoutes: r.namedRoutes}
	return route
}

func (r *Router) UpsertRoute(newRoute *Route) {
	for i, currRoute := range r.routes {
		if reflect.DeepEqual(currRoute.matchers, newRoute.matchers) {
			r.routes[i] = newRoute
			return
		}
	}
	r.routes = append(r.routes, newRoute)
}
