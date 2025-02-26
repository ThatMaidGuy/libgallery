package gelbooru

import "github.com/WheatleyHDD/libgallery"

func init() {
	libgallery.Register("gelbooru", New())
}
