import pytest

from qa_forge.assertions import assert_api_error, assert_json_response
from qa_forge.data import INVALID_EMAILS, INVALID_NAMES


def test_create_user(api_client, user_factory):
    created_user = user_factory(
        name="Kamarik",
        email="kamarik@example.com",
    )

    assert "id" in created_user
    assert isinstance(created_user["id"], int)
    assert created_user["id"] > 0

    assert created_user["name"] == "Kamarik"
    assert created_user["email"] == "kamarik@example.com"

def test_list_users(api_client, user_factory):
    created_user = user_factory(
        name="Kamarik",
        email="kamarik@example.com",
    )

    response = api_client.list_users()

    assert response.status_code == 200
    assert response.headers["Content-Type"] == "application/json"

    users = response.json()

    assert isinstance(users, list)

    user_ids = [user["id"] for user in users]

    assert created_user["id"] in user_ids

    found_user = next(
        user
        for user in users
        if user["id"] == created_user["id"]
    )

    assert found_user["name"] == "Kamarik"
    assert found_user["email"] == "kamarik@example.com"


@pytest.mark.parametrize("name", INVALID_NAMES)
def test_create_user_with_invalid_name(api_client, name):
    response = api_client.create_user(
        name=name,
        email="invalid-name@example.com",
    )

    assert response.status_code == 400


@pytest.mark.parametrize("email", INVALID_EMAILS)
def test_create_user_with_invalid_email(api_client, email):
    response = api_client.create_user(
        name="Kamarik",
        email=email,
    )

    assert response.status_code == 400


def test_get_user(api_client, user_factory):
    created_user = user_factory(
        name="Alice",
        email="alice@example.com",
    )

    user_id = created_user["id"]

    response = api_client.get_user(user_id)

    user = assert_json_response(response, 200)

    assert user["id"] == user_id
    assert user["name"] == "Alice"
    assert user["email"] == "alice@example.com"


def test_get_nonexistent_user(api_client):
    response = api_client.get_user(999999)

    assert_api_error(
        response,
        404,
        "user not found",
    )


@pytest.mark.parametrize(
    "user_id",
    ["abc", 0, -1],
)
def test_get_user_with_invalid_id(api_client, user_id):
    response = api_client.get_user(user_id)

    assert_api_error(
        response,
        400,
        "invalid user id",
    )


def test_update_user(api_client, created_user):
    user_id = created_user["id"]

    response = api_client.update_user(
        user_id=user_id,
        name="Updated Kamarik",
        email="updated@example.com",
    )

    assert response.status_code ==200
    assert response.headers["Content-Type"] == "application/json"

    user = response.json()

    assert user["id"] == user_id
    assert user["name"] == "Updated Kamarik"
    assert user["email"] == "updated@example.com"

@pytest.mark.parametrize(
    "user_id",
    ["abc", 0, -1],
)
def test_update_user_with_invalid_id(api_client, user_id):
    response = api_client.update_user(
        user_id=user_id,
        name="Updated Kamarik",
        email="updated@example.com",
    )

    assert_api_error(
        response,
        400,
        "invalid user id",
    )


def test_update_nonexistent_user(api_client):
    response = api_client.update_user(
        user_id=999999,
        name="Update Kamarik",
        email="update@example.com",
    )

    assert response.status_code == 404
    assert response.headers["Content-Type"] == "application/json"

    body = response.json()

    assert body["error"] == "user not found"


@pytest.mark.parametrize("name", INVALID_NAMES)
def test_update_user_with_ivalid_name(api_client, created_user, name):
    response = api_client.update_user(
        user_id=created_user["id"],
        name=name,
        email="updated@example.com",
    )

    assert response.status_code == 400


@pytest.mark.parametrize("email", INVALID_EMAILS)
def test_update_user_with_invalid_email(api_client, created_user, email):
    response = api_client.update_user(
        user_id=created_user["id"],
        name="Updated Kamarik",
        email=email,
    )

    assert response.status_code == 400


def test_delete_user(api_client, created_user):
    user_id = created_user["id"]

    response = api_client.delete_user(user_id)

    assert response.status_code == 204

    get_response = api_client.get_user(user_id)

    assert get_response.status_code == 404

    body = get_response.json()

    assert body["error"] == "user not found"


@pytest.mark.parametrize(
    "user_id",
    [
        "abc",
        0,
        -1,
    ],
)
def test_delete_user_with_invalid_id(api_client, user_id):
    response = api_client.delete_user(user_id)

    assert_api_error(
        response,
        400,
        "invalid user id",
    )


def test_delete_nonexistent_user(api_client):
    response = api_client.delete_user(999999)

    assert response.status_code == 404
    assert response.headers["Content-Type"] == "application/json"

    body = response.json()

    assert body["error"] == "user not found"


def test_create_multiple_users(user_factory):
    users = []

    for _ in range(3):
        created_user = user_factory()
        users.append(created_user)

    assert len(users) == 3

    ids = [user["id"] for user in users]

    assert len(ids) == len(set(ids))

    names = [user["name"] for user in users]
    emails = [user["email"] for user in users]

    assert names == [
        "User 1",
        "User 2",
        "User 3",
    ]

    assert emails == [
        "user1@example.com",
        "user2@example.com",
        "user3@example.com",
    ]
