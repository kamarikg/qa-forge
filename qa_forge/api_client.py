import requests

class APIClient:
    def __init__(self, base_url):
        self.base_url = base_url.rstrip("/")

    def create_user(self, name=None, email=None):
        payload = {}

        if name is not None:
            payload["name"] = name

        if email is not None:
            payload["email"] = email

        return requests.post(
            f"{self.base_url}/users",
            json=payload,
            timeout=5,
        )

    def list_users(self):
        return requests.get(
            f"{self.base_url}/users",
            timeout=5,
        )

    def get_user(self, user_id):
        return requests.get(
            f"{self.base_url}/users/{user_id}",
            timeout=5,
        )

    def update_user(self, user_id, name=None, email=None):
        payload = {}

        if name is not None:
            payload["name"] = name

        if email is not None:
            payload["email"] = email

        return requests.put(
            f"{self.base_url}/users/{user_id}",
            json=payload,
            timeout=5,
        )

    def delete_user(self, user_id):
        return requests.delete(
            f"{self.base_url}/users/{user_id}",
            timeout=5,
        )
