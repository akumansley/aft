# Compile Regular Expression for valid US Phone Numbers
phone = re.compile(r"^\D?(\d{3})\D?\D?(\d{3})\D?(\d{4})$")

def main(input):
    input = str(input)
    if not phone.match(input):
        fail("Bad phone number: ", input)
    clean = input.replace(" ","").replace("-","")
    return clean.replace("(","").replace(")","")