package main

type Wt_appointment struct {
	Locator          []byte
	StartBlock       uint64
	EndBlock         uint64
	EncryptedBlobLen uint16
	EncryptedBlob    []byte
	AuthTokenLen     uint16
	AuthToken        []byte
	QosLen           uint16
	QosData          string
}

type AppointmentAccepted struct {
	Locator []byte
	Qos     string
	QosLen  uint16
}

type AppointmentRejected struct {
	Locator   []byte
	Rcode     uint16
	Reason    string
	ReasonLen uint16
}

type Mode int

const (
	Altruistic Mode = iota
	Nonaltruistic
)

func (m Mode) String() string {
	return [...]string{"Altruistic", "Nonaltruistic"}[m]
}

type Wt_init struct {
	AcceptedCiphers []string
	Modes           []Mode
	Qos             []string
}
