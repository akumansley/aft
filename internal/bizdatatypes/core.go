package bizdatatypes

import (
	"awans.org/aft/internal/db"
)

var EmailAddressValidator = db.Code{
	ID:                db.MakeID("ed046b08-ade2-4570-ade4-dd1e31078219"),
	Name:              "emailAddress",
	Runtime:           db.Starlark,
	FunctionSignature: db.FromJSON,
	Code: `# Compile Regular Expression for email addresses
email = re.Compile(r"^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$")
def main(input):
    # Check if input matches the regular expression
    if len(input) > 254 or len(input) < 4 or not email.Match(input):
        # If not, raise an error
        error("Invalid email address: %s", input)
    return input

result(main(args))`,
}

var URLValidator = db.Code{
	ID:                db.MakeID("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	Name:              "url",
	Runtime:           db.Starlark,
	FunctionSignature: db.FromJSON,
	Code: `def main(input):
	# Use a built-in to parse an URL
    u, ok = urlparse(input)
    if not ok:
        # If input is bad, raise an error
        error("Invalid url %s", input)
    return input

result(main(args))`,
}

var EmailAddress = db.DatatypeStorage{
	ID:        db.MakeID("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:      "emailAddress",
	Validator: EmailAddressValidator,
	StoredAs:  db.StringStorage,
}

var URL = db.DatatypeStorage{
	ID:        db.MakeID("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:      "url",
	Validator: URLValidator,
	StoredAs:  db.StringStorage,
}
