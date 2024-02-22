package util

import "time"

const (
	Aquarius    = "水瓶座" //"水瓶座" //	01/20 - 02/18
	Pisces      = "双鱼座" //"双鱼座" //	02/19 - 03/20
	Aries       = "白羊座" //"白羊座" //	03/21 - 04/19
	Taurus      = "金牛座" //"金牛座" //	04/20 - 05/20
	Gemini      = "双子座" //"双子座" //	05/21 - 06/21
	Cancer      = "巨蟹座" //"巨蟹座" //	06/22 - 07/22
	Leo         = "狮子座" //"狮子座" //	07/23 - 08/22
	Virgo       = "处女座" //"处女座" //	08/23 - 09/22
	Libra       = "天秤座" //"天秤座" //	09/23 - 10/23
	Scorpio     = "天蝎座" //"天蝎座" //	10/24 - 11/22
	Sagittarius = "射手座" //"射手座" //	11/23 - 12/21
	Capricorn   = "摩羯座" //"摩羯座" //	12/22 - 01/19
)

func GenConstellation(month time.Month, day int) string {
	switch month {
	case time.January:
		if day < 20 {
			return Capricorn
		}
		return Aquarius
	case time.February:
		if day < 19 {
			return Aquarius
		}
		return Pisces
	case time.March:
		if day < 21 {
			return Pisces
		}
		return Aries
	case time.April:
		if day < 20 {
			return Aries
		}
		return Taurus
	case time.May:
		if day < 21 {
			return Taurus
		}
		return Gemini
	case time.June:
		if day < 22 {
			return Gemini
		}
		return Cancer
	case time.July:
		if day < 23 {
			return Cancer
		}
		return Leo
	case time.August:
		if day < 23 {
			return Leo
		}
		return Virgo
	case time.September:
		if day < 23 {
			return Virgo
		}
		return Libra
	case time.October:
		if day < 24 {
			return Libra
		}
		return Scorpio
	case time.November:
		if day < 23 {
			return Scorpio
		}
		return Sagittarius
	case time.December:
		if day < 22 {
			return Sagittarius
		}
		return Capricorn
	default:
		return ""
	}
}
