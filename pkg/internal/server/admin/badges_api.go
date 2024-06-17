package admin

import (
	"fmt"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func grantBadge(c *fiber.Ctx) error {
	if err := utils.CheckPermissions(c, "AdminGrantBadges", true); err != nil {
		return err
	}

	var data struct {
		Type      string         `json:"type" validate:"required"`
		Metadata  map[string]any `json:"metadata"`
		AccountID uint           `json:"account_id"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var account models.Account
	if account, err = services.GetAccount(data.AccountID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("target account was not found: %v", err))
	}

	badge := models.Badge{
		Type:     data.Type,
		Metadata: data.Metadata,
	}

	if err := services.GrantBadge(account, badge); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(badge)
	}
}

func revokeBadge(c *fiber.Ctx) error {
	if err := utils.CheckPermissions(c, "AdminRevokeBadges", true); err != nil {
		return err
	}

	id, _ := c.ParamsInt("badgeId", 0)

	var badge models.Badge
	if err := database.C.Where("id = ?", id).First(&badge).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("target badge was not found: %v", err))
	}

	if err := services.RevokeBadge(badge); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(badge)
	}
}
