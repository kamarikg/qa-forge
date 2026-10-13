def assert_json_response(response, status_code):
    assert response.status_code == status_code
    assert response.headers["Content-Type"] == "application/json"

    return response.json()


def assert_api_error(response, status_code, message):
    body = assert_json_response(response, status_code)

    assert body["error"] == message
