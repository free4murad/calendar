package app

import (
	"calendar/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"time"
)

func Run(fileBytes []byte, startDate time.Time, endDate time.Time, out io.Writer) {
	var calendars []models.Calendar
	if err := json.Unmarshal(fileBytes, &calendars); err != nil {
		log.Fatalf("Invalid input file: %+v", err)
	}

	slots := getAvailableSlots(calendars, startDate, endDate)
	if err := json.NewEncoder(out).Encode(slots); err != nil {
		log.Fatalf("Failed to write output: %+v", err)
	}
}

func getAvailableSlots(calendars []models.Calendar, startDate time.Time, endDate time.Time) []models.Slot {
	meetings := getMeetingsOfInterest(calendars, startDate, endDate)
	if len(meetings) == 0 {
		return []models.Slot{
			{StartTime: startDate, EndTime: endDate},
		}
	}

	sortMeetings(meetings)
	combineOverlappingMeetings(meetings)

	return getSlots(meetings, startDate, endDate)
}

// getMeetingsOfInterest returns the meetings in the concerned time window.
func getMeetingsOfInterest(calendars []models.Calendar, startDate time.Time, endDate time.Time) []models.Meeting {
	meetings := make([]models.Meeting, 0, len(calendars[0].Meetings))
	for _, calendar := range calendars {
		for _, meeting := range calendar.Meetings {
			if meeting.StartTime.After(endDate) || meeting.EndTime.Before(startDate) {
				continue
			}
			meetings = append(meetings, meeting)
		}
	}
	return meetings
}

// sortMeetings will sort meeting first by startTime and ten by endTome
func sortMeetings(meetings []models.Meeting) {
	//for _, meeting := range meetings {
	//	fmt.Println(meeting.StartTime, "=====", meeting.EndTime, "=======", meeting.Subject)
	//}

	sort.Slice(meetings, func(i, j int) bool {
		if meetings[i].StartTime.Before(meetings[j].StartTime) {
			return true
		}

		if meetings[i].EndTime.Before(meetings[j].EndTime) {
			return true
		}

		return false
	})

	//fmt.Println("--------")
	//for _, meeting := range meetings {
	//	fmt.Println(meeting.StartTime, "=====", meeting.EndTime, "=======", meeting.Subject)
	//}
}

//combineOverlappingMeetings will combine overlapping meetings
func combineOverlappingMeetings(meetings []models.Meeting) {
	//for _, meeting := range meetings {
	//	fmt.Println(meeting.StartTime, "=====", meeting.EndTime, "=======", meeting.Subject)
	//}

	for i := 0; i < len(meetings)-1; i++ {
		currentMeeting := meetings[i]
		nextMeeting := meetings[i+1]

		// Combine if nextMeeting ends before currentMeeting. (This is actually remove on nextMeeting)
		if nextMeeting.EndTime.Before(currentMeeting.EndTime) || nextMeeting.EndTime.Equal(currentMeeting.EndTime) {
			meetings[i].Subject = fmt.Sprintf("%s + %s", currentMeeting.Subject, nextMeeting.Subject)
			meetings = append(meetings[:i+1], meetings[i+2:]...)
			i--
			continue
		}

		// Combine if currentMeeting ends after nextMeeting (This is combination of two meetings)
		if currentMeeting.EndTime.After(nextMeeting.StartTime) || currentMeeting.EndTime.Equal(nextMeeting.StartTime) {
			meetings[i].EndTime = nextMeeting.EndTime
			meetings[i].Subject = fmt.Sprintf("%s + %s", currentMeeting.Subject, nextMeeting.Subject)
			meetings = append(meetings[:i+1], meetings[i+2:]...)
			i--
			continue
		}
	}

	//fmt.Println("--------")
	//for _, meeting := range meetings {
	//	fmt.Println(meeting.StartTime, "=====", meeting.EndTime, "=======", meeting.Subject)
	//}
}

// getSlots will form the final array of models.Slot that are available
func getSlots(meetings []models.Meeting, startTime time.Time, endTime time.Time) []models.Slot {
	availableSlots := make([]models.Slot, 0, len(meetings)+1)

	// Add slot before the first meeting
	if meetings[0].StartTime.After(startTime) {
		availableSlots = append(availableSlots, models.Slot{
			StartTime: startTime.UTC(),
			EndTime:   meetings[0].StartTime,
		})
	}

	// Add slots between meetings
	for i := range meetings {
		if i+1 == len(meetings) {
			break
		}

		availableSlots = append(availableSlots, models.Slot{
			StartTime: meetings[i].EndTime.UTC(),
			EndTime:   meetings[i+1].StartTime.UTC(),
		})
	}

	// Add slot after the last meeting
	if meetings[len(meetings)-1].EndTime.Before(endTime) {
		availableSlots = append(availableSlots, models.Slot{
			StartTime: meetings[len(meetings)-1].EndTime.UTC(),
			EndTime:   endTime.UTC(),
		})
	}

	return availableSlots
}
