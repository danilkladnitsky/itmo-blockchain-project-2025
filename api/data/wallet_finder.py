import requests
import json
from typing import List, Dict
import time
import os
import pandas as pd
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
        
    def find_drop_hunters(self, limit: int = 50) -> List[Dict]:
        """
        Поиск адресов дропхантеров и их транзакций
        """
        testnet_addresses = self._get_active_testnet_addresses(limit)
        
        transactions_data = []
        for address in testnet_addresses:
            if self._is_drop_hunter(address):
                txs = self.moralis.get_all_transactions(address)
                for tx in txs:
                    transactions_data.append({
                        'address': address,
                        'label': 'drop_hunter',
                        'timestamp': int(tx.get('block_timestamp', 0)),
                        'to': tx.get('to_address', ''),
                        'value': float(tx.get('value', 0)) / 1e18,  # Конвертируем из wei в ETH
                        'method': tx.get('input', '')[:10] if tx.get('input') else 'transfer'
                    })
                if len(transactions_data) >= limit * 10:  # Примерно 10 транзакций на адрес
                    break
            time.sleep(1)
            
        return transactions_data

    def find_nft_collectors(self, limit: int = 50) -> List[Dict]:
        """
        Поиск адресов NFT коллекторов и их транзакций
        """
        nft_addresses = self._get_nft_marketplace_addresses(limit)
        
        transactions_data = []
        for address in nft_addresses:
            if self._is_nft_collector(address):
                txs = self.moralis.get_all_transactions(address)
                for tx in txs:
                    transactions_data.append({
                        'address': address,
                        'label': 'nft_collector',
                        'timestamp': int(tx.get('block_timestamp', 0)),
                        'to': tx.get('to_address', ''),
                        'value': float(tx.get('value', 0)) / 1e18,
                        'method': self._determine_nft_method(tx)
                    })
                if len(transactions_data) >= limit * 10:
                    break
            time.sleep(1)
            
        return transactions_data

    def find_regular_users(self, limit: int = 50) -> List[Dict]:
        """
        Поиск адресов обычных пользователей и их транзакций
        """
        regular_addresses = self._get_regular_activity_addresses(limit)
        
        transactions_data = []
        for address in regular_addresses:
            if self._is_regular_user(address):
                txs = self.moralis.get_all_transactions(address)
                for tx in txs:
                    transactions_data.append({
                        'address': address,
                        'label': 'regular_user',
                        'timestamp': int(tx.get('block_timestamp', 0)),
                        'to': tx.get('to_address', ''),
                        'value': float(tx.get('value', 0)) / 1e18,
                        'method': 'transfer'  # Для обычных пользователей чаще всего простые переводы
                    })
                if len(transactions_data) >= limit * 10:
                    break
            time.sleep(1)
            
        return transactions_data

    def _determine_nft_method(self, tx: Dict) -> str:
        """Определяет метод транзакции для NFT операций"""
        if tx.get('input', '').startswith('0x23b872dd'):  # transferFrom
            return 'transfer'
        elif tx.get('input', '').startswith('0xa0712d68'):  # mint
            return 'mint'
        elif tx.get('to', '').lower() in ['0x7be8076f4ea4a4ad08075c2508e481d6c946d12b',  # OpenSea
                                        '0x7f268357a8c2552623316e2562d90e642bb538e5']:  # Wyvern
            return 'buy' if float(tx.get('value', 0)) > 0 else 'sell'
        return 'transfer'

    def save_to_csv(self, data: List[Dict], filename: str = 'synthetic_wallet_data_with_tech.csv'):
        """Сохраняет данные в CSV файл"""
        df = pd.DataFrame(data)
        df.to_csv(filename, index=False)
        print(f"Данные сохранены в {filename}")

    def _is_drop_hunter(self, address: str) -> bool:
        """
        Проверяет, является ли адрес дропхантером
        """
        transactions = self.moralis.get_transactions(address)
        if not transactions:
            print(f"Нет транзакций для адреса {address}")
            return False
            
        # Проверяем паттерны дропхантеров:
        # 1. Большое количество транзакций в короткий период
        # 2. Множество взаимодействий с тестовыми контрактами
        # 3. Частые взаимодействия с новыми проектами
        
        # Примерные критерии
        tx_count = len(transactions)
        print(f"Всего транзакций: {tx_count}")
        
        # Собираем уникальные адреса получателей
        unique_contracts = set()
        for tx in transactions:
            if isinstance(tx, dict):
                to_address = tx.get('to')
                if to_address:
                    unique_contracts.add(to_address)
        
        unique_count = len(unique_contracts)
        print(f"Уникальных контрактов: {unique_count}")
        
        # Более мягкие критерии для тестирования
        return tx_count > 20 and unique_count > 10

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
            'address': '0x0000000000000000000000000000000000000000',  # Пример адреса
            'startblock': 0,
            'endblock': 99999999,
            'sort': 'desc',
            'apikey': self.etherscan_api_key
        }
        
        try:
            response = requests.get(self.etherscan_base_url, params=params)
            data = response.json()
            
            if data['status'] == '1' and 'result' in data:
                transactions = data['result']
                if isinstance(transactions, list):
                    # Собираем уникальные адреса отправителей и получателей
                    for tx in transactions:
                        if isinstance(tx, dict):
                            if 'from' in tx:
                                addresses.append(tx['from'])
                            if 'to' in tx:
                                addresses.append(tx['to'])
                    
                    # Удаляем дубликаты и ограничиваем количество
                    addresses = list(set(addresses))[:limit]
                else:
                    print("Ошибка: результат не является списком транзакций")
            else:
                print(f"Ошибка API: {data.get('message', 'Неизвестная ошибка')}")
                
        except Exception as e:
            print(f"Ошибка при получении адресов из тестовой сети: {e}")
            
        return addresses

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
    
    # Собираем данные для каждого класса
    print("Сбор данных о дропхантерах...")
    drop_hunter_data = finder.find_drop_hunters(limit=200)
    
    print("Сбор данных о NFT коллекторах...")
    nft_collector_data = finder.find_nft_collectors(limit=200)
    
    print("Сбор данных об обычных пользователях...")
    regular_user_data = finder.find_regular_users(limit=200)
    
    # Объединяем все данные
    all_data = drop_hunter_data + nft_collector_data + regular_user_data
    
    # Сохраняем в CSV
    finder.save_to_csv(all_data)
    
    print("Сбор данных завершен!")
    print(f"Всего собрано {len(all_data)} транзакций")

if __name__ == "__main__":
    main() 