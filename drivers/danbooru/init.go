package danbooru

import "github.com/ThatMaidGuy/libgallery"

func init() {
	libgallery.Register("danbooru", New("Danbooru", "danbooru.donmai.us"))
	libgallery.Register("thebub.club", New("The Bub Club", "thebub.club"))
}
