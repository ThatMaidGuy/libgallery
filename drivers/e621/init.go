package e621

import "github.com/ThatMaidGuy/libgallery"

func init() {
	libgallery.Register("e621", New())
}
