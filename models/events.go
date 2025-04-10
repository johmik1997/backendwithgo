package models

import (
	"database/sql"
	db "john/database"
	"time"
)

type Event struct {
    ID          int           `json:"id"`
    Name        string        `json:"name"`
    Date        time.Time     `json:"date"`
    Location    string        `json:"location"`
    Host        string        `json:"host"`
    Description string        `json:"description"`
    Participants []Participant `json:"participants"`
}

type Participant struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Class string `json:"class"`
}

func GetUpcomingEvents() ([]Event, error) {
    rows, err := db.DB.Query(`
        SELECT e.id, e.name, e.date, e.location, e.host, e.description,
               p.id, p.name, p.class
        FROM events e
        LEFT JOIN participants p ON e.id = p.event_id
        WHERE e.date > NOW()
        ORDER BY e.date ASC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    eventsMap := make(map[int]*Event)
    for rows.Next() {
        var e Event
        // var p Participant
        var pID sql.NullInt64
        var pName, pClass sql.NullString

        err := rows.Scan(
            &e.ID, &e.Name, &e.Date, &e.Location, &e.Host, &e.Description,
            &pID, &pName, &pClass,
        )
        if err != nil {
            return nil, err
        }

        if _, exists := eventsMap[e.ID]; !exists {
            eventsMap[e.ID] = &e
        }

        if pID.Valid {
            eventsMap[e.ID].Participants = append(eventsMap[e.ID].Participants, Participant{
                ID:    int(pID.Int64),
                Name:  pName.String,
                Class: pClass.String,
            })
        }
    }

    var events []Event
    for _, event := range eventsMap {
        events = append(events, *event)
    }
    return events, nil
}