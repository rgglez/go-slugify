package main

import (
	"fmt"

	"github.com/rgglez/go-slugify"
)

func main() {
	var slug string

	slug, _ = slugify.Slugify("I ♥ Dogs")
	fmt.Println(slug)
	//=> 'i-love-dogs'

	slug, _ = slugify.Slugify("  Déjà Vu!  ")
	fmt.Println(slug)
	//=> 'deja-vu'

	slug, _ = slugify.Slugify("fooBar 123 $#%")
	fmt.Println(slug)
	//=> 'foo-bar-123'

	slug, _ = slugify.Slugify("я люблю единорогов")
	fmt.Println(slug)
	//=> 'ya-lyublyu-edinorogov'
}
