// Package middleware defines the middleware implementation available
package middleware

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	ginJwt "github.com/appleboy/gin-jwt/v2"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	graphModel "employee-management-system/graph/model"
	"employee-management-system/model"
	"employee-management-system/pkg/environment"
	"employee-management-system/pkg/helper"
)

type (
	// Tokens object
	Tokens struct {
		AccessToken        string
		RefreshToken       string
		AccessTokenExpiry  string
		RefreshTokenExpiry string
	}
)

var (
	identityKey   = "id"
	realm         = "c-o-m-p-a-n-y"
	maxRefresh    = jwtRefreshTokenExpiry()
	accessTimeout = jwtAccessTokenExpiry()
	// claimsID id key for middleware claims
	claimsID        = "id"
	claimsExpiry    = "exp"
	claimsCreatedAt = "orig_iat"
	// ErrFailedAuthentication incorrect email or password
	ErrFailedAuthentication = errors.New("incorrect email or password")
	// ErrAccountSuspended user account is suspended
	ErrAccountSuspended = errors.New("user account is suspended")
	// ErrCorruptAdminAccount corrupt admin account
	ErrCorruptAdminAccount = errors.New("corrupt admin account")
	// ErrCorruptAgentAccount incorrect email or password
	ErrCorruptStaffAccount = errors.New("corrupt empployee account")
	// ErrUnexpectedSigningMethod occurs when a token does not conform to the expected signing method
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	// ErrInvalidToken indicates JWT token has expired. Can't refresh.
	ErrInvalidToken = errors.New("token is invalid")
	// ErrUnauthorized reports unauthorized user
	ErrUnauthorized = errors.New("you are not authorized")
)

// JwtAuthenticator authenticates user by username and password
func (m *Middleware) JwtAuthenticator(c *gin.Context, u, p string) (*graphModel.AuthResponse, error) {
	// attempt login
	user, err := m.userStorage.Authenticate(c, u, p)
	if err != nil {
		return nil, ErrFailedAuthentication
	}

	// get relationship values (admin/agent), return error if any occurs
	user, err = m.evalKindForRelationship(c, user)
	if err != nil {
		return nil, err
	}

	tokens, err := m.GenerateTokens(c, user)
	if err != nil {
		return nil, ginJwt.ErrFailedTokenCreation
	}

	if m.jwt.SendCookie {
		maxage := int(time.Now().Add(accessTimeout).Unix() - time.Now().Unix())
		c.SetCookie(
			m.jwt.CookieName,
			tokens.AccessToken,
			maxage,
			"/",
			m.jwt.CookieDomain,
			m.jwt.SecureCookie,
			m.jwt.CookieHTTPOnly,
		)
	}

	return &graphModel.AuthResponse{
		Token:              &tokens.AccessToken,
		Refresh:            &tokens.RefreshToken,
		User:               user,
		AccessTokenExpiry:  &tokens.AccessTokenExpiry,
		RefreshTokenExpiry: &tokens.RefreshTokenExpiry,
	}, nil
}

// GenerateTokens creates both access and refresh tokens
func (m *Middleware) GenerateTokens(c *gin.Context, user *model.User) (*Tokens, error) {
	accessToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	accessClaims := accessToken.Claims.(jwtGo.MapClaims)

	refreshToken := jwtGo.New(jwtGo.GetSigningMethod(m.jwt.SigningAlgorithm))
	refreshClaims := refreshToken.Claims.(jwtGo.MapClaims)

	if m.jwt.PayloadFunc != nil {
		for key, value := range m.jwt.PayloadFunc(user) {
			accessClaims[key] = value
			refreshClaims[key] = value
		}
	}
	accessExpire := time.Now().Add(accessTimeout)
	refreshExpire := time.Now().Add(maxRefresh)

	accessClaims[claimsID] = user.ID
	accessClaims[claimsExpiry] = accessExpire.Unix()
	accessClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	refreshClaims[claimsID] = user.ID
	refreshClaims[claimsExpiry] = refreshExpire.Unix()
	refreshClaims[claimsCreatedAt] = m.jwt.TimeFunc().Unix()

	accessTokenString, err := m.signedString(accessToken)
	if err != nil {
		return nil, err
	}
	refreshTokenString, err := m.signedString(refreshToken)
	if err != nil {
		return nil, err
	}

	// save refresh token in cookie, for future check
	c.SetCookie(
		fmt.Sprintf("%v", user.ID),
		refreshTokenString,
		int(time.Now().Add(maxRefresh).Unix()-time.Now().Unix()),
		"/",
		m.jwt.CookieDomain,
		m.jwt.SecureCookie,
		m.jwt.CookieHTTPOnly,
	)
	// Todo: delete cookie on logout

	return &Tokens{
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		AccessTokenExpiry:  accessExpire.String(),
		RefreshTokenExpiry: refreshExpire.String(),
	}, err
}

// JwtAuthorization returns an authorized User
func (m *Middleware) JwtAuthorization(c *gin.Context) (*model.User, error) {
	var user *model.User
	claims, err := m.jwt.GetClaimsFromJWT(c)
	if err != nil {
		return nil, err
	}

	userID := claims[claimsID].(string)
	if user = m.cache.user(userID); user == nil {
		// get user by ID
		dbUser, err := m.userStorage.GetUserByID(c, userID)
		if err != nil {
			return nil, err
		}

		return m.evalKindForRelationship(c, &dbUser)
	}

	return user, nil
}

// GetGinJWTMiddleware returns GinJWTMiddleware
func (m *Middleware) GetGinJWTMiddleware() *ginJwt.GinJWTMiddleware {
	return m.jwt
}

func (m *Middleware) evalKindForRelationship(ctx context.Context, user *model.User) (*model.User, error) {
	switch user.Kind {
	case model.KindAdministrator:
		admin, err := m.adminStorage.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, ErrCorruptAdminAccount
		}
		user.Admin = &admin
		if !user.Admin.Active {
			return nil, ErrAccountSuspended
		}
	case model.KindAgent:
		var partners []model.Partner
		agent, err := m.agentStorage.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, ErrCorruptAgentAccount
		}
		user.Agent = &agent
		if !user.Agent.Active {
			return nil, ErrAccountSuspended
		}
		agentPartners, err := m.agentPartnerStorage.GetByAgentID(ctx, agent.ID)
		if err != nil {
			return nil, ErrFetchingAgentPartner
		}
		if len(agentPartners) > 0 {
			for _, agentPartner := range agentPartners {
				partner, err := m.partnerStorage.GetPartnerByID(ctx, agentPartner.PartnerID)
				if err != nil {
					return nil, ErrFetchingPartner
				}
				if partner.Active {
					partnerSetting, err := partner.GetSettings()
					if err != nil {
						return nil, ErrFetchingPartner
					}
					if partnerSetting.WorkTime != nil {
						timeFrom := partnerSetting.WorkTime[0]
						timeTo := partnerSetting.WorkTime[1]
						if timeFrom != nil && timeTo != nil {
							timeNow, err := helper.SetGetTimezone(helper.GetEnvironmentVariable("TIMEZONE"), time.Now())
							if err != nil {
								return nil, ErrFetchingPartner
							}
							minute := timeNow.Minute()
							hour := timeNow.Hour()
							splitFromTime := strings.Split(*timeFrom, ":")
							fromTimeHour, err := strconv.Atoi(splitFromTime[0])
							if err != nil {
								return nil, ErrFetchingPartner
							}
							fromTimeMinute, err := strconv.Atoi(splitFromTime[1])
							if err != nil {
								return nil, ErrFetchingPartner
							}

							splitToTime := strings.Split(*timeTo, ":")
							toTimeHour, err := strconv.Atoi(splitToTime[0])
							if err != nil {
								return nil, ErrFetchingPartner
							}
							toTimeMinute, err := strconv.Atoi(splitToTime[1])
							if err != nil {
								return nil, ErrFetchingPartner
							}
							if !((hour > fromTimeHour && hour < toTimeHour) || ((fromTimeHour == hour && minute >= fromTimeMinute) || (toTimeHour == hour && minute <= toTimeMinute))) {
								continue
							}
						}
					}
					partners = append(partners, partner)
				}
			}
		}
		agent.Partners = partners
	}
	// if not suspended
	// cache user - new login
	m.cache.cacheUser(user)

	return user, nil
}

// LogoutHandler can be used by clients to remove the middleware cookie (if set)
func (m *Middleware) LogoutHandler(c *gin.Context) {
	// delete auth cookie
	if m.jwt.SendCookie {
		c.SetCookie(
			m.jwt.CookieName,
			"",
			-1,
			"/",
			m.jwt.CookieDomain,
			m.jwt.SecureCookie,
			m.jwt.CookieHTTPOnly,
		)
	}
}

func (m *Middleware) signedString(token *jwtGo.Token) (string, error) {
	var tokenString string
	var err error
	if m.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(m.pKey)
	} else {
		tokenString, err = token.SignedString(m.jwt.Key)
	}
	return tokenString, err
}

func (m *Middleware) usingPublicKeyAlgo() bool {
	switch m.jwt.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

func jwtAccessTokenExpiry() time.Duration {
	env, _ := environment.New()
	ttl, err := strconv.Atoi(env.Get("JWT_ACCESS_TOKEN_EXPIRY"))
	if err != nil {
		return time.Minute * 15
	}
	return time.Minute * time.Duration(ttl)
}

func jwtRefreshTokenExpiry() time.Duration {
	env, _ := environment.New()
	ttl, err := strconv.Atoi(env.Get("JWT_REFRESH_TOKEN_EXPIRY"))
	if err != nil {
		return time.Hour * 24
	}
	return time.Minute * time.Duration(ttl)
}

// ValidateRefreshToken validates refresh token
func (m *Middleware) ValidateRefreshToken(z zerolog.Logger, c *gin.Context, token string) (*string, error) {
	tokenGotten, err := jwtGo.Parse(token, func(token *jwtGo.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			z.Error().Msgf("RefreshToken unexpected signing method: (%v)", token.Header["alg"])

			return nil, ErrUnexpectedSigningMethod
		}
		env, _ := environment.New()
		return []byte(env.Get("SIGNING_SECRET_KEY")), nil
	})

	//any error may be due to token expiration
	if err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	//is token valid?
	if err = tokenGotten.Claims.Valid(); err != nil {
		z.Err(err).Msgf("RefreshToken error: %v", err)
		return nil, err
	}

	claims, ok := tokenGotten.Claims.(jwtGo.MapClaims)
	claimsUUID := claims[claimsID].(string)
	//get the last refresh token for this user/customer
	refreshTokenCookie, err := c.Cookie(claimsUUID)
	//error may be due to cookie expiration OR a new refresh token has been generated
	if err != nil || refreshTokenCookie != token {
		z.Err(err).Msgf("RefreshToken:Cookie error: %v", err)
		return nil, ErrInvalidToken
	}

	if ok && tokenGotten.Valid {
		//convert the interface to uuid.UUID
		parsedUUID, err := uuid.Parse(claimsUUID)
		if err != nil {
			z.Err(err).Msgf("RefreshToken: Invalid user (%v)", err)
			return nil, err
		}

		return &parsedUUID, nil
	}

	return nil, ErrInvalidToken
}
