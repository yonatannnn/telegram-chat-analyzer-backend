package usecase

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"telegram-chat-analyzer/internal/domain"
	"time"
)

type MessageUsecase interface {
	SeparateMessagesByPerson(chat domain.Chat) (map[string][]domain.Message, string, string)
	CountWords(chat domain.Chat) (map[string]map[string]int, error)
	TopSixFrequentWords(chat domain.Chat, combinedWordCount map[string]int, personOneWordCount map[string]int, personTwoWordCount map[string]int) map[string]map[string]int
	CountMessages(
		chat domain.Chat,
	) (int, int, int)
	GetPersons(chat domain.Chat) (string, string)
	TotalDaysTalked(chat domain.Chat) int
	MessagesPerDay(chat domain.Chat) map[string]map[string]int
	WeeklyStats(chat domain.Chat) map[string]map[string]int
	HourlyStats(chat domain.Chat) map[string]map[string]int
	MostActiveDayOfWeek(chat domain.Chat) map[string]string
	MessageLengthStatistics(chat domain.Chat) map[string]map[string]float64
	ReplyTimeAnalysis(chat domain.Chat) map[string]float64
	CountConversationStartersPerDay(chat domain.Chat) (map[string]int, error)
	CountConsecutiveDays(chat domain.Chat) (map[string][]interface{}, error)
	GetSharedInterests(chat domain.Chat) []string
	AverageMessagesPerDay(chat domain.Chat) map[string]float64
	CountWord(chat domain.Chat) (map[string]int, map[string]int, error)
	RelationshipScore(chat domain.Chat) (float64, error)
	CurrentStreak(chat domain.Chat) (map[string][]interface{}, error)
}

type messageUsecase struct{}

func NewMessageUsecase() MessageUsecase {
	return &messageUsecase{}
}

func (u *messageUsecase) GetPersons(chat domain.Chat) (string, string) {
	_, personOne, personTwo := u.SeparateMessagesByPerson(chat)
	return personOne, personTwo
}

func (u *messageUsecase) SeparateMessagesByPerson(chat domain.Chat) (map[string][]domain.Message, string, string) {
	messagesByPerson := make(map[string][]domain.Message)
	personTwo := strings.TrimSpace(strings.ToLower(chat.Name))
	personTwoID := "user" + strconv.Itoa(chat.ID)
	var personOne string
	for _, message := range chat.Messages {
		sender := message.From
		senderID := message.FromID
		if senderID != personTwoID {
			personOne = sender
		} else {
			personTwo = sender
			personTwoID = senderID
		}

		messagesByPerson[sender] = append(messagesByPerson[sender], message)
	}
	return messagesByPerson, personOne, personTwo

}

func (u *messageUsecase) CountMessages(chat domain.Chat) (int, int, int) {
	messageByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)
	personOneMessageCount := len(messageByPerson[personOne])
	personTwoMessageCount := len(messageByPerson[personTwo])
	totalMessageCount := personOneMessageCount + personTwoMessageCount
	return totalMessageCount, personOneMessageCount, personTwoMessageCount
}

func (u *messageUsecase) TopSixFrequentWords(
	chat domain.Chat,
	combinedWordCount map[string]int,
	personOneWordCount map[string]int,
	personTwoWordCount map[string]int,
) map[string]map[string]int {
	type wordFrequency struct {
		Word       string
		TotalCount int
	}

	var wordFrequencies []wordFrequency
	for word := range combinedWordCount {
		totalCount := personOneWordCount[word] + personTwoWordCount[word]
		wordFrequencies = append(wordFrequencies, wordFrequency{Word: word, TotalCount: totalCount})
	}

	sort.Slice(wordFrequencies, func(i, j int) bool {
		return wordFrequencies[i].TotalCount > wordFrequencies[j].TotalCount
	})

	// Step 3: Extract the top 3 words
	personOne, personTwo := u.GetPersons(chat)
	topWords := map[string]map[string]int{}
	for i := 0; i < len(wordFrequencies) && i < 6; i++ {
		word := wordFrequencies[i].Word
		topWords[word] = map[string]int{
			personOne: personOneWordCount[word],
			personTwo: personTwoWordCount[word],
		}
	}

	return topWords
}

func (u *messageUsecase) CountWords(chat domain.Chat) (map[string]map[string]int, error) {
	combinedWordCount := make(map[string]int)
	personOneWordCount := make(map[string]int)
	personTwoWordCount := make(map[string]int)

	messagesByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	cleanWord := func(word string) string {
		return strings.ToLower(strings.Trim(word, ".,!\""))
	}

	for person, messages := range messagesByPerson {
		for _, msg := range messages {
			text, ok := msg.Text.(string)
			if !ok {
				continue
			}
			words := strings.Fields(text)
			for _, word := range words {
				cleanedWord := cleanWord(word)
				if cleanedWord == "" {
					continue
				}
				combinedWordCount[cleanedWord]++
				if person == personOne {
					personOneWordCount[cleanedWord]++
				} else if person == personTwo {
					personTwoWordCount[cleanedWord]++
				}
			}
		}
	}

	topWords := u.TopSixFrequentWords(chat, combinedWordCount, personOneWordCount, personTwoWordCount)

	return topWords, nil
}

func splitText(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'))
	})
}

func (u *messageUsecase) CountWord(chat domain.Chat) (map[string]int, map[string]int, error) {
	wordCount := make(map[string]int)

	messagesByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)
	average1 := 0
	average2 := 0

	for person, messages := range messagesByPerson {
		for _, msg := range messages {
			text, ok := msg.Text.(string)
			if !ok {
				continue
			}

			words := splitText(text)
			length := len(words)
			wordCount[person] += length
		}
	}

	average1 = wordCount[personOne] / len(messagesByPerson[personOne])
	average2 = wordCount[personTwo] / len(messagesByPerson[personTwo])

	averages := map[string]int{
		personOne: average1,
		personTwo: average2,
	}

	return wordCount, averages, nil
}

func (u *messageUsecase) TotalDaysTalked(chat domain.Chat) int {
	messages := chat.Messages
	dateSet := make(map[string]struct{})
	for _, message := range messages {
		date := strings.Split(message.Date, "T")[0]
		dateSet[date] = struct{}{}
	}
	return len(dateSet)
}

func (u *messageUsecase) MessagesPerDay(chat domain.Chat) map[string]map[string]int {
	messages := chat.Messages
	result := make(map[string]map[string]int)

	_, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	for _, message := range messages {
		date := strings.Split(message.Date, "T")[0]

		if _, exists := result[date]; !exists {
			result[date] = map[string]int{
				personOne: 0,
				personTwo: 0,
			}
		}

		if message.From == personOne {
			result[date][personOne]++
		} else if message.From == personTwo {
			result[date][personTwo]++
		}
	}

	return result
}

func (u *messageUsecase) WeeklyStats(chat domain.Chat) map[string]map[string]int {
	_, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	messages := chat.Messages
	result := map[string]map[string]int{
		"monday":    {personOne: 0, personTwo: 0},
		"tuesday":   {personOne: 0, personTwo: 0},
		"wednesday": {personOne: 0, personTwo: 0},
		"thursday":  {personOne: 0, personTwo: 0},
		"friday":    {personOne: 0, personTwo: 0},
		"saturday":  {personOne: 0, personTwo: 0},
		"sunday":    {personOne: 0, personTwo: 0},
	}

	for _, message := range messages {
		// Parse the date from the message
		dateStr := strings.Split(message.Date, "T")[0]
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // Skip invalid dates
		}

		// Get the day of the week
		dayOfWeek := parsedDate.Weekday().String()

		// Convert to lowercase to match the result map keys
		dayOfWeek = strings.ToLower(dayOfWeek)

		// Increment the count for the respective person
		if message.From == personOne {
			result[dayOfWeek][personOne]++
		} else if message.From == personTwo {
			result[dayOfWeek][personTwo]++
		}
	}

	return result
}

func (u *messageUsecase) HourlyStats(chat domain.Chat) map[string]map[string]int {
	messages := chat.Messages
	result := make(map[string]map[string]int)

	_, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	// Initialize the result map with all 24 hours
	for hour := 0; hour < 24; hour++ {
		hourStr := formatHour(hour)
		result[hourStr] = map[string]int{
			personOne: 0,
			personTwo: 0,
		}
	}

	for _, message := range messages {
		dateTimeParts := strings.Split(message.Date, "T")
		if len(dateTimeParts) < 2 {
			continue
		}
		timePart := dateTimeParts[1]
		parsedTime, err := time.Parse("15:04:05", timePart)
		if err != nil {
			continue
		}
		hour := parsedTime.Hour()
		hourStr := formatHour(hour)
		if message.From == personOne {
			result[hourStr][personOne]++
		} else if message.From == personTwo {
			result[hourStr][personTwo]++
		}
	}

	return result
}

func formatHour(hour int) string {
	return time.Date(0, 0, 0, hour, 0, 0, 0, time.UTC).Format("15")
}

func (u *messageUsecase) MostActiveDayOfWeek(chat domain.Chat) map[string]string {
	personOneCount := make(map[string]int)
	personTwoCount := make(map[string]int)
	overallCount := make(map[string]int)
	_, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	getWeekday := func(dateStr string) string {
		date, err := time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return ""
		}
		return date.Weekday().String()
	}

	for _, message := range chat.Messages {
		weekday := getWeekday(message.Date)
		if weekday == "" {
			continue
		}

		// Count for person and overall
		if message.From == personOne {
			personOneCount[weekday]++
		} else if message.From == personTwo {
			personTwoCount[weekday]++
		}
		overallCount[weekday]++
	}

	findMostActiveDay := func(counts map[string]int) string {
		var mostActiveDay string
		var maxCount int
		for day, count := range counts {
			if count > maxCount {
				maxCount = count
				mostActiveDay = day
			}
		}
		return mostActiveDay
	}

	// Find the most active day for each
	return map[string]string{
		personOne: findMostActiveDay(personOneCount),
		personTwo: findMostActiveDay(personTwoCount),
		"overall": findMostActiveDay(overallCount),
	}
}

func (u *messageUsecase) MessageLengthStatistics(chat domain.Chat) map[string]map[string]float64 {
	messageByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	total1 := 0
	total2 := 0
	max1 := 0
	max2 := 0
	min1 := 10000000
	min2 := 10000000
	messages1 := messageByPerson[personOne]
	messages2 := messageByPerson[personTwo]

	for _, message1 := range messages1 {
		text, ok := message1.Text.(string)
		if !ok {
			continue
		}
		length := len(text)
		total1 += length
		if length > max1 {
			max1 = length
		}
		if length < min1 {
			min1 = length
		}
	}

	average1 := float64(total1) / float64(len(messages1))

	for _, message2 := range messages2 {
		text, ok := message2.Text.(string)
		if !ok {
			continue
		}
		length := len(text)
		total2 += length
		if length > max2 {
			max2 = length
		}
		if length < min2 {
			min2 = length
		}
	}

	personOneData := map[string]float64{
		"total":   float64(total1),
		"max":     float64(max1),
		"min":     float64(min1),
		"average": average1,
	}

	average2 := float64(total2) / float64(len(messages2))

	personTwoData := map[string]float64{
		"total":   float64(total2),
		"max":     float64(max2),
		"min":     float64(min2),
		"average": float64(average2),
	}

	// Calculate stats for each person
	return map[string]map[string]float64{
		personOne: personOneData,
		personTwo: personTwoData,
	}
}

func (u *messageUsecase) ReplyTimeAnalysis(chat domain.Chat) map[string]float64 {
	parseTime := func(date string) time.Time {
		t, err := time.Parse("2006-01-02T15:04:05", date)
		if err != nil {
			return time.Time{}
		}
		return t
	}

	// Group messages by day
	messagesByDay := make(map[string][]domain.Message)
	for _, message := range chat.Messages {
		date := parseTime(message.Date).Format("2006-01-02")
		messagesByDay[date] = append(messagesByDay[date], message)
	}

	// Filter function to exclude sleep hours
	isDuringSleepHours := func(t time.Time) bool {
		hour := t.Hour()
		return hour < 4 || hour >= 23 // Sleep hours: 11 PM to 10 PM
	}

	// Analyze reply times for each day
	replyTimeStatsByDay := make(map[string]map[string]float64)
	var totalReplyTimes []float64

	for day, messages := range messagesByDay {
		var replyTimes []float64

		for i := 1; i < len(messages); i++ {
			currMsg := messages[i]
			prevMsg := messages[i-1]

			// Only calculate reply time if it's between different people
			if currMsg.From == prevMsg.From {
				continue
			}

			currTime := parseTime(currMsg.Date)
			prevTime := parseTime(prevMsg.Date)

			// Skip if reply time is during sleep hours
			if isDuringSleepHours(prevTime) || isDuringSleepHours(currTime) {
				continue
			}

			replyTime := currTime.Sub(prevTime).Minutes()
			if replyTime > 300 {
				continue
			}
			replyTimes = append(replyTimes, replyTime)
			totalReplyTimes = append(totalReplyTimes, replyTime)
		}
		replyTimeStatsByDay[day] = calculateStats(replyTimes)
	}

	// Calculate overall statistics
	totalStats := calculateStats(totalReplyTimes)

	return totalStats
}

func calculateStats(replyTimes []float64) map[string]float64 {
	if len(replyTimes) == 0 {
		return map[string]float64{"average": 0, "min": 0, "max": 0}
	}

	sum := 0.0
	min := replyTimes[0]
	max := replyTimes[0]

	for _, rt := range replyTimes {
		sum += rt
		if rt < min {
			min = rt
		}
		if rt > max {
			max = rt
		}
	}

	average := sum / float64(len(replyTimes))
	return map[string]float64{"average": average, "min": min, "max": max}
}

func (u *messageUsecase) CountConversationStartersPerDay(chat domain.Chat) (map[string]int, error) {
	messagesByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	conversationStarters := map[string]int{
		personOne: 0,
		personTwo: 0,
	}

	getDate := func(dateUnixtime string) string {
		dateParts := strings.Split(dateUnixtime, "T")
		return dateParts[0]
	}
	firstMessagesOfDay := make(map[string]bool)
	// Iterate over each person's messages
	for person, messages := range messagesByPerson {

		for _, msg := range messages {
			date := getDate(msg.Date)
			if _, exists := firstMessagesOfDay[date]; !exists {
				firstMessagesOfDay[date] = true
				conversationStarters[person]++
			}
		}
	}

	return conversationStarters, nil
}

func (u *messageUsecase) CountConsecutiveDays(chat domain.Chat) (map[string][]interface{}, error) {
	messageByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)
	consecutiveDays := map[string][]interface{}{
		personOne: {0, "", ""},
		personTwo: {0, "", ""},
		"overall": {0, "", ""},
	}

	messages := chat.Messages
	sort.Slice(messages, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02T15:04:05", messages[i].Date)
		dateJ, _ := time.Parse("2006-01-02T15:04:05", messages[j].Date)
		return dateI.Before(dateJ)
	})

	prevDate := ""
	days := make(map[string]bool)
	start := 0
	end := 0
	startDate := ""
	endDate := ""
	number := 0

	for _, message := range messages {
		date := message.Date
		realDate := strings.Split(date, "T")[0]
		if _, exists := days[realDate]; !exists {
			days[realDate] = true
		} else {
			continue
		}

		if prevDate == "" {
			startDate = realDate
		}

		if prevDate != "" && isConsecutive(prevDate, date) {
			end += 1
			endDate = realDate
		} else {

			number = end - start + 1
			if number > consecutiveDays["overall"][0].(int) {
				consecutiveDays["overall"][0] = number
				consecutiveDays["overall"][1] = startDate
				consecutiveDays["overall"][2] = endDate
			}
			start = 0
			end = 0
			startDate = realDate
		}
		prevDate = date
	}

	number = end - start + 1
	if number >= consecutiveDays["overall"][0].(int) {
		consecutiveDays["overall"][0] = number
		consecutiveDays["overall"][1] = startDate
		consecutiveDays["overall"][2] = endDate
	}

	for person := range messageByPerson {
		if _, exists := consecutiveDays[person]; !exists {
			// Initialize the entry if it doesn't exist
			consecutiveDays[person] = []interface{}{0, "", ""}
		}

		messages := messageByPerson[person]
		sort.Slice(messages, func(i, j int) bool {
			dateI, _ := time.Parse("2006-01-02T15:04:05", messages[i].Date)
			dateJ, _ := time.Parse("2006-01-02T15:04:05", messages[j].Date)
			return dateI.Before(dateJ)
		})

		prevDate := ""
		days := make(map[string]bool)

		start := 0
		end := 0
		startDate := ""
		endDate := ""
		number := 0

		for _, message := range messages {
			date := message.Date
			realDate := strings.Split(date, "T")[0]

			if _, exists := days[realDate]; !exists {
				days[realDate] = true
			} else {
				continue
			}

			if prevDate == "" {
				startDate = realDate
			}

			if prevDate != "" && isConsecutive(prevDate, date) {
				end += 1
				endDate = realDate
			} else {
				number = end - start + 1
				if number > consecutiveDays[person][0].(int) {
					consecutiveDays[person][0] = number
					consecutiveDays[person][1] = startDate
					consecutiveDays[person][2] = endDate
				}

				start = 0
				end = 0
				startDate = realDate
			}
			prevDate = date
		}

		number = end - start + 1
		if number >= consecutiveDays[person][0].(int) {
			consecutiveDays[person][0] = number
			consecutiveDays[person][1] = startDate
			consecutiveDays[person][2] = endDate
		}
	}

	return consecutiveDays, nil
}

func findMax(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	max := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}

func isConsecutive(date1, date2 string) bool {
	layout := "2006-01-02"
	t1, err1 := time.Parse(layout, strings.Split(date1, "T")[0])
	t2, err2 := time.Parse(layout, strings.Split(date2, "T")[0])

	if err1 != nil || err2 != nil {
		return false
	}
	return t2.Sub(t1).Hours() <= 24
}

func (u *messageUsecase) CurrentStreak(chat domain.Chat) (map[string][]interface{}, error) {
	consecutiveDays := map[string][]interface{}{
		"overall": {0, "", ""},
	}

	messages := chat.Messages
	sort.Slice(messages, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02T15:04:05", messages[i].Date)
		dateJ, _ := time.Parse("2006-01-02T15:04:05", messages[j].Date)
		return dateI.Before(dateJ)
	})

	today := time.Now().Format("2006-01-02")
	if len(messages) == 0 || strings.Split(messages[len(messages)-1].Date, "T")[0] != today {
		return consecutiveDays, nil
	}

	prevDate := ""
	days := make(map[string]bool)
	start := 0
	end := 0
	startDate := ""
	endDate := ""
	number := 0
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		date := message.Date
		realDate := strings.Split(date, "T")[0]
		if _, exists := days[realDate]; !exists {
			days[realDate] = true
		} else {
			continue
		}

		if prevDate == "" {
			startDate = realDate
		}

		if prevDate != "" && isConsecutive(prevDate, date) {
			end += 1
			endDate = realDate
		} else {
			break
		}
		prevDate = date
	}

	number = end - start + 1
	if number > consecutiveDays["overall"][0].(int) {
		consecutiveDays["overall"][0] = number
		consecutiveDays["overall"][1] = startDate
		consecutiveDays["overall"][2] = endDate
	}
	return consecutiveDays, nil
}

func countWords2(messages []domain.Message) map[string]int {
	wordCount := make(map[string]int)
	cleanWord := func(word string) string {
		return strings.ToLower(strings.Trim(word, ".,!\""))
	}

	for _, msg := range messages {
		text, ok := msg.Text.(string)
		if !ok {
			continue
		}
		words := strings.Fields(text)
		for _, word := range words {
			cleanedWord := cleanWord(word)
			if cleanedWord == "" {
				continue
			}
			wordCount[cleanedWord]++
		}
	}

	return wordCount
}

func sharedInterests(personOneWordCount, personTwoWordCount map[string]int) map[string]int {
	shared := make(map[string]int)

	for word, countOne := range personOneWordCount {
		if countTwo, exists := personTwoWordCount[word]; exists {
			shared[word] = countOne + countTwo // You could also combine frequencies in a different way.
		}
	}

	return shared
}

func sortSharedInterests(shared map[string]int) []string {
	type pair struct {
		word  string
		count int
	}

	var sorted []pair
	for word, count := range shared {
		sorted = append(sorted, pair{word, count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	// Returning just the words sorted by frequency
	var result []string
	for _, p := range sorted {
		result = append(result, p.word)
	}

	return result
}

func (u *messageUsecase) GetSharedInterests(chat domain.Chat) []string {
	messagesByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	personOneWordCount := countWords2(messagesByPerson[personOne])
	personTwoWordCount := countWords2(messagesByPerson[personTwo])

	shared := sharedInterests(personOneWordCount, personTwoWordCount)
	return sortSharedInterests(shared)
}

func (u *messageUsecase) AverageMessagesPerDay(chat domain.Chat) map[string]float64 {
	messagesByPerson, personOne, personTwo := u.SeparateMessagesByPerson(chat)

	personOneMessages := len(messagesByPerson[personOne])
	personTwoMessages := len(messagesByPerson[personTwo])

	totalMessages := personOneMessages + personTwoMessages

	totalDays := u.TotalDaysTalked(chat)

	averageMessagesPerDay := map[string]float64{
		personOne: float64(personOneMessages) / float64(totalDays),
		personTwo: float64(personTwoMessages) / float64(totalDays),
		"overall": float64(totalMessages) / float64(totalDays),
	}

	return averageMessagesPerDay
}

func (u *messageUsecase) RelationshipScore(chat domain.Chat) (float64, error) {
	// Get basic details
	personOne, personTwo := u.GetPersons(chat)

	totalMessages, personOneMessages, personTwoMessages := u.CountMessages(chat)
	per1 := float64(personOneMessages) * 100 / float64(totalMessages)
	per2 := float64(personTwoMessages) * 100 / float64(totalMessages)
	fmt.Println("per1", per1)
	fmt.Println("per2", per2)
	fmt.Println(" ")

	s1 := 0.0
	if per1 < per2 {
		s1 = (1 - (per1 / per2)) * 15
	} else {
		s1 = (1 - (per2 / per1)) * 15
	}

	totalDaysTalked := u.TotalDaysTalked(chat)
	consecutiveDays, err := u.CountConsecutiveDays(chat)
	if err != nil {
		return 0, fmt.Errorf("failed to count consecutive days: %v", err)
	}

	personOneConsecutiveDays := consecutiveDays[personOne][0].(int)
	personTwoConsecutiveDays := consecutiveDays[personTwo][0].(int)
	overallConsecutiveDays := consecutiveDays["overall"][0].(int)

	f1 := float64(personOneConsecutiveDays) / float64(overallConsecutiveDays)
	f2 := float64(personTwoConsecutiveDays) / float64(overallConsecutiveDays)

	fmt.Println("personOneConsecutiveDays", personOneConsecutiveDays)
	fmt.Println("personTwoConsecutiveDays", personTwoConsecutiveDays)
	fmt.Println("overall", overallConsecutiveDays)
	fmt.Println(" ")
	dif1 := math.Abs(f1-f2) * 10
	dif2 := (math.Abs(1-f2) + math.Abs(1-f1)) * 5
	s2 := dif1 + dif2

	s3 := 0.0
	if totalDaysTalked > 365 {
		s3 = 0
	} else {
		s3 = (1 - (float64(overallConsecutiveDays) / float64(totalDaysTalked))) * 2
	}

	acvtiveDay := u.MostActiveDayOfWeek(chat)
	personOneActiveDay := acvtiveDay[personOne]
	personTwoAvtiveDay := acvtiveDay[personTwo]
	fmt.Println("personOneActiveDays", personOneActiveDay)
	fmt.Println("personTwoActiveDays", personTwoAvtiveDay)
	fmt.Println(" ")

	s4 := 0.0
	if personOneActiveDay == personTwoAvtiveDay {
		s4 = 1
	}

	replyTime := u.ReplyTimeAnalysis(chat)
	averageReplyTime := replyTime["average"]
	s5 := 0.25 * averageReplyTime
	fmt.Println("averageReplyTime", averageReplyTime)

	_, average, err := u.CountWord(chat)

	if err != nil {
		return 0, fmt.Errorf("failed to count words: %v", err)
	}

	personOneWordCount := average[personOne]
	personTwoWordCount := average[personTwo]

	fmt.Println("personOneWordCount", personOneWordCount)
	fmt.Println("personTwoWordCount", personTwoWordCount)
	fmt.Println(" ")

	s6 := 0.0
	if personOneWordCount < personTwoWordCount {
		s6 = (1 - (float64(personOneWordCount) / float64(personTwoWordCount))) * 5
	} else {
		s6 = (1 - (float64(personTwoWordCount) / float64(personOneWordCount))) * 5
	}

	averageMessages := u.AverageMessagesPerDay(chat)
	personOneAverageMessages := averageMessages[personOne]
	personTwoAverageMessages := averageMessages[personTwo]
	fmt.Println("personOneAverageMessages", personOneAverageMessages)
	fmt.Println("personTwoAverageMessages", personTwoAverageMessages)
	total := personOneAverageMessages + personTwoAverageMessages
	p1 := personOneAverageMessages / total
	p2 := personTwoAverageMessages / total
	pd := math.Abs(p1 - p2)

	fmt.Println(" ")
	s7 := pd * 20

	fmt.Println(s1, s2, s3, s4, s5, s6, s7)

	relationshipScore := 100 - (s1 + s2 + s3 - s4 + s5 + s6 + s7)
	return math.Round(relationshipScore), nil

}
