package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*mux.Router
	shutdown   chan os.Signal
	mw         []Middleware
	pathPrefix string
}

func NewApp(shutdown chan os.Signal, pathPrefix string, mw ...Middleware) *App {
	return &App{
		Router:     mux.NewRouter(),
		shutdown:   shutdown,
		mw:         mw,
		pathPrefix: pathPrefix,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) HandleNoMiddleWareFunc(path string, h Handler) *mux.Route {

	return a.handleFunc(h, path)
}

// HandleNoMiddleware sets a handler function for a given HTTP method and path pair
// to the application server mux. Does not include the application middleware.
func (a *App) HandleFunc(path string, h Handler, mw ...Middleware) *mux.Route {

	h = wrapMiddleWare(mw, h)
	h = wrapMiddleWare(a.mw, h)

	return a.handleFunc(h, path)
}

// ===========================================================================

func (a *App) handleFunc(h Handler, path string) *mux.Route {
	f := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now(),
		}

		ctx := SetValues(r.Context(), &v)

		if err := h(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}

	}
	routes := a.Router.PathPrefix(a.pathPrefix).Subrouter()

	return routes.HandleFunc(path, f)
}

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
