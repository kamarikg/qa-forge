import pytest


@pytest.fixture(scope="module")
def database():
    print("\n>>> SETUP database")
    yield "DB"
    print(">>> TEARDOWN database")

@pytest.fixture(scope="function")
def user(database):
    print(">>> SETUP user")

    data = {
        "name": "Alice",
        "db": database,
    }

    yield data

    print(">>> TEARDOWN user")


def test_one(user):
    print(f">>> TEST ONE BEFOR: {user}")

    user["name"] = "Bob"

    print(f">>> TEST ONE AFTER: {user}")

    assert user["name"] == "Bob"


def test_two(user):
    print(f">>> TEST TWO: {user}")

    assert user["name"] == "Alice"
