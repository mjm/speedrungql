package speedrungql

func FindLink(links []Link, rel string) string {
	for _, link := range links {
		if link.Rel == rel {
			return link.URI
		}
	}

	return ""
}
