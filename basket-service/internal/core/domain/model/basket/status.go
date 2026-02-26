package basket

const (
	StatusEmpty     Status = ""
	StatusCreated   Status = "Created"
	StatusConfirmed Status = "Confirmed"
)

type Status string

func (s Status) Equal(other Status) bool {
	return s == other
}

func (s Status) IsValid() bool {
	return s != StatusEmpty
}

func (s Status) String() string {
	return string(s)
}
