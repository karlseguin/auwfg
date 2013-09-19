package auwfg

import(
  "strings"
  "net/http"
  // "users/server/middleware"
)

// var middlewareRunner = func(context *web.Context, route *web.Route) web.Response { return middleware.Run(context, route) }

type Router struct {
  *Configuration
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
  route, params := r.loadRouteAndParams(req)
  if route == nil {
    reply(writer, r.notFound)
    return
  }
  bc := &BaseContext{Route: route, Params: params, Req: req,}
  context := r.contextFactory(bc)
  // reply(writer, middlewareRunner(context, route))
  reply(writer, r.dispatcher(route, context))
}

func reply(writer http.ResponseWriter, res Response) {
  h := writer.Header()
  for k, v := range res.Header() { h[k] = v }
  writer.WriteHeader(res.Status())
  writer.Write(res.Body())
}

func (r *Router) loadRouteAndParams(req *http.Request) (*Route, *Params) {
  path := req.URL.Path
  if len(path) < 4 { return nil, nil }

  end := strings.LastIndex(path, ".")
  if end == -1 { end = len(path) }

  parts := strings.Split(path[1:end], "/")
  l := len(parts)
  if l < 2 || l > 5 { return nil, nil }

  for index, part := range parts {
    parts[index] = strings.ToLower(part)
  }

  version, exists := r.routes[parts[0]]
  if exists == false { return nil, nil }

  params := loadParams(parts[1:])
  controller, exists := version[params.Resource]
  if exists == false { return nil, nil }

  m := req.Method
  if m == "GET" && len(params.Id) == 0 { m = "LIST" }

  route, exists := controller[m]
  if exists == false { return nil, nil }

  params.Version = parts[0]
  return route, params
}

func loadParams(parts []string) *Params {
  params := new(Params)
  switch len(parts) {
  case 1:
    params.Resource = parts[0]
  case 2:
    params.Resource = parts[0]
    params.Id = parts[1]
  case 3:
    params.ParentResource = parts[0]
    params.ParentId = parts[1]
    params.Resource = parts[2]
  case 4:
    params.ParentResource = parts[0]
    params.ParentId = parts[1]
    params.Resource = parts[2]
    params.Id = parts[3]
  }
  return params
}
