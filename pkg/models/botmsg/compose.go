package botmsg

import (
	"strconv"

	. "personaerpgcompanion/pkg/models"
)

func ComposeWeaponMessage(stats Weapon) string {
	message := ""
	// Name group
	message = stats.Name

	message += "\n"
	for i := 0; i < len(stats.Name); i++ {
		message += "="
	}

	// Main stats
	message += "\nУрон: " + strconv.Itoa(stats.DMG)
	message += "\nСмертельность: " + strconv.Itoa(stats.DLS)

	message += "\n"

	// Skill
	message += "\nНавык: " + stats.Skill
	// Grip
	if stats.Hand1 != "X" {
		if stats.Hand1 == "O" {
			message += "\nОдноручное"
		} else {
			message += "\nВ одной руке: " + stats.Hand1
		}
	}
	if stats.Hand2 != "X" {
		if stats.Hand2 == "O" {
			message += "\nДвуручное"
		} else {
			message += "\nВ двух руках: " + stats.Hand2
		}
	}

	message += "\n"

	// Rarity
	message += "\nРедкость: " + strconv.Itoa(stats.Rarity)
	// Price
	message += "\nЦена: " + strconv.Itoa(stats.Price)
	switch stats.Curr {
	case "z":
		message += " зени (медь)"
		break
	case "b":
		message += " бу (серебро)"
		break
	case "k":
		message += " коку (золото)"
		break
	}

	message += "\n"

	// Qualities
	if stats.Qualities != "-" {
		message += "\nСвойства: " + stats.Qualities
	}
	message += "\n"
	if stats.Additional != "-" {
		message += "\nДополнительно: \n" + stats.Additional
	}

	// Picture
	if len(stats.Pic) > 0 {
		message += "\n" + stats.Pic
	}

	return message
}

func ComposeArmorMessage(stats Armor) string {
	message := ""
	// Name group
	message = stats.Name

	message += "\n"
	for i := 0; i < len(stats.Name); i++ {
		message += "="
	}

	// Main stats
	if stats.Phys+stats.Super > 0 {
		message += "\nЗащита"
		if stats.Phys > 0 {
			message += "\nФизическая: " + strconv.Itoa(stats.Phys)
		}
		if stats.Super > 0 {
			message += "\nСверхъестественная: " + strconv.Itoa(stats.Super)
		}
	}

	message += "\n"

	// Rarity
	message += "\nРедкость: " + strconv.Itoa(stats.Rarity)
	// Price
	message += "\nЦена: " + strconv.Itoa(stats.Price)
	switch stats.Curr {
	case "z":
		message += " зени (медь)"
		break
	case "b":
		message += " бу (серебро)"
		break
	case "k":
		message += " коку (золото)"
		break
	}

	message += "\n"

	// Qualities
	if stats.Qualities != "-" {
		message += "\nСвойства: " + stats.Qualities
	}
	message += "\n"
	if stats.Additional != "-" {
		message += "\nДополнительно: \n" + stats.Additional
	}

	// Picture
	if len(stats.Pic) > 0 {
		message += "\n" + stats.Pic
	}

	return message
}
