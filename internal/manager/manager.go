package manager

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/unanet/eve-bot/internal/config"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/identity"
	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/middleware"
	"go.uber.org/zap"
)

// Key to use when setting the request ID.
type ctxKeyTokenClaimsID int

// TokenClaimsRequestIDKey is the key that holds the unique Token Claims ID in a request context.
const TokenClaimsRequestIDKey ctxKeyTokenClaimsID = 0

func OpenIDConnectOpt(id *identity.Service) Option {
	return func(svc *Service) {
		svc.oidc = id
	}
}

type Option func(*Service)

type Service struct {
	cbstate string
	cfg     *config.Config
	oidc    *identity.Service
}

func (s *Service) OpenIDService() *identity.Service {
	return s.oidc
}

func NewService(cfg *config.Config, opts ...Option) *Service {
	svc := &Service{cfg: cfg, cbstate: "eve-bot"}

	for _, opt := range opts {
		opt(svc)
	}

	return svc

}

func (s *Service) AuthCodeURL(chatUser string) string {
	return s.oidc.AuthCodeURL(chatUser)
	//return s.oidc.AuthCodeURL(s.cbstate)
}

func (s *Service) AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			unknownToken := jwtauth.TokenFromHeader(r)

			middleware.Log(r.Context()).Debug("in middleware " + r.URL.String())

			if len(unknownToken) == 0 {
				render.Respond(w, r, errors.ErrUnauthorized)
				return
			}

			verifiedToken, err := s.oidc.Verify(ctx, unknownToken)
			if err != nil {
				middleware.Log(ctx).Debug("invalid token", zap.Error(err))
				http.Redirect(w, r, s.oidc.AuthCodeURL(s.cbstate), http.StatusFound)
				return
			}

			//var idTokenClaims = new(json.RawMessage)
			var claims = new(jwt.MapClaims)
			if err := verifiedToken.Claims(&claims); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, TokenClaimsRequestIDKey, claims)))
		}
		return http.HandlerFunc(hfn)
	}
}

// GetTokenClaims returns the verified token claims
// Returns nil if unknown
func (s *Service) GetTokenClaims(ctx context.Context) jwt.MapClaims {
	if ctx == nil {
		return nil
	}
	if claims, ok := ctx.Value(TokenClaimsRequestIDKey).(jwt.MapClaims); ok {
		return claims
	}
	return nil
}

func (s *Service) ReadOnlyMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log.Logger.Info("API Listener", zap.String("in middleware", r.URL.String()))

			if s.cfg.ReadOnly && r.Method != http.MethodGet {
				err := errors.NewRestError(http.StatusServiceUnavailable, "Unable to perform action. API is in read only mode")
				middleware.Log(ctx).Debug("invalid token", zap.Error(err))
				http.Error(w, err.Error(), err.Code)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
