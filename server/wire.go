package main 

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
    Modes []Mode
    Qos []string
}
