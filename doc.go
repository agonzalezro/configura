// Configura will allow you to store all your configuration in environment
// variables, but it will provide some other utilities as default values as well.

// If you already know [twelve-factor](http://12factor.net/) methodology you will
// be more than aware of their third point: http://12factor.net/config were they
// highly recommend to use environment variables to all the configs in your app. I
// will also add that using docker this way of storing the configuration is quite
// handy and some times allows you to avoid other techniques like: ansible,
// puppet, chef...
package configura
