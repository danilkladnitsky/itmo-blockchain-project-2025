import json
import pandas as pd
from typing import List, Dict
from moralis_api import MoralisAPI
from features import WalletFeatures
import time

class DataCollector:
    def __init__(self):
        self.moralis = MoralisAPI()
        self.data = []

    def collect_wallet_data(self, address: str, label: str) -> Dict:
        """
        Собирает данные для одного кошелька
        """
        print(f"Сбор данных для кошелька {address}...")
        
        # Получаем данные через API
        transactions = self.moralis.get_transactions(address)
        token_transfers = self.moralis.get_token_transfers(address)
        
        # Извлекаем фичи
        features = WalletFeatures(transactions, token_transfers).extract_features()
        
        # Добавляем метку
        features['label'] = label
        features['address'] = address
        
        return features

    def collect_dataset(self, addresses: Dict[str, List[str]]):
        """
        Собирает датасет для всех предоставленных адресов
        """
        for label, wallet_addresses in addresses.items():
            for address in wallet_addresses:
                try:
                    wallet_data = self.collect_wallet_data(address, label)
                    self.data.append(wallet_data)
                    print(f"Данные для {address} ({label}) успешно собраны")
                    time.sleep(1)  # Задержка для соблюдения лимитов API
                except Exception as e:
                    print(f"Ошибка при сборе данных для {address}: {e}")
                    continue

    def save_dataset(self, filename: str = "wallet_dataset.csv"):
        """
        Сохраняет собранный датасет в CSV файл
        """
        if not self.data:
            print("Нет данных для сохранения")
            return
            
        df = pd.DataFrame(self.data)
        df.to_csv(filename, index=False)
        print(f"Датасет сохранен в {filename}")

def main():
    # Пример адресов для сбора данных
    addresses = {
        "drop_hunter": [
            "0x123...",  # Замените на реальные адреса дропхантеров
        ],
        "nft_collector": [
            "0x456...",  # Замените на реальные адреса NFT коллекторов
        ],
        "regular_user": [
            "0x789...",  # Замените на реальные адреса обычных пользователей
        ]
    }

    collector = DataCollector()
    collector.collect_dataset(addresses)
    collector.save_dataset()

if __name__ == "__main__":
    main() 