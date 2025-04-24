# Blockchain Wallet Analyzer API

FastAPI сервис для анализа и классификации кошельков в блокчейне.

## Установка

1. Установите зависимости:
```bash
pip install -r requirements.txt
```

2. Убедитесь, что у вас есть обученные модели:
- `train/blockchain_classifier.joblib`
- `train/scaler.joblib`
- `wallet_index.faiss`

## Запуск

Запустите сервис с помощью uvicorn:
```bash
uvicorn app:app --host 0.0.0.0 --port 8000
```

## API Endpoints

### POST /analyze
Анализирует кошелек и возвращает его классификацию.

**Request Body:**
```json
{
    "address": "0x...",
    "transactions": [
        {
            "timestamp": 1234567890,
            "to": "0x...",
            "value": 0.1,
            "method": "transfer"
        },
        ...
    ]
}
```

**Response:**
```json
{
    "predicted_class": "drop_hunter",
    "confidence": 0.85,
    "similar_wallets": [
        {
            "address": "0x...",
            "label": "drop_hunter",
            "transaction_count": 100,
            "description": "..."
        },
        ...
    ]
}
```

### GET /health
Проверка работоспособности сервиса.

## Особенности

- Классифицирует кошельки на 4 типа: дропхантеры, NFT-коллекторы, обычные пользователи, технические аккаунты
- Если уверенность классификации ниже 0.7, возвращает список похожих кошельков
- Использует предобученные модели для быстрого инференса 