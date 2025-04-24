import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import classification_report, confusion_matrix, accuracy_score
import wandb
import os
from dotenv import load_dotenv
from datetime import datetime

load_dotenv()

# Инициализация wandb
wandb.init(project="blockchain-user-classification", entity=os.getenv('WANDB_ENTITY'))

def load_data():
    """Загрузка и подготовка данных из CSV файла"""
    df = pd.read_csv('../data/synthetic_wallet_data_with_tech.csv')
    
    # Группируем данные по адресам для создания признаков
    grouped = df.groupby(['address', 'label']).agg({
        'timestamp': ['count', 'min', 'max'],
        'value': ['mean', 'sum', 'std'],
        'method': lambda x: x.nunique()
    }).reset_index()
    
    # Переименовываем колонки
    grouped.columns = ['address', 'label', 'transaction_count', 'first_transaction', 
                      'last_transaction', 'mean_value', 'total_value', 'value_std', 
                      'unique_methods']
    
    # Добавляем временные признаки
    grouped['transaction_duration'] = (grouped['last_transaction'] - grouped['first_transaction']).dt.total_seconds()
    grouped['avg_transaction_interval'] = grouped['transaction_duration'] / grouped['transaction_count']
    
    return grouped

def extract_features(row):
    """Извлечение признаков из строки данных"""
    features = [
        row['transaction_count'],
        row['mean_value'],
        row['total_value'],
        row['value_std'],
        row['unique_methods'],
        row['transaction_duration'],
        row['avg_transaction_interval'],
        # Добавляем признаки из адреса
        len(row['address']),
        int(row['address'][2:4], 16),  # первые два символа после 0x
        int(row['address'][-4:], 16),  # последние два символа
    ]
    return features

def train_model():
    """Обучение модели"""
    # Загрузка данных
    df = load_data()
    
    # Извлечение признаков
    X = np.array([extract_features(row) for _, row in df.iterrows()])
    y = df['label']
    
    # Разделение на train/val/test
    X_train, X_temp, y_train, y_temp = train_test_split(X, y, test_size=0.3, random_state=42)
    X_val, X_test, y_val, y_test = train_test_split(X_temp, y_temp, test_size=0.5, random_state=42)
    
    # Масштабирование признаков
    scaler = StandardScaler()
    X_train = scaler.fit_transform(X_train)
    X_val = scaler.transform(X_val)
    X_test = scaler.transform(X_test)
    
    # Создание и обучение модели
    model = RandomForestClassifier(
        n_estimators=200,
        max_depth=10,
        min_samples_split=5,
        min_samples_leaf=2,
        random_state=42
    )
    model.fit(X_train, y_train)
    
    # Оценка на валидационном наборе
    val_pred = model.predict(X_val)
    val_accuracy = accuracy_score(y_val, val_pred)
    
    # Оценка на тестовом наборе
    test_pred = model.predict(X_test)
    test_accuracy = accuracy_score(y_test, test_pred)
    
    # Логирование метрик в wandb
    wandb.log({
        "val_accuracy": val_accuracy,
        "test_accuracy": test_accuracy,
        "confusion_matrix": wandb.plot.confusion_matrix(
            y_true=y_test,
            preds=test_pred,
            class_names=model.classes_
        )
    })
    
    # Вывод отчета о классификации
    print("\nClassification Report:")
    print(classification_report(y_test, test_pred))
    
    # Сохранение модели
    import joblib
    joblib.dump(model, 'blockchain_classifier.joblib')
    joblib.dump(scaler, 'scaler.joblib')
    
    return model, scaler

if __name__ == "__main__":
    train_model()
    wandb.finish() 