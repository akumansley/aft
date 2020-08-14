def main(input):
    input = str(input)
    u, ok = urlparse(input)
    if not ok:
        fail("Invalid url: ", input)
    return input