import requests
import json
from typing import List, Dict
import time
import os
from moralis_api import MoralisAPI
from web3 import Web3
from dotenv import load_dotenv

load_dotenv()

class WalletFinder:
    def __init__(self):
        self.moralis = MoralisAPI()
        self.etherscan_api_key = os.getenv('ETHERSCAN_API_KEY')
        self.opensea_api_key = os.getenv('OPENSEA_API_KEY')
        self.etherscan_base_url = "https://api.etherscan.io/api"
        self.opensea_base_url = "https://api.opensea.io/api/v1"
        
    def find_drop_hunters(self, limit: int = 50) -> List[str]:
        """
        Поиск адресов дропхантеров
        """
        # Ищем адреса с высокой активностью в тестовых сетях
        testnet_addresses = self._get_active_testnet_addresses(limit)
        
        # Фильтруем адреса по паттернам дропхантеров
        drop_hunters = []
        for address in testnet_addresses:
            if self._is_drop_hunter(address):
                drop_hunters.append(address)
                if len(drop_hunters) >= limit:
                    break
            time.sleep(1)  # Задержка для API
            
        return drop_hunters

    def find_nft_collectors(self, limit: int = 50) -> List[str]:
        """
        Поиск адресов NFT коллекторов
        """
        # Ищем адреса с высокой активностью на NFT маркетплейсах
        nft_addresses = self._get_nft_marketplace_addresses(limit)
        
        collectors = []
        for address in nft_addresses:
            if self._is_nft_collector(address):
                collectors.append(address)
                if len(collectors) >= limit:
                    break
            time.sleep(1)
            
        return collectors

    def find_regular_users(self, limit: int = 50) -> List[str]:
        """
        Поиск адресов обычных пользователей
        """
        # Ищем адреса с умеренной активностью
        regular_addresses = self._get_regular_activity_addresses(limit)
        
        users = []
        for address in regular_addresses:
            if self._is_regular_user(address):
                users.append(address)
                if len(users) >= limit:
                    break
            time.sleep(1)
            
        return users

    def _is_drop_hunter(self, address: str) -> bool:
        """
        Проверяет, является ли адрес дропхантером
        """
        transactions = self.moralis.get_transactions(address)
        if not transactions:
            return False
            
        # Проверяем паттерны дропхантеров:
        # 1. Большое количество транзакций в короткий период
        # 2. Множество взаимодействий с тестовыми контрактами
        # 3. Частые взаимодействия с новыми проектами
        
        # Примерные критерии
        tx_count = len(transactions)
        unique_contracts = len(set(tx.get('to_address') for tx in transactions))
        
        return tx_count > 100 and unique_contracts > 50

    def _is_nft_collector(self, address: str) -> bool:
        """
        Проверяет, является ли адрес NFT коллектором
        """
        token_transfers = self.moralis.get_token_transfers(address)
        if not token_transfers:
            return False
            
        # Проверяем паттерны NFT коллекторов:
        # 1. Большое количество NFT транзакций
        # 2. Взаимодействия с известными NFT маркетплейсами
        # 3. Длительное хранение NFT
        
        nft_transfers = [tx for tx in token_transfers if tx.get('token_type') == 'ERC721']
        return len(nft_transfers) > 20

    def _is_regular_user(self, address: str) -> bool:
        """
        Проверяет, является ли адрес обычным пользователем
        """
        transactions = self.moralis.get_transactions(address)
        if not transactions:
            return False
            
        # Проверяем паттерны обычных пользователей:
        # 1. Умеренное количество транзакций
        # 2. Регулярная активность
        # 3. Отсутствие специфических паттернов дропхантеров/NFT коллекторов
        
        tx_count = len(transactions)
        return 10 <= tx_count <= 100 and not self._is_drop_hunter(address) and not self._is_nft_collector(address)

    def _get_active_testnet_addresses(self, limit: int) -> List[str]:
        """
        Получает список активных адресов в тестовых сетях
        """
        addresses = []
        
        # Получаем адреса из Goerli тестовой сети
        params = {
            'module': 'account',
            'action': 'txlist',
            'startblock': 0,
            'endblock': 99999999,
            'sort': 'desc',
            'apikey': self.etherscan_api_key
        }
        
        try:
            response = requests.get(self.etherscan_base_url, params=params)
            data = response.json()
            
            if data['status'] == '1':
                transactions = data['result']
                # Собираем уникальные адреса отправителей
                addresses = list(set(tx['from'] for tx in transactions[:limit*2]))
                
        except Exception as e:
            print(f"Ошибка при получении адресов из тестовой сети: {e}")
            
        return addresses[:limit]

    def _get_nft_marketplace_addresses(self, limit: int) -> List[str]:
        """
        Получает список адресов с активностью на NFT маркетплейсах
        """
        addresses = []
        
        # Получаем адреса с активностью на OpenSea
        headers = {
            'X-API-KEY': self.opensea_api_key
        }
        
        try:
            # Получаем последние продажи
            response = requests.get(
                f"{self.opensea_base_url}/events",
                headers=headers,
                params={'event_type': 'successful', 'limit': limit*2}
            )
            data = response.json()
            
            if 'asset_events' in data:
                # Собираем уникальные адреса покупателей
                addresses = list(set(
                    event['winner_account']['address'] 
                    for event in data['asset_events']
                    if 'winner_account' in event
                ))
                
        except Exception as e:
            print(f"Ошибка при получении адресов с NFT активностью: {e}")
            
        return addresses[:limit]

    def _get_regular_activity_addresses(self, limit: int) -> List[str]:
        """
        Получает список адресов с регулярной активностью
        """
        addresses = []
        
        # Получаем адреса с умеренной активностью из Etherscan
        params = {
            'module': 'account',
            'action': 'txlist',
            'startblock': 0,
            'endblock': 99999999,
            'sort': 'desc',
            'apikey': self.etherscan_api_key
        }
        
        try:
            response = requests.get(self.etherscan_base_url, params=params)
            data = response.json()
            
            if data['status'] == '1':
                transactions = data['result']
                # Фильтруем адреса с умеренной активностью
                address_tx_counts = {}
                for tx in transactions:
                    addr = tx['from']
                    address_tx_counts[addr] = address_tx_counts.get(addr, 0) + 1
                
                # Выбираем адреса с 10-100 транзакциями
                addresses = [
                    addr for addr, count in address_tx_counts.items()
                    if 10 <= count <= 100
                ][:limit]
                
        except Exception as e:
            print(f"Ошибка при получении адресов с регулярной активностью: {e}")
            
        return addresses[:limit]

def main():
    finder = WalletFinder()
    
    # Собираем адреса для каждого класса
    drop_hunters = finder.find_drop_hunters(limit=20)
    nft_collectors = finder.find_nft_collectors(limit=20)
    regular_users = finder.find_regular_users(limit=20)
    
    # Сохраняем найденные адреса
    addresses = {
        "drop_hunter": drop_hunters,
        "nft_collector": nft_collectors,
        "regular_user": regular_users
    }
    
    with open("wallet_addresses.json", "w") as f:
        json.dump(addresses, f, indent=4)
    
    print("Найдены адреса:")
    print(f"Дропхантеры: {len(drop_hunters)}")
    print(f"NFT коллекторы: {len(nft_collectors)}")
    print(f"Обычные пользователи: {len(regular_users)}")

if __name__ == "__main__":
    main() 