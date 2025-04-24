import requests
import time
from typing import List, Dict, Any
import os
from dotenv import load_dotenv

load_dotenv()

class MoralisAPI:
    def __init__(self):
        self.api_key = os.getenv('MORALIS_API_KEY')
        self.base_url = "https://deep-index.moralis.io/api/v2"
        self.headers = {
            "X-API-Key": self.api_key,
            "Content-Type": "application/json"
        }

    def get_transactions(self, address: str, chain: str = "eth", limit: int = 100, cursor: str = None) -> Dict[str, Any]:
        """
        Получает историю транзакций для указанного адреса с поддержкой пагинаци/transactionsи
        
        Args:
            address: Ethereum адрес
            chain: ID блокчейна (по умолчанию "eth")
            limit: Количество транзакций на странице (максимум 100)
            cursor: Курсор для пагинации
         
            print(f"Получено {len(data)} транзакций для адреса {address}")         Dict с транзакциями и метаданными пагинации
        """
        url = f"{self.base_url}/{address}/transactions"
        params = {
            "chain": chain,
            "limit": min(limit, 100)  # Moralis ограничивает максимум 100 транзакций на странице
        }
        
        if cursor:
            params["cursor"] = cursor
            
        try:
            response = requests.get(url, headers=self.headers, params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Ошибка при получении транзакций: {e}")
            return {"result": [], "cursor": None}

    def get_all_transactions(self, address: str, chain: str = "eth", max_pages: int = 10) -> List[Dict[str, Any]]:
        """
        Получает все транзакции для адреса, автоматически обрабатывая пагинацию
        
        Args:
            address: Ethereum адрес
            chain: ID блокчейна
            max_pages: Максимальное количество страниц для получения
            
        Returns:
            List всех транзакций
        """
        all_transactions = []
        cursor = None
        page = 0
        
        while page < max_pages:
            response = self.get_transactions(address, chain, cursor=cursor)
            transactions = response.get("result", [])
            all_transactions.extend(transactions)
            
            cursor = response.get("cursor")
            if not cursor or not transactions:
                break
                
            page += 1
            time.sleep(0.2)  # Небольшая задержка между запросами
            
        return all_transactions

    def get_token_transfers(self, address: str, chain: str = "eth") -> List[Dict[str, Any]]:
        """
        Получает историю трансферов токенов для указанного адреса
        """
        url = f"{self.base_url}/{address}/erc20/transfers"
        params = {
            "chain": chain
        }
        
        try:
            response = requests.get(url, headers=self.headers, params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Ошибка при получении трансферов токенов: {e}")
            return [] 