package gsservice

type Msg struct {
	Tp      byte
	BinData []byte
}

func NewMsg(tp byte, data []byte) *Msg {
	return &Msg{
		Tp:      tp,
		BinData: data,
	}
}

func (m *Msg) GetMsgTp() byte {
	return m.Tp
}

func (m *Msg) GetMsgData() []byte {
	return m.BinData
}

func (m *Msg) GetMsgLen() int {
	return len(m.BinData)
}
