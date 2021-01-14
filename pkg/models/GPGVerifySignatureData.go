package models

type GPGVerifySignatureData struct {
	Base64Data string `example:"SGVsbG8gd29ybGQK"`
	Signature  string `example:"0551F452ABE463A4_SHA512_wsDcBAABCgAQBQJf+LYnCRAFUfRSq+RjpAAA7/oMACHJPMtQs4rr0uxX4AMZ8akb+x2p5ZYL+uRug+zctp82sJEJmL76HG++UyzDmMUCagJ+LBWp2RcCQvfsIhX5MqD7lPkEdtl0uNCIU40apvzn1+0kndl7LnFtzyHMWrHrRqEFGJ0E2APPqv7g1pehVKeusMOkTNUmmsJNgZBYrluZxHnai/Rudoe9jBxihY4ALF0eOyTCHbtWy0z6fll3Bo/iPe777kplDXmTBzCEM8uD3/VZmY6pGn6oXUov/z8Dcrg2x5qT4i5DgdF8OSLbsxVW2OIV8DwCicQCT2tK95fctBqJ22vfmhNlxI3KzI9ShxeV6Eci5p5Zydgoh77pDiWDysrq1dOZ+o7T+ij72K3s63w3loERFVoDxDuKG3jS3+fj+ggqqtpUpm957+9+4QlnJqZk0v9TKT661HnoH4MfZR3muBir8/dgF4mNtuQLSswOxdVs1sHSC3ssTIzzpQqeI2iy3m8Svgl5unAdv2QE81EM/wT5brc2R/abSRz52A===J34T"`
}

// Used only for documentation
type GPGVerifySignatureDataNonQuanto struct {
	Base64Data string `example:"SGVsbG8gd29ybGQK"`
	Signature  string `example:"-----BEGIN PGP SIGNATURE-----\n\nwsDcBAABCgAQBQJf+LriCRAFUfRSq+RjpAAAuL0MAGGrSJfK/tnMkwZ2Rkh3JcvF\nE8WU8jwc8quz+0p9gMDscby0jShJ2G2XXMm3WAYXW88J6h8u2E/lTb6l3oBq/FPb\n15gTM5Ie0p0kHBUlgP5bkV9EF9+VQif40fhVX7OPrS27jWtVNP374ARzSIgKMLa6\nKBZhV1eQecLIlEYXahUP9jyt4cR4A4d9P+YJS/L6d/tQT4g9DBo66hYt5lu4sagG\nDHsW2HK9I7fizCBaE8azLtQd3RRFTWZshln7OGVypwcdbzWbYr5uEhituxAnZKS4\nSWwI0hgj1OkZeOhKwaydtITnaeH+nmlLBzhGKQWjCiLlsDNkkp3/4FKOuYJkYXeZ\nm61GV6G5ZpW/gFVJXXyPz6ElNfWCorZQvxLbY4YWTBLdLyblHnp9kshav6dnexN1\nwQyBDk8jxucmKNE8kCu591dPj/g/H38/zpGZQhj8Firb0rCFumqsAwxFeyTEFjVI\ncyDHa5K+ytmSrITIdQUUsp1M4UQiRH63c1HYOLQurw==\n=BRZt\n-----END PGP SIGNATURE-----"`
}
