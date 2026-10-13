
from qa_forge.data import user_data



def test_user_data_default():
    data = user_data(index=1)

    assert data == {
        "name": "User 1",
        "email": "user1@example.com",
    }


def test_user_data_custom_values():
    data = user_data(
        index=5,
        name="Alice",
        email="alice@example.com",
    )

    assert data == {
        "name": "Alice",
        "email": "alice@example.com",
    }


def test_user_data_is_deterministic():
    first = user_data(index=10)
    second = user_data(index=10)

    assert first == second


    assert first == second
