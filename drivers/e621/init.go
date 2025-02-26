package e621

import "github.com/WheatleyHDD/libgallery"

func init() {
	libgallery.Register("e621", New())
}
