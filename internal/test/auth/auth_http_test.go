package auth

import (
	"cro_test/internal/delivery/http"
	"cro_test/internal/model"
	"cro_test/internal/test/factory"
	"cro_test/pkg/config"
	"cro_test/pkg/middleware"
)

func (s *AuthSuite) TestHttp_Auth() {
	user := factory.User.MustBuild().(*model.User)
	pwd := "test_123456"
	u, err := s.Service.SignUp(s.Ctx, user.Email, pwd)
	s.Require().NoError(err)
	req, resp := s.HttpHelper.BuildRequest("POST", "/api/v1/auth", map[string]interface{}{
		"eamil":    user.Email,
		"password": pwd,
	})
	s.Echo.Server.Handler.ServeHTTP(resp, req)
	token := http.TokenResponse{}
	err = s.HttpHelper.GetResponseData(resp, &token)
	s.Require().NoError(err)
	s.NotEmpty(token.AccessToken)
	s.T().Log(token.AccessToken)
	cliams, err := middleware.ParseWithClaims(token.AccessToken, config.GetJwtSecret())
	s.Require().NoError(err)
	s.Equal(cliams.UserID, u.ID)
}
