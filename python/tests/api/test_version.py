import requests


def test_version_endpoint(base_url):
    response = requests.get(
        f"{base_url}/version",
        timeout=5,
    )

    assert response.status_code == 200
    assert response.headers["Content-Type"] == "application/json"
    assert response.json() == {"version": "0.1.0"}
