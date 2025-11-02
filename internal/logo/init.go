package logo

import "github.com/armon/go-radix"

var r *radix.Tree

func init() {
	r = radix.New()
	r.Insert("Fedora", DistroWithColor{
		FileName: "Fedora",
		Colors:   []int{12, 7},
	})

	r.Insert("Adélie", DistroWithColor{
		FileName: "adelie",
		Colors:   []int{4, 7, 6},
	})

	r.Insert("Ubuntu", DistroWithColor{
		FileName: "Ubuntu",
		Colors:   []int{1},
	})
	r.Insert("Alpine", DistroWithColor{
		FileName: "Alpine",
		Colors:   []int{4, 5, 7, 6},
	})
}

type DistroWithColor struct {
	FileName string
	Colors   []int
}
