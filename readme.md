# Another useless web framework for Go
AUWFG is an ugly web stack for Go used for building web services.

AUWFG has two goals.

First, to provide a generic framework while still levarging strong typing. Many Go-based frameworks are built as a black box resulting in the inability to be cleanly expanded and/or exposing untyped data (map[string]string) where it doesn't make any sense.

Second, to cleanly separate the framework infrastructure from the actual web application.

Ultimately, these goals result in actions which are clean, testable and look like:

    package "sessions"
    //context is specific to your application
    //auwfg.Response is an interface which maps well http.ResponseWriter
    func Create(context *Context) auwfg.Response {
      input := context.Body.(*LoginInput)
      usermame := input.Username
      password := input.Rassword
      remember := input.LongLived

      //or create your own response types
      return auwfg.JsonResponse(`{"token":"blah"`}, 201)
    }

## Configuration
You can start the server, using default settings, by calling `auwfg.Run(auwfg.Configure())`. However, `Configure()` exposes a fluent interface which lets specify various settings:

    c := auwfg.Configure().Address("127.0.0.1:3944").NotFound('{"error":404}')
    auwfg.Run(c)

Possible configuration options are:

- `Address`: the address to listen to
- `NotFound`: the body to reply with when a route isn't found
- `NotFoundResponse`: the response to reply with when a route isn't found (gives more control than the simpler `NotFound`). Note that once set, this response is available as `auwfg.NotFound` for use in your own application
- `BodyTooLarge`: the body to reply with when an input body is too large (see `BodyPool`)
- `BodyTooLargeResponse`: the response to reply with when an input body is too large (gives more control than the simpler `BodyTooLarge`)
- `InvalidFormat`: the body to reply with when an input format is invalid (not valid json)
- `InvalidFormatResponse`: the response to reply with when an input format is invalid (gives more control than the simpler `InvalidFormat`)
- `InternalServerError`: the body to reply with when an error occurs
- `InternalServerErrorResponse`: the response to reply with when an error occurs (gives more control than the simpler `InternalServerError`). Note that once set, this response is available as `auwfg.InternalServerError` for use in your own application
- `ContextFactory`: explained below
- `Dispatcher`: explained below
- `Route`: explained below
- `BodyPool`: explained below

## ContextFactory and Dispatcher
Having strongly-typed, application specific context relies on specifying a custom `ContextFactory` and `Dispatcher`.

The `ContextFactory` takes an instance of `*auwfg.BaseContext` and returns whatever context instance you care about. Whatever object you return need not have a relationship with the provided `BaseContext`. However, most usage will likely want to make use of Go's implicit composition (embedding):

    type Context struct {
      User *User
      // more app-specific fields
      *auwfg.BaseContext
    }

    func ContextFactory(c *auwfg.BaseContext) interface{} {
      return &Context{loadUser(c.Req), c}
    }

    c := auwfg.Configure().ContextFactory(ContextFactory)
    auwfg.Run(c)

Unfortunately, on its own, `ContextFactory` is not enough. While we've created a specific type, said type information is unknown to AUWFG itself. Rather than having each of action deal with typeless data, we use a custom `Dispatcher`:

    func Dispatcher(route *auwfg.Route, context interface{}) auwfg.Response {
      //you probably won't need to do anything with route
      return route.Action.(func(*Context) auwfg.Response)(context.(*Context))
    }
    ...
    c := auwfg.Configure().ContextFactory(ContextFactory).Dispatcher(Dispatcher)
    auwfg.Run(c)

**Your** dispatcher knows that **your** actions take a parameter of type `Context`, and thus can safely cast both the action and the parameter.

## Routing
Routing is strict and configured via the `Route` method of the configuration. Every route must have a mimimum of 4 components:

- method
- version
- resource
- action

The first three are `strings`, the last is the action handler. Let's look at some configuration example:

    c := auwfg.Configure()

    //matches GET /v1/gholas.ext
    c.Route(auwfg.R("LIST", "v1", "gholas", gholas.List))

    //matches GET /v1/gholas/SOMEID.EXT
    c.Route(auwfg.R("GET", "v1", "gholas", gholas.Show))

    //matches POST /v1/gholas.ext
    c.Route(auwfg.R("POST", "v1", "gholas", gholas.Create))

    //matches PUT /v1/gholas/SOMEID.ext
    c.Route(auwfg.R("PUT", "v1", "gholas", gholas.Update))

    //matches DELETE /v1/gholas/SOMEID.ext
    c.Route(auwfg.R("DELETE", "v1", "gholas", gholas.Delete))

Note that the method for a `GET` with no id is called `LIST`. The `gholas.List|Show|Create|Update|Delete` actions are any function which can be invoked by your dispatcher.

The captured parameters, "v1", "gholas" and "SOMEID" are available in the BaseContext.Params (Version, Resource and Id fields respectively), which you've hopefully preserved in your own context.

## Parsing Request Bodies
Any route can have an associated `BodyFactory`, configured as:

    c.Route(auwfg.R("POST", "v1", "gholas", gholas.Create).BodyFactory(func() { return new(GholaInput)} ))

A `BodyFactory` simply returns an instance of an object which will be used with the `encoding/json` package to convert the request body into the desired input.

Sadly, type information is lost, in actions must currently cast the `BaseContext.Body` to the appropriate type:

    func Create(context *Context) auwfg.Response {
      input := context.Body.(*GholaInput)
      ...
    }

When no body is present, an empty instance is made available. If, however, there's an error parsing the input, and `InvalidFormat` is returned.

### Pool Size
A [fixed-length byte pool](https://github.com/viki-org/bytepool) is used to parse input bodies. The size of this pool is configured with the `BodyPool` configuration method. For example, to specify a max input body size of 64K and to pre-allocate 512 slots, you'd use:

    c := auwfg.Configure().BodyPool(64 * 1024, 512)
