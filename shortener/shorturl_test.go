package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkGenerator(t *testing.T) {
	UserId := "test-user-id"

	initialLink_1 := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	shortLink_1 := GenerateShortLink(initialLink_1, UserId)

	initialLink_2 := "https://www.eddywm.com/lets-build-a-url-shortener-in-go-with-redis-part-2-storage-layer/"
	shortLink_2 := GenerateShortLink(initialLink_2, UserId)

	initialLink_3 := "https://spectrum.ieee.org/automaton/robotics/home-robots/hello-robots-stretch-mobile-manipulator"
	shortLink_3 := GenerateShortLink(initialLink_3, UserId)

	assert.Equal(t, shortLink_1, "KuszwPtD")
	assert.Equal(t, shortLink_2, "DTVmGtb7")
	assert.Equal(t, shortLink_3, "ZL5YSmF4")

	t.Log("shortLink_1:", shortLink_1)
	t.Log("shortLink_2:", shortLink_2)
	t.Log("shortLink_3:", shortLink_3)
}
