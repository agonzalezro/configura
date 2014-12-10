configura
=========

![circleci](https://circleci.com/gh/agonzalezro/configura.png)

Configura will allow you to store all your configuration in environment
variables, but it will also let you set some defaults.

History
-------

If you already know the [twelve-factor](http://12factor.net/) methodology you
will be more than aware of their third point: http://12factor.net/config were
they highly recommend to use environment variables to all the configs in your
app. I will also add that using docker this way of storing the configuration is
quite handy and some times allows you to avoid other techniques like: ansible,
puppet, chef...

The idea of configura is keeping all this configuration loading easy-peasy, let
me show you that I am not lying:

The most basic thing ever
-------------------------

    type Config struct {
        SomeString   string
    }

    c := Config{}
    err := Load("TWAPP_", &c)

The example showed above is the simplest one that you can find. You will need
to have a env var called `TWAPP_SOMESTRING` before launching your program or at
the moment of loading the conf it will miserably fail (actually, it will just
return an error, you can do whatever you want with it :)

Please, bear in mind that `TWAPP_` prefix is something that you specify at
loading time, but you can use whatever prefix you want, or even don't use
anything.

But what else?
--------------

It supports:

- `int`s
- `string`s
- `bool`s
- `time.Duration`s
- `float32`s & `float64`s

Of course, to use some of those types, your env variables should have parsable
values (aka, don't set it to "Alex" when you are expecting an int) and that's
all!

### Example

This struct:

    type Config struct {
        SomeString   string
        SomeInt      int
        SomeBool     bool
        SomeDuration time.Duration
    }

Will be expecting env vars like this one (if we continue using the `TWAPP_`
prefix):

    TWAPP_SOMESTRING=string
    TWAPP_SOMEINT=1
    TWAPP_SOMEBOOL=y # or true, or yes... whatever accepted by ParseBool
    TWAPP_SOMEDURATION=1s # or whatever accepted by ParseDuration


More complex struct
-------------------

Using struct tags you will be able to change the behaviour of the package:

- The var name to be looked up on the system can be override using
  `configura:"OVERRIDE"`.
- In case that the var was not found on the system you will be able to set some
  default with: `configura:",defaultvalue"`.
- And finally, if you want to use both: `configura:"OVERRIDE,defaultvalue"`.

### Example:

For doing so the struct will look like this:

    type Config struct {
        // Will read the value from TWAPP_SOMESTRING
        SomeString   string

        // Will read it from ANOTHERINT and if it's not there will fail
        SomeInt      int `configura:"ANOTHERINT"`

        // Will default to true reading the value from TWAPP_SOMEBOOL
        SomeBool     bool `configura:",true"`

        // Will read it from ANTOHERDURATION but defaulting to 1s
        SomeDuration time.Duration `configura:"ANOTHERDURATION,1s"`
    }

If you need more help, you can check the [package
documentation](https://godoc.org/github.com/agonzalezro/configura) or [ping me
on twitter](http://twitter.com/agonzalezro).
