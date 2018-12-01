package auxiliary

// Domain contains the domain name and a map of subdomains
type Domain struct {
	Domain     string
	Subdomains map[string]bool
	Source     string
}
