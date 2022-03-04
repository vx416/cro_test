package auth

import (
	"cro_test/internal/model"
)

func (s *AuthSuite) TestSvc_Singup() {
	tcs := []struct {
		name   string
		email  string
		pwd    string
		hasErr bool
	}{
		{
			"signup_ok",
			"vic@gmail.com", "test1234", false,
		},
		{
			"email duplicated",
			"vic@gmail.com", "test1234", true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			u, err := s.Service.SignUp(s.Ctx, tc.email, tc.pwd)
			if tc.hasErr {
				s.Error(err)
				return
			}
			s.Require().NoError(err)

			dbU, err := s.MySQLDao.GetUser(s.Ctx, model.GetUserOpts{
				Email: tc.email,
			})
			s.Require().NoError(err)
			s.Equal(u.PasswordHash, dbU.PasswordHash)
			s.NoError(dbU.CmpPwd(tc.pwd))
		})
	}
}
