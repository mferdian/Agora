package websocket

// Message adalah struktur data yang dikirim dan diterima antar client.
type Message struct {
	Type       string `json:"type"`        // contoh: "comment"
	ProposalID string `json:"proposal_id"` // ID proposal
	Content    string `json:"content"`     // isi komentar
	UserID     string `json:"user_id"`     // ID user (opsional, bisa diambil dari token nanti)
}
