package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/samber/lo"
)

func DetectRisk(user models.Account, ip, ua string) bool {
	var secureFactor int64
	if err := database.C.Where(models.AuthTicket{
		AccountID: user.ID,
		IpAddress: ip,
	}).Model(models.AuthTicket{}).Count(&secureFactor).Error; err == nil {
		if secureFactor >= 1 {
			return false
		}
	}

	return true
}

func NewTicket(user models.Account, ip, ua string) (models.AuthTicket, error) {
	var ticket models.AuthTicket
	if err := database.C.Where(models.AuthTicket{
		AccountID: user.ID,
	}).First(&ticket).Error; err == nil {
		return ticket, nil
	}

	ticket = models.AuthTicket{
		Claims:      []string{"*"},
		Audiences:   []string{"passport"},
		IpAddress:   ip,
		UserAgent:   ua,
		RequireMFA:  DetectRisk(user, ip, ua),
		ExpiredAt:   lo.ToPtr(time.Now().Add(2 * time.Hour)),
		AvailableAt: nil,
		AccountID:   user.ID,
	}

	err := database.C.Save(&ticket).Error

	return ticket, err
}

func NewOauthTicket(
	user models.Account,
	client models.ThirdClient,
	claims, audiences []string,
	ip, ua string,
) (models.AuthTicket, error) {
	ticket := models.AuthTicket{
		Claims:       claims,
		Audiences:    audiences,
		IpAddress:    ip,
		UserAgent:    ua,
		RequireMFA:   DetectRisk(user, ip, ua),
		GrantToken:   lo.ToPtr(uuid.NewString()),
		AccessToken:  lo.ToPtr(uuid.NewString()),
		RefreshToken: lo.ToPtr(uuid.NewString()),
		AvailableAt:  lo.ToPtr(time.Now()),
		ExpiredAt:    lo.ToPtr(time.Now()),
		ClientID:     &client.ID,
		AccountID:    user.ID,
	}

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func ActiveTicketWithPassword(ticket models.AuthTicket, password string) (models.AuthTicket, error) {
	if ticket.AvailableAt != nil {
		return ticket, nil
	} else if !ticket.RequireAuthenticate {
		return ticket, fmt.Errorf("detected risk, multi factor authentication required")
	}

	if factor, err := GetPasswordFactor(ticket.AccountID); err != nil {
		return ticket, fmt.Errorf("unable to active ticket: %v", err)
	} else if err = CheckFactor(factor, password); err != nil {
		return ticket, err
	}

	ticket.AvailableAt = lo.ToPtr(time.Now())

	if !ticket.RequireAuthenticate && !ticket.RequireMFA {
		ticket.AvailableAt = lo.ToPtr(time.Now())
	}

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func ActiveTicketWithMFA(ticket models.AuthTicket, factor models.AuthFactor, code string) (models.AuthTicket, error) {
	if ticket.AvailableAt != nil {
		return ticket, nil
	} else if !ticket.RequireMFA {
		return ticket, nil
	}

	if err := CheckFactor(factor, code); err != nil {
		return ticket, fmt.Errorf("invalid code: %v", err)
	}

	ticket.RequireMFA = false

	if !ticket.RequireAuthenticate && !ticket.RequireMFA {
		ticket.AvailableAt = lo.ToPtr(time.Now())
	}

	if err := database.C.Save(&ticket).Error; err != nil {
		return ticket, err
	}

	return ticket, nil
}

func RegenSession(session models.AuthTicket) (models.AuthTicket, error) {
	session.GrantToken = lo.ToPtr(uuid.NewString())
	session.AccessToken = lo.ToPtr(uuid.NewString())
	session.RefreshToken = lo.ToPtr(uuid.NewString())
	err := database.C.Save(&session).Error
	return session, err
}