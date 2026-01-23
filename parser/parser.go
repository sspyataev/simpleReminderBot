package parser

import (
	"regexp"
	"strings"
)

type ReminderParser struct {
	patterns []*regexp.Regexp
}

func NewReminderParser() *ReminderParser {
	return &ReminderParser{
		// Шаблоны для поиска временных указаний
		patterns: []*regexp.Regexp{
			// через ... (минут, часов, дней)
			regexp.MustCompile(`(?i)\bчерез\s+[\w\s]+?(?:секунд|минут|часов|дней|недель|месяцев|лет|секунду|минуту|час|день|недел[уяь]|месяц|год)\b.*$`),
			
			// в <время> (в 10, в 18:30, в 7 вечера)
			regexp.MustCompile(`(?i)\bв\s+(?:[0-2]?[0-9]:[0-5][0-9]|[\d]+\s+вечера|утра|дня|ночи|час[аов]?\s*\d+|[\d]+\s*(?:утра|вечера|дня|ночи))\b.*$`),
			
			// завтра, сегодня, послезавтра + опционально время
			regexp.MustCompile(`(?i)\b(?:сегодня|завтра|послезавтра|на\s+неделе|на\s+выходных|в\s+понедельник|в\s+вторник|в\s+среду|в\s+четверг|в\s+пятницу|в\s+субботу|в\s+воскресенье)\b.*$`),
			
			// после <события> — можно расширить
			regexp.MustCompile(`(?i)\bпосле\s+[\w\s]+?\b.*$`),
		},
	}
}

func (p *ReminderParser) Parse(input string) (text, timePart string) {
	input = strings.TrimSpace(input)

	for _, re := range p.patterns {
		loc := re.FindStringIndex(input)
		if loc != nil {
			start, end := loc[0], loc[1]
			return strings.TrimSpace(input[:start]), strings.TrimSpace(input[start:end])
		}
	}
	// Если временная часть не найдена
	return input, ""
}

