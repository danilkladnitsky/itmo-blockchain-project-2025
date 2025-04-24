import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split, cross_val_score, StratifiedKFold
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier
from sklearn.svm import SVC
from sklearn.metrics import classification_report, confusion_matrix, accuracy_score, log_loss
import wandb
import os
from dotenv import load_dotenv
from datetime import datetime
import matplotlib.pyplot as plt
from sklearn.calibration import CalibratedClassifierCV
from sklearn.base import clone

load_dotenv()

# Инициализация wandb
wandb.init(project="blockchain-user-classification", entity=os.getenv('WANDB_ENTITY'))

class LossTracker:
    def __init__(self, model_name):
        self.model_name = model_name
        self.train_losses = []
        self.val_losses = []
        self.iterations = []
    
    def update(self, iteration, train_loss, val_loss):
        self.iterations.append(iteration)
        self.train_losses.append(train_loss)
        self.val_losses.append(val_loss)
        wandb.log({
            f"{self.model_name}_train_loss": train_loss,
            f"{self.model_name}_val_loss": val_loss,
            "iteration": iteration
        })
    
    def plot_losses(self):
        plt.figure(figsize=(10, 6))
        plt.plot(self.iterations, self.train_losses, label='Training Loss')
        plt.plot(self.iterations, self.val_losses, label='Validation Loss')
        plt.xlabel('Iteration')
        plt.ylabel('Loss')
        plt.title(f'{self.model_name} Loss Curves')
        plt.legend()
        wandb.log({f"{self.model_name}_loss_curves": wandb.Image(plt)})
        plt.close()

def load_data():
    """Загрузка и подготовка данных из CSV файла"""
    df = pd.read_csv('../data/data.csv')

    # Преобразуем timestamp в datetime
    df['timestamp'] = pd.to_datetime(df['timestamp'], unit='s', errors='coerce')

    # Удалим строки с битым временем
    df = df.dropna(subset=['timestamp'])

    # Группируем по адресу и метке
    grouped = df.groupby(['address', 'label']).agg({
        'timestamp': ['count', 'min', 'max'],
        'value': ['mean', 'sum', 'std', 'min', 'max'],
        'method': lambda x: x.nunique()
    }).reset_index()

    # Переименование колонок
    grouped.columns = ['address', 'label', 'transaction_count', 'first_transaction', 
                       'last_transaction', 'mean_value', 'total_value', 'value_std', 
                       'min_value', 'max_value', 'unique_methods']

    # Временные признаки
    grouped['transaction_duration'] = (grouped['last_transaction'] - grouped['first_transaction']).dt.total_seconds()
    grouped['avg_transaction_interval'] = grouped['transaction_duration'] / grouped['transaction_count']

    # Безопасное деление для value_std_norm
    grouped['value_range'] = grouped['max_value'] - grouped['min_value']
    grouped['value_std_norm'] = np.where(
        grouped['mean_value'].abs() > 1e-8,
        grouped['value_std'] / grouped['mean_value'],
        0.0
    )

    # Интенсивность активности
    grouped['transaction_intensity'] = np.where(
        grouped['transaction_duration'] > 0,
        grouped['transaction_count'] / grouped['transaction_duration'],
        0.0
    )

    # Добавим шум
    np.random.seed(42)
    noise_scale = 0.1
    numeric_columns = [
        'transaction_count', 'mean_value', 'total_value', 'value_std',
        'min_value', 'max_value', 'value_range', 'value_std_norm',
        'transaction_intensity'
    ]

    for col in numeric_columns:
        noise = 1 + noise_scale * np.random.randn(len(grouped))
        grouped[col] = grouped[col] * noise

    # Очистка от бесконечностей и NaN
    grouped.replace([np.inf, -np.inf], np.nan, inplace=True)
    grouped.dropna(inplace=True)

    return grouped

def extract_features(row):
    """Извлечение признаков из строки данных"""
    features = [
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
        # Добавляем признаки из адреса
        len(row['address']),
        int(row['address'][2:4], 16),  # первые два символа после 0x
        int(row['address'][-4:], 16),  # последние два символа
    ]
    return features

def train_and_evaluate_model(model, X_train, y_train, X_val, y_val, X_test, y_test, model_name):
    """Обучение и оценка модели с отслеживанием лосса"""
    tracker = LossTracker(model_name)
    
    # Для RandomForest и GradientBoosting можем отслеживать лосс во время обучения
    if isinstance(model, (RandomForestClassifier, GradientBoostingClassifier)):
        # Создаем копию модели для отслеживания прогресса
        temp_model = clone(model)
        
        # Для RandomForest
        if isinstance(model, RandomForestClassifier):
            n_estimators = model.n_estimators
            for i in range(1, n_estimators + 1):
                temp_model.n_estimators = i
                temp_model.fit(X_train, y_train)
                
                # Предсказания и лосс
                train_pred_proba = temp_model.predict_proba(X_train)
                val_pred_proba = temp_model.predict_proba(X_val)
                
                train_loss = log_loss(y_train, train_pred_proba)
                val_loss = log_loss(y_val, val_pred_proba)
                
                tracker.update(i, train_loss, val_loss)
            
            # Обучаем основную модель
            model.fit(X_train, y_train)
        
        # Для GradientBoosting
        elif isinstance(model, GradientBoostingClassifier):
            n_estimators = model.n_estimators
            for i in range(1, n_estimators + 1):
                temp_model.n_estimators = i
                temp_model.fit(X_train, y_train)
                
                train_pred_proba = temp_model.predict_proba(X_train)
                val_pred_proba = temp_model.predict_proba(X_val)
                
                train_loss = log_loss(y_train, train_pred_proba)
                val_loss = log_loss(y_val, val_pred_proba)
                
                tracker.update(i, train_loss, val_loss)
            
            # Обучаем основную модель
            model.fit(X_train, y_train)
    
    # Для SVM просто обучаем и оцениваем
    else:
        model.fit(X_train, y_train)
        train_pred_proba = model.predict_proba(X_train)
        val_pred_proba = model.predict_proba(X_val)
        
        train_loss = log_loss(y_train, train_pred_proba)
        val_loss = log_loss(y_val, val_pred_proba)
        
        tracker.update(1, train_loss, val_loss)
    
    # Окончательная оценка
    test_pred = model.predict(X_test)
    test_pred_proba = model.predict_proba(X_test)
    test_accuracy = accuracy_score(y_test, test_pred)
    test_loss = log_loss(y_test, test_pred_proba)
    
    # Логирование финальных метрик
    wandb.log({
        f"{model_name}_final_test_accuracy": test_accuracy,
        f"{model_name}_final_test_loss": test_loss
    })
    
    # Построение графиков лосса
    tracker.plot_losses()
    
    return model, val_loss, test_loss

def train_model():
    """Обучение и сравнение моделей"""
    # Загрузка данных
    df = load_data()
    
    # Проверяем уникальные метки
    unique_labels = df['label'].unique()
    print("Уникальные метки в данных:", unique_labels)
    
    # Извлечение признаков
    X = np.array([extract_features(row) for _, row in df.iterrows()])
    y = df['label']
    
    # Разделение на train/val/test
    X_train, X_temp, y_train, y_temp = train_test_split(X, y, test_size=0.3, random_state=42, stratify=y)
    X_val, X_test, y_val, y_test = train_test_split(X_temp, y_temp, test_size=0.5, random_state=42, stratify=y_temp)
    
    # Масштабирование признаков
    scaler = StandardScaler()
    X_train = scaler.fit_transform(X_train)
    X_val = scaler.transform(X_val)
    X_test = scaler.transform(X_test)
    
    # Определяем модели для сравнения
    models = {
        "RandomForest": RandomForestClassifier(
            n_estimators=200,
            max_depth=10,
            min_samples_split=5,
            min_samples_leaf=2,
            random_state=42,
            class_weight='balanced'
        ),
        "GradientBoosting": GradientBoostingClassifier(
            n_estimators=200,
            learning_rate=0.1,
            max_depth=5,
            random_state=42
        ),
        "SVM": CalibratedClassifierCV(
            SVC(probability=True, random_state=42),
            cv=5
        )
    }
    
    # Обучаем и оцениваем каждую модель
    results = {}
    for name, model in models.items():
        print(f"\nTraining {name}...")
        trained_model, val_loss, test_loss = train_and_evaluate_model(
            model, X_train, y_train, X_val, y_val, X_test, y_test, name
        )
        results[name] = {
            'model': trained_model,
            'val_loss': val_loss,
            'test_loss': test_loss
        }
    
    # Создаем график сравнения моделей
    plt.figure(figsize=(10, 6))
    models_names = list(results.keys())
    val_losses = [results[name]['val_loss'] for name in models_names]
    test_losses = [results[name]['test_loss'] for name in models_names]
    
    x = np.arange(len(models_names))
    width = 0.35
    
    plt.bar(x - width/2, val_losses, width, label='Validation Loss')
    plt.bar(x + width/2, test_losses, width, label='Test Loss')
    
    plt.xlabel('Models')
    plt.ylabel('Loss')
    plt.title('Comparison of Models')
    plt.xticks(x, models_names)
    plt.legend()
    
    # Сохраняем график в wandb
    wandb.log({"model_comparison": wandb.Image(plt)})
    plt.close()
    
    # Сохраняем лучшую модель
    best_model_name = min(results, key=lambda x: results[x]['val_loss'])
    best_model = results[best_model_name]['model']
    
    import joblib
    joblib.dump(best_model, 'best_classifier.joblib')
    joblib.dump(scaler, 'scaler.joblib')
    
    return best_model, scaler

if __name__ == "__main__":
    train_model()
    wandb.finish() 