package setup

type DTOFieldSlice []*DTOField

func (d DTOFieldSlice) ContainsFieldByName(name string) bool {
	for _, f := range d {
		if f.Name == name {
			return true
		}
	}
	return false
}
