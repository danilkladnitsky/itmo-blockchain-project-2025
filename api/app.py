from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Optional
import joblib
import numpy as np
from wallet_similarity import WalletSimilarity
import pandas as pd
from datetime import datetime
import logging
import traceback

# Настройка логирования
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = FastAPI(title="Blockchain Wallet Analyzer")

# Загрузка моделей
try:
    logger.info("Loading models...")
    classifier = joblib.load('train/blockchain_classifier.joblib')
    scaler = joblib.load('train/scaler.joblib')
    similarity_engine = WalletSimilarity()
    similarity_engine.load_data()
    logger.info("Building similarity index...")
    similarity_engine.build_index()  # Строим индекс
    logger.info("Models loaded successfully")
except Exception as e:
    logger.error(f"Error loading models: {str(e)}")
    logger.error(traceback.format_exc())
    raise

class WalletRequest(BaseModel):
    address: str
    transactions: List[Dict]

class ClassificationResult(BaseModel):
    predicted_class: str
    confidence: float
    similar_wallets: Optional[List[Dict]] = None

def extract_features(transactions: List[Dict], wallet_address: str) -> np.ndarray:
    """Извлекает признаки из транзакций для классификации"""
    try:
        logger.debug(f"Extracting features from {len(transactions)} transactions")
        # Добавляем адрес в каждую транзакцию
        for tx in transactions:
            tx['address'] = wallet_address
            
        df = pd.DataFrame(transactions)
        logger.debug(f"DataFrame created with columns: {df.columns.tolist()}")
        
        # Преобразуем timestamp в datetime
        df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s')
        
        # Группируем данные по адресу
        grouped = df.groupby('address').agg({
            'timestamp': ['count', 'min', 'max'],
            'value': ['mean', 'sum', 'std', 'min', 'max'],
            'method': lambda x: x.nunique()
        }).reset_index()
        
        logger.debug(f"Grouped data shape: {grouped.shape}")
        
        # Переименовываем колонки
        grouped.columns = ['address', 'transaction_count', 'first_transaction', 
                          'last_transaction', 'mean_value', 'total_value', 'value_std', 
                          'min_value', 'max_value', 'unique_methods']
        
        # Добавляем временные признаки
        grouped['transaction_duration'] = (grouped['last_transaction'] - grouped['first_transaction']).dt.total_seconds()
        grouped['transaction_duration'] = grouped['transaction_duration'].fillna(0)
        grouped['avg_transaction_interval'] = grouped['transaction_duration'] / grouped['transaction_count']
        grouped['avg_transaction_interval'] = grouped['avg_transaction_interval'].fillna(0)
        
        # Добавляем дополнительные признаки как в train_classifier.py
        grouped['value_range'] = grouped['max_value'] - grouped['min_value']
        grouped['value_std_norm'] = np.where(
            grouped['mean_value'].abs() > 1e-8,
            grouped['value_std'] / grouped['mean_value'],
            0.0
        )
        
        grouped['transaction_intensity'] = np.where(
            grouped['transaction_duration'] > 0,
            grouped['transaction_count'] / grouped['transaction_duration'],
            0.0
        )
        
        # Извлекаем признаки из адреса
        features = []
        for _, row in grouped.iterrows():
            try:
                features.append([
                    row['transaction_count'],
                    row['mean_value'],
                    row['total_value'],
                    row['value_std'],
                    row['min_value'],
                    row['max_value'],
                    row['value_range'],
                    row['value_std_norm'],
                    row['unique_methods'],
                    row['transaction_duration'],
                    row['avg_transaction_interval'],
                    row['transaction_intensity'],
                    len(row['address']),
                    int(row['address'][2:4], 16),
                    int(row['address'][-4:], 16)
                ])
            except Exception as e:
                logger.error(f"Error processing row: {row}")
                logger.error(f"Error details: {str(e)}")
                raise
        
        logger.debug(f"Extracted features shape: {np.array(features).shape}")
        return np.array(features)
        
    except Exception as e:
        logger.error(f"Error in extract_features: {str(e)}")
        logger.error(traceback.format_exc())
        raise

@app.post("/analyze", response_model=ClassificationResult)
async def analyze_wallet(wallet_request: WalletRequest):
    try:
        logger.info(f"Analyzing wallet: {wallet_request.address}")
        logger.debug(f"Request data: {wallet_request.dict()}")
        
        # Извлекаем признаки
        features = extract_features(wallet_request.transactions, wallet_request.address)
        logger.debug(f"Features extracted: {features}")
        
        # Масштабируем признаки
        scaled_features = scaler.transform(features)
        logger.debug(f"Features scaled: {scaled_features}")
        
        # Получаем предсказания и вероятности
        prediction = classifier.predict(scaled_features)[0]
        probabilities = classifier.predict_proba(scaled_features)[0]
        confidence = max(probabilities)
        
        logger.info(f"Prediction: {prediction}, Confidence: {confidence}")
        
        # Если уверенность низкая, ищем похожие кошельки
        if confidence < 0.4:
            logger.info("Low confidence, searching for similar wallets")
            similar_wallets = similarity_engine.find_similar_wallets(wallet_request.address, k=1)  # Берем только самый похожий
            
            if similar_wallets:
                # Берем метку от самого похожего кошелька
                most_similar_label = similar_wallets[0]['label']
                logger.info(f"Using label from most similar wallet: {most_similar_label}")
                prediction = most_similar_label
                confidence = 0.4  # Устанавливаем уверенность на пороговое значение
            else:
                logger.warning("No similar wallets found, keeping original prediction")
        
        result = ClassificationResult(
            predicted_class=prediction,
            confidence=float(confidence)
        )
        
        return result
        
    except Exception as e:
        logger.error(f"Error in analyze_wallet: {str(e)}")
        logger.error(traceback.format_exc())
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "healthy"} 