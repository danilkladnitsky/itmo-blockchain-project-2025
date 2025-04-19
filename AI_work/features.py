from typing import Dict, List, Any
from datetime import datetime
import pandas as pd
import numpy as np

class WalletFeatures:
    def __init__(self, transactions: List[Dict[str, Any]], token_transfers: List[Dict[str, Any]]):
        self.transactions = transactions
        self.token_transfers = token_transfers
        self.features = {}

    def extract_features(self) -> Dict[str, Any]:
        """
        Извлекает все фичи из транзакций и трансферов токенов
        """
        self._extract_basic_features()
        self._extract_token_features()
        self._extract_time_features()
        self._extract_nft_features()
        return self.features

    def _extract_basic_features(self):
        """Базовые фичи транзакций"""
        self.features['total_transactions'] = len(self.transactions)
        self.features['total_token_transfers'] = len(self.token_transfers)
        
        # Подсчет входящих и исходящих транзакций
        incoming = sum(1 for tx in self.transactions if tx.get('to_address'))
        outgoing = sum(1 for tx in self.transactions if tx.get('from_address'))
        self.features['incoming_transactions'] = incoming
        self.features['outgoing_transactions'] = outgoing
        self.features['in_out_ratio'] = incoming / (outgoing + 1)  # +1 для избежания деления на 0

    def _extract_token_features(self):
        """Фичи, связанные с токенами"""
        unique_tokens = set()
        token_balances = {}
        
        for transfer in self.token_transfers:
            token_address = transfer.get('token_address')
            value = float(transfer.get('value', 0))
            unique_tokens.add(token_address)
            
            if token_address in token_balances:
                token_balances[token_address] += value
            else:
                token_balances[token_address] = value

        self.features['unique_tokens'] = len(unique_tokens)
        self.features['avg_token_balance'] = np.mean(list(token_balances.values())) if token_balances else 0

    def _extract_time_features(self):
        """Временные фичи"""
        if not self.transactions:
            return
            
        timestamps = [int(tx.get('block_timestamp', 0)) for tx in self.transactions]
        if timestamps:
            dates = [datetime.fromtimestamp(ts) for ts in timestamps]
            time_diffs = [(dates[i+1] - dates[i]).total_seconds() for i in range(len(dates)-1)]
            
            self.features['avg_time_between_tx'] = np.mean(time_diffs) if time_diffs else 0
            self.features['std_time_between_tx'] = np.std(time_diffs) if time_diffs else 0

    def _extract_nft_features(self):
        """Фичи, связанные с NFT"""
        nft_transfers = [tx for tx in self.token_transfers if tx.get('token_type') == 'ERC721']
        self.features['nft_transfers'] = len(nft_transfers)
        self.features['unique_nft_contracts'] = len(set(tx.get('token_address') for tx in nft_transfers)) 