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

    def get_transactions(self, address: str, chain: str = "eth") -> List[Dict[str, Any]]:
        """
        Получает историю транзакций для указанного адреса
        """
        url = f"{self.base_url}/{address}"
        params = {
            "chain": chain
        }
        
        try:
            response = requests.get(url, headers=self.headers, params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"Ошибка при получении транзакций: {e}")
            return []

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