# Compile Regular Expression for email addresses
email = re.compile(r"^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$")

def main(input):
    input = str(input)
    # Check if input matches the regular expression
    if len(input) > 254 or len(input) < 4 or not email.match(input):
        fail("Invalid email address: ", input)
    return input