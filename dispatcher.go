package auwfg

type Dispatcher func(route *Route, context interface{}) Response

func genericDispatcher(route *Route, context interface{}) Response {
  return route.Action.(func(interface{}) Response)(context)
}
