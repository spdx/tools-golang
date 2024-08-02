package v3_0

type Relationships struct {
	relationships []IRelationship
}

func (r Relationships) From(element IElement) Relationships {
	var out []IRelationship
	for _, r := range r.relationships {
		if r.GetFrom() == element {
			out = append(out, r)
		}
	}
	return Relationships{out}
}
