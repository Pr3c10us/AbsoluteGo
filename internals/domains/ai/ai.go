package ai

type Response struct {
	Response string  `json:"response"`
	Dollars  float64 `json:"dollars"`
}

type File struct {
	Path     string `json:"path"`
	MIMEType string `json:"mimeType"`
}

type UploadedFile struct {
	URI      string `json:"uri"`
	MIMEType string `json:"mimeType"`
}

type Voice string

const (
	Zephyr        Voice = "Zephyr"
	Puck          Voice = "Puck"
	Charon        Voice = "Charon"
	Kore          Voice = "Kore"
	Fenrir        Voice = "Fenrir"
	Leda          Voice = "Leda"
	Orus          Voice = "Orus"
	Aoede         Voice = "Aoede"
	Callirrhoe    Voice = "Callirrhoe"
	Autonoe       Voice = "Autonoe"
	Enceladus     Voice = "Enceladus"
	Iapetus       Voice = "Iapetus"
	Umbriel       Voice = "Umbriel"
	Algieba       Voice = "Algieba"
	Despina       Voice = "Despina"
	Erinome       Voice = "Erinome"
	Algenib       Voice = "Algenib"
	Rasalgethi    Voice = "Rasalgethi"
	Laomedeia     Voice = "Laomedeia"
	Achernar      Voice = "Achernar"
	Alnilam       Voice = "Alnilam"
	Schedar       Voice = "Schedar"
	Gacrux        Voice = "Gacrux"
	Pulcherrima   Voice = "Pulcherrima"
	Achird        Voice = "Achird"
	Zubenelgenubi Voice = "Zubenelgenubi"
	Vindemiatrix  Voice = "Vindemiatrix"
	Sadachbia     Voice = "Sadachbia"
	Sadaltager    Voice = "Sadaltager"
	Sulafat       Voice = "Sulafat"
)
