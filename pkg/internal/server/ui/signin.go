package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/sujit-baniya/flash"
)

func signinPage(c *fiber.Ctx) error {
	localizer := c.Locals("localizer").(*i18n.Localizer)

	next, _ := localizer.LocalizeMessage(&i18n.Message{ID: "next"})
	username, _ := localizer.LocalizeMessage(&i18n.Message{ID: "username"})
	password, _ := localizer.LocalizeMessage(&i18n.Message{ID: "password"})
	signup, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signupTitle"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinCaption"})
	requiredNotify, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinRequired"})

	var info any
	if flash.Get(c)["message"] != nil {
		info = flash.Get(c)["message"]
	} else {
		info = requiredNotify
	}

	return c.Render("views/signin", fiber.Map{
		"info": info,
		"i18n": fiber.Map{
			"next":     next,
			"username": username,
			"password": password,
			"signup":   signup,
			"title":    title,
			"caption":  caption,
		},
	}, "views/layouts/auth")
}

func signinAction(c *fiber.Ctx) error {
	var data struct {
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/sign-in")
	}

	user, err := services.LookupAccount(data.Username)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("account was not found: %v", err.Error()),
		}).Redirect("/sign-in")
	}

	ticket, err := services.NewTicket(user, c.IP(), c.Get(fiber.HeaderUserAgent))
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable setup ticket: %v", err.Error()),
		}).Redirect("/sign-in")
	}

	ticket, err = services.ActiveTicketWithPassword(ticket, data.Password)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("invalid password: %v", err.Error()),
		}).Redirect("/sign-in")
	}

	if ticket.IsAvailable() != nil {
		return flash.WithData(c, fiber.Map{
			"redirect_uri": utils.GetRedirectUri(c),
		}).Redirect(fmt.Sprintf("/mfa?ticket=%d", ticket.ID))
	}

	access, refresh, err := services.ExchangeToken(*ticket.GrantToken)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("failed to exchange token: %v", err.Error()),
		}).Redirect("/sign-in")
	} else {
		services.SetJwtCookieSet(c, access, refresh)
	}

	return c.Redirect(lo.FromPtr(utils.GetRedirectUri(c, "/users/me")))
}
