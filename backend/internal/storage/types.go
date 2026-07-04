package storage

type ChecksumAlgo string

const (
	ChecksumNone   ChecksumAlgo = "none"
	ChecksumSha256 ChecksumAlgo = "sha256"
)

// ==========options==========
type PutOptions struct {
	ContentType string
	Checksum    ChecksumAlgo
}

type DeleteOptions struct {
	ForceDelete      bool
	GovernanceBypass bool
	VersionID        string
}

// ==========models/DTOs==========
type CompletedPart struct {
	PartNumber    int
	ETag          string
	ChecksumValue *string
	ChecksumAlgo  ChecksumAlgo
}
