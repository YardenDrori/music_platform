package storage

type ChecksumAlgo string

const (
	ChecksumNone   ChecksumAlgo = "none"
	ChecksumSha256 ChecksumAlgo = "sha256"
)

type PutOptions struct {
	ContentType string
	Checksum    ChecksumAlgo
}

type DeleteOptions struct {
	ForceDelete      bool
	GovernanceBypass bool
	VersionID        string
}

type CompletedPart struct {
	PartNumber    int
	ETag          string
	ChecksumValue string
	ChecksumAlgo  ChecksumAlgo
}
