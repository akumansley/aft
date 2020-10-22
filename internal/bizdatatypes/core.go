package bizdatatypes

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
)

var EmailAddressValidator = starlark.MakeStarlarkFunction(
	db.MakeID("ed046b08-ade2-4570-ade4-dd1e31078219"),
	"emailAddressValidator",
	1,
	`# Compile Regular Expression for email addresses
email = re.compile(r"^([a-zA-Z0-9_\+\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$")

def main(input):
    # Check if input matches the regular expression
    if len(input) > 254 or len(input) < 4 or not email.match(input):
        # If not, raise an error
        fail("Invalid email address: ", input)
    return input`,
)

var URLValidator = starlark.MakeStarlarkFunction(
	db.MakeID("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	"urlValidator",
	1,
	`def main(input):
	# Use a built-in to parse an URL
    u, ok = urlparse(input)
    if not ok:
        # If input is bad, raise an error
        error("Invalid url %s", input)
    return input
`)

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
	1,
	`# Compile Regular Expression for valid US Phone Numbers
phone = re.compile(r"^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$")

def main(input):
    if not phone.match(input):
        fail("Bad phone number: ", input)
    # Otherwise, return it stripped of formatting
    clean = input.replace(" ","").replace("-","")
    return clean.replace("(","").replace(")","")`,
)

var Phone = db.MakeCoreDatatype(
	db.MakeID("d5b7bc19-9eec-4bf9-b362-1a642458060f"),
	"phone",
	db.IntStorage,
	PhoneValidator,
)
