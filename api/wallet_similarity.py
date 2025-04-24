import pandas as pd
import numpy as np
import faiss
from sentence_transformers import SentenceTransformer
from typing import List, Dict
import json
from collections import defaultdict

class WalletSimilarity:
    def __init__(self, model_name: str = 'all-MiniLM-L6-v2'):
        self.model = SentenceTransformer(model_name)
        self.index = None
        self.wallet_descriptions = {}
        self.wallet_data = None
        
    def load_data(self, csv_path: str = 'data/data.csv'):
        """Загружает данные из CSV и создает описания кошельков"""
        print(f"Loading data from {csv_path}")
        df = pd.read_csv(csv_path)
        print(f"Loaded {len(df)} transactions")
        print(f"Columns in data: {df.columns.tolist()}")
        
        # Проверяем наличие необходимых колонок
        required_columns = ['address', 'to', 'value', 'method']
        missing_columns = [col for col in required_columns if col not in df.columns]
        if missing_columns:
            print(f"Warning: Missing required columns: {missing_columns}")
        
        if 'label' not in df.columns:
            print("Warning: No 'label' column found in data")
        
        self.wallet_data = df
        
        # Группируем транзакции по кошелькам
        wallet_groups = df.groupby('address')
        print(f"Found {len(wallet_groups)} unique wallets")
        
        # Создаем описания для каждого кошелька
        for address, group in wallet_groups:
            description = self._create_wallet_description(group)
            self.wallet_descriptions[address] = description
        
        print(f"Created descriptions for {len(self.wallet_descriptions)} wallets")
        
    def _create_wallet_description(self, transactions: pd.DataFrame) -> str:
        """Создает текстовое описание активности кошелька"""
        # Считаем статистику
        total_txs = len(transactions)
        unique_contracts = transactions['to'].nunique()
        avg_value = transactions['value'].mean()
        methods = transactions['method'].value_counts().to_dict()
        
        # Создаем описание
        description = f"""
        Кошелек совершил {total_txs} транзакций.
        Взаимодействовал с {unique_contracts} уникальными контрактами.
        Среднее значение транзакции: {avg_value:.6f} ETH.
        Методы транзакций: {', '.join([f'{k}: {v}' for k, v in methods.items()])}.
        """
        
        return description
        
    def build_index(self):
        """Строит индекс Faiss для быстрого поиска похожих кошельков"""
        # Получаем эмбеддинги для всех описаний
        descriptions = list(self.wallet_descriptions.values())
        embeddings = self.model.encode(descriptions)
        
        # Создаем и обучаем индекс
        dimension = embeddings.shape[1]
        self.index = faiss.IndexFlatL2(dimension)
        self.index.add(embeddings.astype('float32'))
        
    def find_similar_wallets(self, address: str, k: int = 5) -> List[Dict]:
        """Находит k наиболее похожих кошельков"""
        print(f"Searching for {k} most similar wallets")
        
        # Создаем описание для запрашиваемого кошелька
        description = self._create_wallet_description(pd.DataFrame({
            'to': ['0xfriendwallet000000000000000000000', '0xuniswaprouter000000000000000001'],
            'value': [0.01, 0.05],
            'method': ['send', 'swap']
        }))
        
        # Получаем эмбеддинг для запрашиваемого кошелька
        print("Encoding wallet description...")
        query_embedding = self.model.encode([description])
        
        # Ищем похожие кошельки
        print("Searching in index...")
        distances, indices = self.index.search(query_embedding.astype('float32'), k)
        print(f"Found distances: {distances}, indices: {indices}")
        
        # Преобразуем индексы в адреса
        addresses = list(self.wallet_descriptions.keys())
        similar_addresses = [addresses[i] for i in indices[0]]  # Берем все k ближайших
        print(f"Similar addresses: {similar_addresses}")
        
        # Собираем информацию о похожих кошельках
        results = []
        for addr in similar_addresses:
            wallet_data = self.wallet_data[self.wallet_data['address'] == addr]
            if len(wallet_data) == 0:
                print(f"No data found for address {addr}")
                continue
            
            # Проверяем наличие метки класса
            if 'label' not in wallet_data.columns:
                print(f"Warning: No 'label' column found for address {addr}")
                label = "unknown"
            else:
                label = wallet_data['label'].iloc[0]
                print(f"Found label '{label}' for address {addr}")
            
            results.append({
                'address': addr,
                'label': label,
                'transaction_count': len(wallet_data),
                'description': self.wallet_descriptions[addr]
            })
        
        print(f"Returning {len(results)} similar wallets")
        return results
        
    def save_index(self, path: str = 'wallet_index.faiss'):
        """Сохраняет индекс Faiss"""
        if self.index:
            faiss.write_index(self.index, path)
            
    def load_index(self, path: str = 'wallet_index.faiss'):
        """Загружает индекс Faiss"""
        self.index = faiss.read_index(path)

def main():
    # Инициализируем поисковик
    similarity = WalletSimilarity()
    
    # Загружаем данные
    print("Загрузка данных...")
    similarity.load_data()
    
    # Строим индекс
    print("Построение индекса...")
    similarity.build_index()
    
    # Сохраняем индекс
    similarity.save_index()
    
    # Пример поиска похожих кошельков
    test_address = similarity.wallet_data['address'].iloc[0]
    print(f"\nПоиск похожих кошельков для {test_address}:")
    similar = similarity.find_similar_wallets(test_address)
    
    for wallet in similar:
        print(f"\nАдрес: {wallet['address']}")
        print(f"Тип: {wallet['label']}")
        print(f"Количество транзакций: {wallet['transaction_count']}")
        print(f"Описание: {wallet['description']}")

if __name__ == "__main__":
    main() 