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
        df = pd.read_csv(csv_path)
        self.wallet_data = df
        
        # Группируем транзакции по кошелькам
        wallet_groups = df.groupby('address')
        
        # Создаем описания для каждого кошелька
        for address, group in wallet_groups:
            description = self._create_wallet_description(group)
            self.wallet_descriptions[address] = description
            
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
        if address not in self.wallet_descriptions:
            return []
            
        # Получаем эмбеддинг для запрашиваемого кошелька
        query_embedding = self.model.encode([self.wallet_descriptions[address]])
        
        # Ищем похожие кошельки
        distances, indices = self.index.search(query_embedding.astype('float32'), k + 1)
        
        # Преобразуем индексы в адреса
        addresses = list(self.wallet_descriptions.keys())
        similar_addresses = [addresses[i] for i in indices[0][1:]]  # Пропускаем первый (сам кошелек)
        
        # Собираем информацию о похожих кошельках
        results = []
        for addr in similar_addresses:
            wallet_data = self.wallet_data[self.wallet_data['address'] == addr]
            results.append({
                'address': addr,
                'label': wallet_data['label'].iloc[0],
                'transaction_count': len(wallet_data),
                'description': self.wallet_descriptions[addr]
            })
            
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