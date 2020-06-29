package bizdatatypes

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
)

var EmailAddressValidator = starlark.MakeStarlarkFunction(
	db.MakeID("ed046b08-ade2-4570-ade4-dd1e31078219"),
	"emailAddressValidator",
	db.FromJSON,
	`# Compile Regular Expression for email addresses
email = re.Compile(r"^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$")
def validator(input):
    # Check if input matches the regular expression
    if len(input) > 254 or len(input) < 4 or not email.Match(input):
        # If not, raise an error
        error("Invalid email address: %s", input)
    return input
    `)

var URLValidator = starlark.MakeStarlarkFunction(
	db.MakeID("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	"urlValidator",
	db.FromJSON,
	`def validator(input):
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
