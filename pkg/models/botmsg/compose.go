package msg

import (
	"strconv"

	. "personaerpgcompanion/pkg/models"
)

func WelcomeMessage() string {

	msg := "Вас приветствует бот Personae — ваш компаньон в мире НРИ!"
	msg += "\n\nСейчас бот находится на стадии разработки, а функционал непрерывно дорабатывается. "
	msg += "Чтобы предложить свои идеи или обозначить ошибку в работе бота, свяжитесь с разработчиком @gewiahr."
	msg += "\n\nДля того чтобы узнать текущий функционал введите команду /help."

	return msg
}

func CommandNotFoundMessage() string {

	msg := "Команда не опознана!"
	msg += "\nВоспользуйтесь командой /help чтобы узнать о командах бота."

	return msg
}

func HelpMessage() string {

	msg := "⚔️ Для поиска оружия используйте команду /w и название оружия."
	msg += "\n👘 Для поиска брони и одежды используйте команду /a и название брони или одежды."
	msg += "\n"
	msg += "\nПока что бот понимает только оригинальные названия на английском языке."

	return msg
}

func ComposeWeaponMessage(stats Weapon) string {
	msg := ""
	// Name group
	msg = stats.Name

	msg += "\n"
	for i := 0; i < len(stats.Name); i++ {
		msg += "="
	}

	// Main stats
	msg += "\nУрон: " + strconv.Itoa(stats.DMG)
	msg += "\nСмертельность: " + strconv.Itoa(stats.DLS)
	// Range
	msg += "\nДистанция: " + stats.RNG

	msg += "\n"

	// Skill
	msg += "\nНавык: " + stats.Skill
	// Grip
	if stats.Hand1 != "X" {
		if stats.Hand1 == "O" {
			msg += "\nОдноручное"
		} else {
			msg += "\nВ одной руке: " + stats.Hand1
		}
	}
	if stats.Hand2 != "X" {
		if stats.Hand2 == "O" {
			msg += "\nДвуручное"
		} else {
			msg += "\nВ двух руках: " + stats.Hand2
		}
	}

	msg += "\n"

	// Rarity
	msg += "\nРедкость: " + strconv.Itoa(stats.Rarity)
	// Price
	msg += "\nЦена: " + strconv.Itoa(stats.Price)
	switch stats.Curr {
	case "z":
		msg += " зени (медь)"
		break
	case "b":
		msg += " бу (серебро)"
		break
	case "k":
		msg += " коку (золото)"
		break
	}

	msg += "\n"

	// Qualities
	if stats.Qualities != "-" {
		msg += "\nСвойства: " + stats.Qualities
	}
	msg += "\n"
	if stats.Additional != "-" {
		msg += "\nДополнительно: \n" + stats.Additional
	}

	// Picture
	if len(stats.Pic) > 0 {
		msg += "\n" + stats.Pic
	}

	return msg
}

func ComposeArmorMessage(stats Armor) string {
	msg := ""
	// Name group
	msg = stats.Name

	msg += "\n"
	for i := 0; i < len(stats.Name); i++ {
		msg += "="
	}

	// Main stats
	if stats.Phys+stats.Super > 0 {
		msg += "\nЗащита"
		if stats.Phys > 0 {
			msg += "\nФизическая: " + strconv.Itoa(stats.Phys)
		}
		if stats.Super > 0 {
			msg += "\nСверхъестественная: " + strconv.Itoa(stats.Super)
		}
	}

	msg += "\n"

	// Rarity
	msg += "\nРедкость: " + strconv.Itoa(stats.Rarity)
	// Price
	msg += "\nЦена: " + strconv.Itoa(stats.Price)
	switch stats.Curr {
	case "z":
		msg += " зени (медь)"
		break
	case "b":
		msg += " бу (серебро)"
		break
	case "k":
		msg += " коку (золото)"
		break
	}

	msg += "\n"

	// Qualities
	if stats.Qualities != "-" {
		msg += "\nСвойства: " + stats.Qualities
	}
	msg += "\n"
	if stats.Additional != "-" {
		msg += "\nДополнительно: \n" + stats.Additional
	}

	// Picture
	if len(stats.Pic) > 0 {
		msg += "\n" + stats.Pic
	}

	return msg
}

func ComposeQualityMessage(stats Quality) string {
	msg := ""
	// Name group
	msg = stats.Name

	msg += "\n"
	for i := 0; i < len(stats.Name); i++ {
		msg += "="
	}

	msg += "\n"

	// Main
	msg += stats.General

	return msg
}
