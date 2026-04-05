package discovery

type CandidateKind string

const (
	CandidateKindImage   CandidateKind = "image"
	CandidateKindSidecar CandidateKind = "sidecar"
)

type Candidate struct {
	Kind         CandidateKind
	Path         string
	RelativePath string
}

type Warning struct {
	Path    string
	Message string
}

type Result struct {
	Candidates       []Candidate
	Warnings         []Warning
	FilesSeen        int
	DirectoriesSeen  int
	CandidatesFound  int
	WarningsObserved int
}
