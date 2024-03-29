package services

import (
	"fmt"

	"git.solsynth.dev/hydrogen/identity/pkg/database"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const EmailPasswordTemplate = `Dear %s,

We hope this message finds you well.
As part of our ongoing commitment to ensuring the security of your account, we require you to complete the login process by entering the verification code below:

Your Login Verification Code: %s

Please use the provided code within the next 2 hours to complete your login. 
If you did not request this code, please update your information, maybe your username or email has been leak.

Thank you for your cooperation in helping us maintain the security of your account.

Best regards,
%s`

func LookupFactor(id uint) (models.AuthFactor, error) {
	var factor models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		BaseModel: models.BaseModel{ID: id},
	}).First(&factor).Error

	return factor, err
}

func LookupFactorsByUser(uid uint) ([]models.AuthFactor, error) {
	var factors []models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		AccountID: uid,
	}).Find(&factors).Error

	return factors, err
}

func GetFactorCode(factor models.AuthFactor) (bool, error) {
	switch factor.Type {
	case models.EmailPasswordFactor:
		var user models.Account
		if err := database.C.Where(&models.Account{
			BaseModel: models.BaseModel{ID: factor.AccountID},
		}).Preload("Contacts").First(&user).Error; err != nil {
			return true, err
		}

		factor.Secret = uuid.NewString()[:6]
		if err := database.C.Save(&factor).Error; err != nil {
			return true, err
		}

		subject := fmt.Sprintf("[%s] Login verification code", viper.GetString("name"))
		content := fmt.Sprintf(EmailPasswordTemplate, user.Name, factor.Secret, viper.GetString("maintainer"))
		if err := SendMail(user.GetPrimaryEmail().Content, subject, content); err != nil {
			return true, err
		}
		return true, nil

	default:
		return false, nil
	}
}
