from pyngrok import ngrok

# Установи токен, если ещё не установила
ngrok.set_auth_token("2UOoLUe9uBEqvN4C5frNQRVXcF9_3YQ8xd3PczggxmFWkK9ME")

# Проброс локального порта 8000
public_url = ngrok.connect(8000)

print(f"Public URL: {public_url}")