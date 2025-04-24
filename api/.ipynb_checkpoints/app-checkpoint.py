from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Optional
import joblib
import numpy as np
from wallet_similarity import WalletSimilarity
import pandas as pd
from datetime import datetime

app = FastAPI(title="Blockchain Wallet Analyzer")

# Загрузка моделей
classifier = joblib.load('train/blockchain_classifier.joblib')
scaler = joblib.load('train/scaler.joblib')
similarity_engine = WalletSimilarity()
similarity_engine.load_data()
similarity_engine.load_index()

class WalletRequest(BaseModel):
    address: str
    transactions: List[Dict]

class ClassificationResult(BaseModel):
    predicted_class: str
    confidence: float
    similar_wallets: Optional[List[Dict]] = None

def extract_features(transactions: List[Dict]) -> np.ndarray:
    """Извлекает признаки из транзакций для классификации"""
    df = pd.DataFrame(transactions)
    
    # Группируем данные по адресу
    grouped = df.groupby('address').agg({
        'timestamp': ['count', 'min', 'max'],
        'value': ['mean', 'sum', 'std'],
        'method': lambda x: x.nunique()
    }).reset_index()
    
    # Переименовываем колонки
    grouped.columns = ['address', 'transaction_count', 'first_transaction', 
                      'last_transaction', 'mean_value', 'total_value', 'value_std', 
                      'unique_methods']
    
    # Добавляем временные признаки
    grouped['transaction_duration'] = (grouped['last_transaction'] - grouped['first_transaction']).dt.total_seconds()
    grouped['transaction_duration'] = grouped['transaction_duration'].fillna(0)
    grouped['avg_transaction_interval'] = grouped['transaction_duration'] / grouped['transaction_count']
    grouped['avg_transaction_interval'] = grouped['avg_transaction_interval'].fillna(0)
    
    # Извлекаем признаки из адреса
    features = []
    for _, row in grouped.iterrows():
        features.append([
            row['transaction_count'],
            row['mean_value'],
            row['total_value'],
            row['value_std'],
            row['unique_methods'],
            row['transaction_duration'],
            row['avg_transaction_interval'],
            len(row['address']),
            int(row['address'][2:4], 16),
            int(row['address'][-4:], 16),
        ])
    
    return np.array(features)

@app.post("/analyze", response_model=ClassificationResult)
async def analyze_wallet(wallet_request: WalletRequest):
    try:
        # Извлекаем признаки
        features = extract_features(wallet_request.transactions)
        
        # Масштабируем признаки
        scaled_features = scaler.transform(features)
        
        # Получаем предсказания и вероятности
        prediction = classifier.predict(scaled_features)[0]
        probabilities = classifier.predict_proba(scaled_features)[0]
        confidence = max(probabilities)
        
        result = ClassificationResult(
            predicted_class=prediction,
            confidence=float(confidence)
        )
        
        # Если уверенность низкая, ищем похожие кошельки
        if confidence < 0.7:
            similar_wallets = similarity_engine.find_similar_wallets(wallet_request.address)
            result.similar_wallets = similar_wallets
            
        return result
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "healthy"} 