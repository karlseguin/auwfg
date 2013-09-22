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
      return auwfg.Json(`{"token":"blah"`}, 201)
    }

## Configuration
You can start the server, using default settings, by calling `auwfg.Run(auwfg.Configure())`. However, `Configure()` exposes a fluent interface which lets you specify various settings:

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
- `ContextFactory`: [explained below](#cfd)
- `Dispatcher`: [explained below](#cfd)
- `Route`: [explained below](#routing)
- `BodyPool`: [explained below](#bodypool)
- `InvalidPool`: [explained below](#invalidpool)

<a id="cfd"></a>
## ContextFactory and Dispatcher
Having strongly-typed, application specific context relies on specifying a custom `ContextFactory` and `Dispatcher`.

The `ContextFactory` takes an instance of `*auwfg.BaseContext` and returns whatever context instance you care about. Whatever object you return need not have a relationship with the provided `BaseContext`. However, most usage will likely want to make use of Go's implicit composition (embedding):

    type Context struct {
      User *User
      // more app-specific fields
      // ...
      *auwfg.BaseContext
    }

    func ContextFactory(c *auwfg.BaseContext) interface{} {
      return &Context{loadUser(c.Req), c}
    }

    c := auwfg.Configure().ContextFactory(ContextFactory)
    auwfg.Run(c)

Unfortunately, on its own, `ContextFactory` is not enough. While we've created a specific type, said type information is unknown to AUWFG itself. Rather than having each action deal with typeless data, we use a custom `Dispatcher`:

    func Dispatcher(route *auwfg.Route, context interface{}) auwfg.Response {
      //you probably won't need to do anything with route
      return route.Action.(func(*Context) auwfg.Response)(context.(*Context))
    }
    ...
    c := auwfg.Configure().ContextFactory(ContextFactory).Dispatcher(Dispatcher)
    auwfg.Run(c)

**Your** dispatcher knows that **your** actions take a parameter of type `Context`, and thus can safely cast both the action and the parameter. Also, the dispatcher is a great place to execute pre and post filters for all actions.

<a id="routing"></a>
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

Note that the method for a `GET` with no id is called `LIST`. The `gholas.List|Show|Create|Update|Delete` actions are function which can be invoked by your dispatcher.

The captured parameters, "v1", "gholas" and "SOMEID" are available in the `BaseContext.Params` (`Version`, `Resourc`e and `Id` fields respectively), which you've hopefully preserved in your own context.

### Nested Routes
A nested route, is configured by specifying the route's `Parent`:

    c.Route(auwfg.R("LIST", "v1", "history", gholas.Show).Parent("gholas"))

Will match `GET /v1/gholas/SOMEID/history.json`

`SOMEID` is available in `Params.ParentId`

## Parsing Request Bodies
Any route can have an associated `BodyFactory`, configured as:

    c.Route(auwfg.R("POST", "v1", "gholas", gholas.Create).BodyFactory(func() { return new(GholaInput)} ))

A `BodyFactory` simply returns an instance of an object which will be used with the `encoding/json` package to convert the request body into the desired input.

Sadly, type information is lost; actions must currently cast the `BaseContext.Body` to the appropriate type:

    func Create(context *Context) auwfg.Response {
      input := context.Body.(*GholaInput)
      ...
    }

When no body is present, an empty instance is made available. If, however, there's an error parsing the input, and `InvalidFormat` is returned.

<a id="bodypool"></a>
### Pool Size
A [fixed-length byte pool](https://github.com/viki-org/bytepool) is used to parse input bodies. The size of this pool is configured with the `BodyPool` configuration method. For example, to specify a max input body size of 64K and to pre-allocate 512 slots, you'd use:

    c := auwfg.Configure().BodyPool(512, 64 * 1024)

While the maximum size of the body is fixed (64K with the above configuration), the number of available buffers (512 above) will grow as needed.

## Responses
A `auwfg.Response` must implement three members:

- `Body() []byte`
- `Status() int`
- `Header() http.Header`
- `Close()`

The `Json(body string, status int)` helper should prove helpful.

The `Fatal(err error)` helper should be used when an `InternalServerError` should be returned and an error logged (using the standard logger)

The `auwfg.Deleted` variable is a response which returns 204 and a content-length of 0

### Closing Responses
The reason for responses to implement `Close` is to make it possible to use buffer pools when generating responses. If you're creating your own `Response`, it would not be unusual for `Close` to do nothing. Furthermore, under normal conditions, AUWFG will take care of closing responses. However, if you intercept a `auwfg.Response` and return a different response, you should call `Close()` on it. For example:

    if res, valid := validate(input); valid == false {
      res.Close() // <- since auwfg will never see this response, you need to take care of closing it yourself
      return Json("why?!")
    }

## Validation
AUWFG has some basic input validation facilities. Validation works in two phases. The first phase is to define error messages. The second phase is to validate the actual data:

    func init() {
      auwfg.Define("required_username").Field("username").Message("Username is required")
    }

    func Create(context *Context) auwfg.Response {
      input := context.Body.(*LoginInput)
      if res, valid := validate(input); valid == false {
        return res
      }
      ...
    }

    func validate(input *LoginInput) (auwfg.Response, bool) {
      validator := auwfg.Validator()
      validator.Required(input.UserName, "required_username")
      return validator.Response()
    }

The last parameter of every validation function is the `definitionId`. The `definitionId` maps to the parameter given to `Define`. If the `definitionId` was not `Defined`, the validation will panic.

The following validation methods are currently support:

- `Required(s, definitionId string)`
- `Len(s string, min, max int, definitionId string)`
- `MinLen(s string, min int, definitionId string)`
- `MaxLen(s string, max int, definitionId string)`
- `Same(a, b string, definitionId string)`

In addition to calling `Response`, which returns `(auwfg.Response, bool)`, `IsValid() bool` is also available.

### Custom Validation
Calling `AddError(definitionId string)` is currently the only way to have "custom errors":

    func validate(input *LoginInput) (auwfg.Response, bool) {
      validator := auwfg.Validator()
      if isTooCool(input) { validator.AddError("too_cool") }
      ...
      return validator.Response()
    }

<a id="invalidpool"></a>
### Pool Size
A [fixed-length byte pool](https://github.com/viki-org/bytepool) is used to generate validation response. The size of this pool is configured with the `InvalidPool` configuration method. For example, to specify a max response size of 32 and to pre-allocate 256 slots, you'd use:

    c := auwfg.Configure().InvalidPool(256, 32 * 1024)

While the maximum size is fixed (32K with the above configuration), the number of available buffers (256 above) will grow as needed.
