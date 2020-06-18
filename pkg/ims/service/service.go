package service

import (
	"crypto/rand"
	"github.com/jmoiron/sqlx"
	"math/big"
	"strconv"

	"github.com/omc-college/management-system/pkg/ims/models"
	"github.com/omc-college/management-system/pkg/ims/repository/postgresql"
	"github.com/omc-college/management-system/pkg/pwd"
)

// token returns random int for verification token
func token() string {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(999999999999999)))
	rand := strconv.Itoa(int(r.Int64()))
	return rand
}

type SignUpService struct {
	db *sqlx.DB
}

func NewSignUpService(DB *sqlx.DB) *SignUpService {
	return &SignUpService{
		db: DB,
	}
}

func (service *SignUpService) SignUp (request *models.SignupRequest) error {
	var cred models.Credentials
	var tok models.EmailVerificationTokens
	var err error
	var user models.User

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.Email

	tok.VerificationToken = token()

	cred.Salt = pwd.Salt(256 - len(request.Password))

	cred.PasswordHash, err = pwd.Hash(request.Password, cred.Salt)
	if err != nil {
		return err
	}

	tx, err := service.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	signUpRepo := postgresql.NewSignUpRepository(tx)
	credRepo := postgresql.NewCredentialsRepository(tx)
	userRepo := postgresql.NewUsersRepository(tx)

	err = userRepo.InsertUser(&user)
	if err != nil {
		return err
	}

	tok.ID = user.ID
	cred.ID = user.ID

	err = credRepo.InsertCredentials(&cred)
	if err != nil {
		return err
	}

	err = signUpRepo.InsertEmailVerificationToken(&tok)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (service *SignUpService) EmailAvailable (email string) (bool, error) {
	var user *models.User
	var exist bool
	var err error

	userRepo := postgresql.NewUsersRepository(service.db)

	user, err = userRepo.GetUserByEmail(email)
	if err != nil {
		return exist, err
	}

	if user.Email != "" {
		exist = true
	}

	return exist, nil
}

func (service *SignUpService) EmailVerificationToken (token *models.EmailVerificationTokens) error {
	var user models.User
	var err error

	tx, err := service.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userRepo := postgresql.NewUsersRepository(tx)
	signUpRepo := postgresql.NewSignUpRepository(tx)

	user.Email, err = userRepo.GetEmailByToken(token)
	if err != nil {
		return err
	}

	err = signUpRepo.SetUserVerified(&user)
	if err != nil {
		return err
	}

	err = signUpRepo.DeleteEmailVerificationToken(token)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
	func (service *SignUpService) ResetPassword(request *models.SignupRequest) error {
	var user *models.User
	var err error
	var tok models.EmailVerificationTokens

	user.Email = request.Email
	tok.VerificationToken = token()

	tx, err := service.db.Beginx()
	if err != nil {
		return err
	}

    userRepo := postgresql.NewUsersRepository(tx)
	user, err = userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	from := "noreply@managementsystem.com"
	host := "mail.managementsystem.com"
	auth := smtp.PlainAuth("", from, "password", host)
	message := fmt.Sprintf(`To: "Some User" <someuser@example.com>
From: "mailbot" <%s>
Subject: Password reset
"Dear %s,

You recently requested to reset your password for your managementsystem account. Click the link below to reset it.
%s
If you did not request a password reset, please ignore this email or reply to let us know. This password reset link is only valid for the next 2 hours."
`,from,user.Email,fmt.Sprintf("http://managementsystem.com/confirmedreset?token=%s",tok))
	if user.Email != "" {
		smtp.SendMail(
			host+":25", auth, from,[]string{user.Email}, []byte(message))

	}
	signUpRepo := postgresql.NewSignUpRepository(tx)
	err = signUpRepo.InsertEmailVerificationToken(&tok)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func (service *SignUpService) ConfirmedReset(request *models.SignupRequest,tok *models.EmailVerificationTokens) error {
	var user *models.User
	var cred *models.Credentials
	var err error

	tx, err := service.db.Beginx()
	if err != nil {
		return err
	}

	userRepo := postgresql.NewUsersRepository(tx)
	credRepo := postgresql.NewCredentialsRepository(tx)


	user.Email, err = userRepo.GetEmailByToken(tok)
	if err != nil {
		return err
	}

	user, err = userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}

    cred, err = credRepo.GetCredentialByUserID(string(user.ID))
	if err != nil {
		return err
	}

	cred.Salt = pwd.Salt(256 - len(request.Password))

	cred.PasswordHash, err = pwd.Hash(request.Password, cred.Salt)
	if err != nil {
		return err
	}

    credRepo.UpdateCredentials(cred)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Send Email Verification Token function must be here!
