package proxy

type Hooks struct {
	OnUpstream2xx func(chosenBaseURL string)
}
