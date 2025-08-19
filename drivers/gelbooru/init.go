package gelbooru

import "github.com/ThatMaidGuy/libgallery"

func init() {
	libgallery.Register("gelbooru", New("safebooru", "safebooru.org"))
}
