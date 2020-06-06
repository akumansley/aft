package bizdatatypes

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var emailAddressValidator = db.Code{
	ID:      uuid.MustParse("ed046b08-ade2-4570-ade4-dd1e31078219"),
	Name:    "emailAddressValidator",
	Runtime: db.Starlark,
	Code: `# Compile Regular Expression for email addresses
email = re.Compile(r"^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$")
def validator(input):
    # Check if input matches the regular expression
    if len(input) > 254 or len(input) < 4 or not email.Match(input):
        # If not, raise an error
        error("Invalid email address: %s", input)
    return input
`,
}

var URLValidator = db.Code{
	ID:      uuid.MustParse("259d9049-b21e-44a4-abc5-79b0420cda5f"),
	Name:    "urlValidator",
	Runtime: db.Starlark,
	Code: `def validator(input):
	# Use a built-in to parse an URL
    u, ok = urlparse(input)
    if not ok:
        # If input is bad, raise an error
        error("Invalid url %s", input)
    return input
`,
}

var EmailAddress = db.Datatype{
	ID:        uuid.MustParse("6c5e513b-9965-4463-931f-dd29751f5ae1"),
	Name:      "emailAddress",
	Validator: emailAddressValidator,
	StoredAs:  db.StringStorage,
}

var URL = db.Datatype{
	ID:        uuid.MustParse("84c8c2c5-ff1a-4599-9605-b56134417dd7"),
	Name:      "url",
	Validator: URLValidator,
	StoredAs:  db.StringStorage,
}
