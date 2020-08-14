package bizdatatypes

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"io/ioutil"
	"os"
	"strings"
)

func email() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	emailAddressB, _ := ioutil.ReadFile(base + "aft/internal/bizdatatypes/emailAddress.star")
	return string(emailAddressB)
}

func url() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	urlB, _ := ioutil.ReadFile(base + "aft/internal/bizdatatypes/url.star")
	return string(urlB)
}

func phone() string {
    wd, _ := os.Getwd()
    base := strings.Split(wd, "aft")[0]
	phoneB, _ := ioutil.ReadFile(base + "aft/internal/bizdatatypes/phone.star")
	return string(phoneB)
}

var EmailAddressValidator = starlark.MakeStarlarkFunction(
	db.MakeID("ed046b08-ade2-4570-ade4-dd1e31078219"),
	"emailAddressValidator",
	db.FromJSON, email())

var URLValidator = starlark.MakeStarlarkFunction(
	db.MakeID("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	"urlValidator",
	db.FromJSON, url())

var EmailAddress = db.MakeCoreDatatype(
	db.MakeID("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	"emailAddress",
	db.StringStorage,
	EmailAddressValidator,
)

var URL = db.MakeCoreDatatype(
	db.MakeID("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	"url",
	db.StringStorage,
	URLValidator,
)

var PhoneValidator = starlark.MakeStarlarkFunction(
	db.MakeID("f720efdc-3694-429f-9d4e-c2150388bd30"),
	"phone",
	db.FromJSON, phone())

var Phone = db.MakeCoreDatatype(
	db.MakeID("d5b7bc19-9eec-4bf9-b362-1a642458060f"),
	"phone",
	db.IntStorage,
	PhoneValidator,
)
