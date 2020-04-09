package interfaces

import "github.com/quan-to/chevron/internal/models"

type PKSInterface interface {
	PKSGetKey(fingerPrint string) string
	PKSSearchByName(name string, pageStart, pageEnd int) []models.GPGKey
	PKSSearchByFingerPrint(fingerPrint string, pageStart, pageEnd int) []models.GPGKey
	PKSSearchByEmail(email string, pageStart, pageEnd int) []models.GPGKey
	PKSSearch(value string, pageStart, pageEnd int) []models.GPGKey
	PKSAdd(pubKey string) string
}
