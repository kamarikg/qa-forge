def user_data(
    index=None,
    name=None,
    email=None,
):
    if name is None:
        name = f"User {index}"

    if email is None:
        email = f"user{index}@example.com"

    return {
        "name": name,
        "email": email,
    }

INVALID_NAMES = [
    None,
    "",
    "   ",
]

INVALID_EMAILS = [
    None,
    "",
    "   ",
]
