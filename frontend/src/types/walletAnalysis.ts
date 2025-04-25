export interface Transaction {
    hash: string;
    nonce: string;
    transaction_index: string;
    from_address: string;
    from_address_label: string | null;
    to_address: string;
    to_address_label: string | null;
    value: string;
    gas: string;
    gas_price: string;
    receipt_status: string;
    block_timestamp: string;
    block_number: string;
    block_hash: string;
    transaction_fee: string;
}

export interface MLResult {
    confidence: number;
    predicted_class: string;
    similar_wallets: string[] | null;
}

export interface WalletAnalysis {
    status: string;
    message: string;
    cursor: string;
    page_size: number;
    transactions: Transaction[];
    ml_result: MLResult;
}
