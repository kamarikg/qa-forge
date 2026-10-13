import requests


def test_health_endpoint(base_url):
    response = requests.get(
        f"{base_url}/health",
        timeout=5,
    )

    assert response.status_code == 200
    assert response.headers["Content-Type"] == "application/json"
    assert response.json() == {"status": "ok"}


def test_health_rejects_post(base_url):
    responses = requests.post(
        f"{base_url}/health",
        timeout=5,
    )

    assert responses.status_code == 405
