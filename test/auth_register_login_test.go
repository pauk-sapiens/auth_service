package test

import (
	"auth/test/suite"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	auth1 "github.com/pauk-sapiens/protos/gen/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	email := gofakeit.Email()
	password := randomFakePW()

	respReq, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{Email: email, Password: password})
	require.NoError(t, err)
	assert.NotEmpty(t, respReq.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{Email: email, Password: password, AppId: appID})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return []byte(appSecret), nil })
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReq.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["appID"].(float64)))

}

func randomFakePW() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
