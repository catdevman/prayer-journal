package models

import (
	"time"
)

// PrayerStatus defines the state of a prayer
type PrayerStatus string

const (
	StatusActive   PrayerStatus = "ACTIVE"
	StatusAnswered PrayerStatus = "ANSWERED"
	StatusArchived PrayerStatus = "ARCHIVED"
)

type Prayer struct {
	ID        string    `json:"id" dynamodbav:"id"`
	UserID    string    `json:"userId" dynamodbav:"pk"`    // Partition Key: The user who owns this journal entry
	CreatedAt time.Time `json:"createdAt" dynamodbav:"sk"` // Sort Key: For easy sorting by date
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updated_at"`

	Title  string       `json:"title" dynamodbav:"title"`
	Body   string       `json:"body" dynamodbav:"body"`
	Status PrayerStatus `json:"status" dynamodbav:"status"`

	// Who is this prayer for? (e.g., "Aunt Sally", "My Country", "Myself")
	Target string `json:"target" dynamodbav:"target"`

	// Sharing Metadata
	// If this prayer was imported from a shared link, this field contains the contact info/name of the sharer.
	SharedBy string `json:"sharedBy,omitempty" dynamodbav:"shared_by,omitempty"`
	IsShared bool   `json:"isShared" dynamodbav:"is_shared"`
}
