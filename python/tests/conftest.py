import os
import pytest

from qa_forge.api_client import APIClient
from qa_forge.data import user_data


@pytest.fixture
def base_url():
    return os.getenv(
	"QA_FORGE_BASE_URL",
	"http://127.0.0.1:8080",
    )


@pytest.fixture
def api_client(base_url):
    return APIClient(base_url)


@pytest.fixture
def cleanup_user(api_client):
    def cleanup(user_id):
        response = api_client.get_user(user_id)

        if response.status_code == 200:
            delete_response = api_client.delete_user(user_id)

            assert delete_response.status_code == 204

        elif response.status_code == 404:
            return

        else:
            pytest.fail(
                f"cleanup check failed for user {user_id}: "
                f"expected 200 or 404 , got {response.status_code}"
            )

    return cleanup


@pytest.fixture
def created_user(api_client, cleanup_user):
    response = api_client.create_user(
        name="Kamarik",
        email="kamarik@example.com",
    )

    assert response.status_code == 201

    user = response.json()

    yield user

    cleanup_user(user["id"])


@pytest.fixture
def user_factory(api_client, cleanup_user):
    created_users = []
    user_index = 0

    def create_user(name=None, email=None):
        nonlocal user_index

        user_index += 1

        data = user_data(
            index=user_index,
            name=name,
            email=email,
        )

        response = api_client.create_user(
            name=data["name"],
            email=data["email"],
        )

        assert response.status_code == 201

        user = response.json()

        created_users.append(user)

        return user

    yield create_user

    for user in created_users:
        cleanup_user(user["id"])
