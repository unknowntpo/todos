package middleware

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/unknowntpo/todos/config"
	"github.com/unknowntpo/todos/internal/domain"
	"github.com/unknowntpo/todos/internal/helpers"
	"github.com/unknowntpo/todos/internal/logger"
	"github.com/unknowntpo/todos/pkg/validator"

	"github.com/felixge/httpsnoop"
	"github.com/pkg/errors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

type Middleware struct {
	config  *config.Config
	usecase domain.UserUsecase
	logger  logger.Logger
}

func New(cfg *config.Config, uu domain.UserUsecase, logger logger.Logger) *Middleware {
	return &Middleware{config: cfg, usecase: uu, logger: logger}
}

func (mid *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// TODO: How to handle panic?
				// Try to log stack trace ?
				mid.logger.PrintError(errors.New(fmt.Sprintf("%s", err)), nil)
				helpers.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	// Declare a mutex and a map to hold the clients' IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once
	// every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if mid.config.Limiter.Enabled {
			// Use the realip.FromRequest() function to get the client's real IP address.
			ip := realip.FromRequest(r)

			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't, then
			// initialize a new rate limiter and add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(mid.config.Limiter.Rps), mid.config.Limiter.Burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the rate limiter for the current IP address. If
			// the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
			// response, just like before.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				helpers.RateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This indicates to any caches that the response may vary
		// based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = helpers.ContextSetUser(r, domain.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Expect the value of the Authorization header to be in the format
		// "Bearer <token>".
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			helpers.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if domain.ValidateTokenPlaintext(v, token); !v.Valid() {
			helpers.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		ctx := r.Context()
		user, err := mid.usecase.GetForToken(ctx, domain.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				helpers.InvalidAuthenticationTokenResponse(w, r)
			default:
				helpers.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = helpers.ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

// RequireAuthenticatedUser checks that a user is not anonymous.
func (mid *Middleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.ContextGetUser(r)

		if user.IsAnonymous() {
			helpers.AuthenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireActivatedUser checks if a user is both authenticated and activated.
func (mid *Middleware) RequireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user information from the request context.
		user := helpers.ContextGetUser(r)

		if !user.Activated {
			helpers.InactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	// Use mid.RequireAuthenticatedUser to check if a user is authenticated.
	return mid.RequireAuthenticatedUser(fn)
}

func (mid *Middleware) Metrics(next http.Handler) http.Handler {
	// Initialize the new expvar variables when the middleware chain is first built.
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	// etotal_responses_sent_by_status holds the count of responses for each HTTP status
	// code.
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	// The following code will be run for every request...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the requests received count, like before.
		totalRequestsReceived.Add(1)

		// Call the httpsnoop.CaptureMetrics() function, passing in the next handler in
		// the chain along with the existing http.ResponseWriter and http.Request. This
		// returns the metrics struct that we saw above.
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		// Increment the response sent count, like before.
		totalResponsesSent.Add(1)

		// Get the request processing time in microseconds from httpsnoop and increment
		// the cumulative processing time.
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		// Use the Add() method to increment the count for the given status code by 1.
		// Note that the expvar map is string-keyed, so we need to use the strconv.Itoa()
		// function to convert the status code (which is an integer) to a string.
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

func (mid *Middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Origin" header.
		w.Header().Add("Vary", "Origin")

		// Get the value of the request's Origin header.
		origin := r.Header.Get("Origin")

		// Only run this if there's an Origin request header present AND at least one
		// trusted origin is configured.
		if origin != "" && len(mid.config.Cors.TrustedOrigins) != 0 {
			// Loop through the list of trusted origins, checking to see if the request
			// origin exactly matches one of them.
			for i := range mid.config.Cors.TrustedOrigins {
				if origin == mid.config.Cors.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
