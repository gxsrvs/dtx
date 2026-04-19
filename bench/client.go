package bench

import (
	"time"

	"github.com/gxsrvs/dtx/types"
)

// Profile describes how heavily nullable fields are populated in a
// generated instance.
type Profile int

const (
	// ProfileAllValid populates every nullable field with a value.
	ProfileAllValid Profile = iota
	// ProfileAllNull leaves every nullable field empty / nil.
	ProfileAllNull
	// ProfileMixed populates roughly half of the nullable fields, to
	// approximate a realistic record where some optional information
	// is present and some is missing.
	ProfileMixed
)

// String returns the human-readable profile name used in benchmark
// sub-test labels.
func (p Profile) String() string {
	switch p {
	case ProfileAllValid:
		return "AllValid"
	case ProfileAllNull:
		return "AllNull"
	case ProfileMixed:
		return "Mixed"
	default:
		return "Unknown"
	}
}

// LeadSource is shared by both variants — every field is non-null, so
// there is no representation difference between the dtx and pointer
// styles for this entity.
type LeadSource struct {
	Id   int16  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// -----------------------------------------------------------------------------
// Variant (a): dtx Null* wrappers alongside plain atomic fields.
// -----------------------------------------------------------------------------

// DocumentNull is the dtx variant of Document.
type DocumentNull struct {
	TypeId     int16            `json:"typeId"`
	Number     types.NullString `json:"number,omitzero"`
	DateIssue  types.NullDate   `json:"dateIssue,omitzero"`
	DateExpiry types.NullDate   `json:"dateExpiry,omitzero"`
}

// ClientNull is the dtx variant of Client.
type ClientNull struct {
	Id        int32            `json:"id"`
	FirstName types.NullString `json:"firstName,omitzero"`
	LastName  types.NullString `json:"lastName,omitzero"`
	Source    *LeadSource      `json:"source,omitempty"`
	Document  *DocumentNull    `json:"document,omitempty"`
	BirthDate types.NullDate   `json:"birthDate,omitzero"`
}

// -----------------------------------------------------------------------------
// Variant (b): plain Go pointers for every nullable field.
// -----------------------------------------------------------------------------

// DocumentPtr is the pointer variant of Document.
type DocumentPtr struct {
	TypeId     int16      `json:"typeId"`
	Number     *string    `json:"number,omitempty"`
	DateIssue  *time.Time `json:"dateIssue,omitempty"`
	DateExpiry *time.Time `json:"dateExpiry,omitempty"`
}

// ClientPtr is the pointer variant of Client.
type ClientPtr struct {
	Id        int32        `json:"id"`
	FirstName *string      `json:"firstName,omitempty"`
	LastName  *string      `json:"lastName,omitempty"`
	Source    *LeadSource  `json:"source,omitempty"`
	Document  *DocumentPtr `json:"document,omitempty"`
	BirthDate *time.Time   `json:"birthDate,omitempty"`
}

// -----------------------------------------------------------------------------
// Generators.
// -----------------------------------------------------------------------------

// Sample values are anchored to Yuri Gagarin so the data echoes the
// rest of the project's examples; nothing in the benchmark depends on
// the specific values themselves.
var (
	sampleBirth     = time.Date(1934, 3, 9, 0, 0, 0, 0, time.UTC)
	sampleDocIssue  = time.Date(1957, 9, 1, 0, 0, 0, 0, time.UTC)
	sampleDocExpiry = time.Date(1968, 3, 27, 0, 0, 0, 0, time.UTC)

	sampleSource = LeadSource{
		Id:   7,
		Code: "REFERRAL",
		Name: "Customer referral programme",
	}
)

// NewClientNull builds a ClientNull instance for the given profile.
func NewClientNull(p Profile) ClientNull {
	switch p {
	case ProfileAllValid:
		return ClientNull{
			Id:        42,
			FirstName: types.NewNullString("Yuri"),
			LastName:  types.NewNullString("Gagarin"),
			Source:    new(sampleSource),
			Document: &DocumentNull{
				TypeId:     1,
				Number:     types.NewNullString("VVS-1957-09-01"),
				DateIssue:  types.NewNullDate(sampleDocIssue),
				DateExpiry: types.NewNullDate(sampleDocExpiry),
			},
			BirthDate: types.NewNullDate(sampleBirth),
		}
	case ProfileAllNull:
		return ClientNull{
			Id:        42,
			FirstName: types.NewNullStringEmpty(),
			LastName:  types.NewNullStringEmpty(),
			BirthDate: types.NewNullDateEmpty(),
		}
	case ProfileMixed:
		// FirstName and Source are present; LastName, Document and
		// BirthDate are absent — about half the optional payload.
		return ClientNull{
			Id:        42,
			FirstName: types.NewNullString("Yuri"),
			LastName:  types.NewNullStringEmpty(),
			Source:    new(sampleSource),
			BirthDate: types.NewNullDateEmpty(),
		}
	default:
		panic("bench: unknown profile")
	}
}

// NewClientPtr builds a ClientPtr instance for the given profile.
func NewClientPtr(p Profile) ClientPtr {
	switch p {
	case ProfileAllValid:
		return ClientPtr{
			Id:        42,
			FirstName: new("Yuri"),
			LastName:  new("Gagarin"),
			Source:    new(sampleSource),
			Document: &DocumentPtr{
				TypeId:     1,
				Number:     new("VVS-1957-09-01"),
				DateIssue:  new(sampleDocIssue),
				DateExpiry: new(sampleDocExpiry),
			},
			BirthDate: new(sampleBirth),
		}
	case ProfileAllNull:
		return ClientPtr{
			Id: 42,
		}
	case ProfileMixed:
		return ClientPtr{
			Id:        42,
			FirstName: new("Yuri"),
			Source:    new(sampleSource),
		}
	default:
		panic("bench: unknown profile")
	}
}
